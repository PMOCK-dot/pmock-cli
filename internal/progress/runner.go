package progress

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/PMOCK-dot/pmock-cli/internal/client"
)

// ErrTimeout means the task did not reach a terminal state before the deadline.
// The runner still emits a failed event first, so the stream never just stops.
var ErrTimeout = errors.New("timed out")

const (
	defaultInterval  = 3 * time.Second
	defaultHeartbeat = 30 * time.Second
	defaultTimeout   = 20 * time.Minute
)

// Runner polls one adapter and emits events when the answer changes.
//
// Three behaviours matter more than the polling itself:
//
//   - unchanged polls are silent, so a 10 minute flow does not produce 200 lines
//     that all say "still running";
//   - silence for longer than Heartbeat produces an explicit "nothing new"
//     event, because a quiet stream and a dead stream look identical to a
//     customer;
//   - a timeout emits a terminal failed event instead of ending the stream with
//     no explanation.
type Runner struct {
	Adapter   Adapter
	Interval  time.Duration
	Heartbeat time.Duration // negative disables the heartbeat
	Timeout   time.Duration
	Once      bool
	// Observe turns Run into a long poll: it waits up to this long for the
	// state to change, then returns the current event either way. Agents use
	// it to narrate a long task without hammering the API or going silent.
	Observe time.Duration
	Emit    func(Event) error
	Now     func() time.Time
	Sleep   func(context.Context, time.Duration) bool
}

// Run polls until a terminal state, the deadline, a cancelled context or an
// error. It returns the last event it emitted.
func (r *Runner) Run(ctx context.Context) (Event, error) {
	now := r.Now
	if now == nil {
		now = time.Now
	}
	sleep := r.Sleep
	if sleep == nil {
		sleep = sleepCtx
	}
	interval := r.Interval
	if interval <= 0 {
		interval = defaultInterval
	}
	heartbeat := r.Heartbeat
	if heartbeat == 0 {
		heartbeat = defaultHeartbeat
	}
	if heartbeat < 0 {
		heartbeat = 0
	}
	timeout := r.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	if r.Observe > 0 {
		return r.observe(ctx, now, sleep, interval)
	}

	start := now()
	var last Event
	var lastEmitAt time.Time
	seq := 0
	first := true

	for {
		obs, err := r.Adapter.Poll(ctx)
		if err != nil {
			return r.pollFailed(last, seq, now, err)
		}
		if cerr := ctx.Err(); cerr != nil {
			return last, cerr
		}

		ev := r.event(obs, now())
		if ev.State == StateRunning && ev.Step > 0 && ev.Total > ev.Step {
			if elapsed := now().Sub(start).Seconds(); elapsed > 0 {
				eta := int(math.Ceil(elapsed / float64(ev.Step) * float64(ev.Total-ev.Step)))
				if eta > 0 {
					ev.ETA = &eta
				}
			}
		}

		out := ev
		emit := first || materialChange(last, ev)
		if !emit && ev.Terminal() {
			emit = true
		}
		if !emit && ev.State == StateRunning && heartbeat > 0 && now().Sub(lastEmitAt) >= heartbeat {
			// Nothing new to report: say so, and drop every specific claim.
			out.Evidence = EvidenceNone
			out.Step, out.Total, out.Label, out.Counts, out.ETA = 0, 0, "", nil, nil
			emit = true
		}

		if emit {
			seq++
			out.Seq = seq
			out.Say = Copy(out)
			if r.Emit != nil {
				if err := r.Emit(out); err != nil {
					return out, err
				}
			}
			lastEmitAt = now()
			first = false
			last = ev
			// needs_you is not terminal, but watching must stop all the same:
			// the only way forward is a human reply.
			if out.Terminal() || out.State == StateNeedsYou || r.Once {
				return out, nil
			}
		} else {
			last = ev
		}

		if now().Sub(start) >= timeout {
			late := Event{
				At:       now(),
				Kind:     r.Adapter.Kind(),
				Seq:      seq + 1,
				State:    StateFailed,
				Evidence: EvidenceStatus,
				Reason:   fmt.Sprintf("%s 内没跑完", Human(int(timeout.Seconds()))),
			}
			late.Say = Copy(late)
			if r.Emit != nil {
				if err := r.Emit(late); err != nil {
					return late, err
				}
			}
			return late, ErrTimeout
		}
		if !sleep(ctx, interval) {
			return last, ctx.Err()
		}
	}
}

func evidenceOf(o Observation) Evidence {
	switch {
	case o.Counts != nil:
		return EvidenceStructured
	case o.Step > 0 || o.Total > 0:
		return EvidenceStep
	default:
		return EvidenceStatus
	}
}

// pollFailed turns a failed poll into a terminal failed event before the error
// reaches the caller. Without it a watch that cannot read the task at all would
// exit non-zero with an empty stdout, and a consumer (the MCP server above all)
// would see "no event" rather than "this task failed" — the opposite of the
// promise that the stream never stops without an explanation.
//
// A cancelled context stays silent: the caller asked to stop, that is not a
// failure of the work.
func (r *Runner) pollFailed(last Event, seq int, now func() time.Time, err error) (Event, error) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return last, err
	}
	ev := Event{
		At:       now(),
		Kind:     r.Adapter.Kind(),
		Seq:      seq + 1,
		State:    StateFailed,
		Evidence: EvidenceStatus,
		Reason:   pollFailureReason(err),
	}
	ev.Say = Copy(ev)
	if r.Emit != nil {
		if emitErr := r.Emit(ev); emitErr != nil {
			return ev, emitErr
		}
	}
	return ev, err
}

// pollFailureReason keeps the customer sentence readable: an API error already
// carries its message, and its code and trace id belong in a log, not in copy.
func pollFailureReason(err error) string {
	var apiErr *client.APIError
	if errors.As(err, &apiErr) && strings.TrimSpace(apiErr.Msg) != "" {
		return strings.TrimSpace(apiErr.Msg)
	}
	msg := strings.TrimSpace(err.Error())
	if msg == "" || msg == "null" {
		return ""
	}
	return msg
}

func materialChange(prev, next Event) bool {
	if prev.State != next.State || prev.Step != next.Step || prev.Total != next.Total {
		return true
	}
	if prev.Label != next.Label || prev.Reason != next.Reason {
		return true
	}
	if (prev.Counts == nil) != (next.Counts == nil) {
		return true
	}
	if prev.Counts != nil && next.Counts != nil && *prev.Counts != *next.Counts {
		return true
	}
	return false
}

func sleepCtx(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

// event builds the normalized event for one observation.
func (r *Runner) event(obs Observation, at time.Time) Event {
	return Event{
		At:       at,
		Kind:     r.Adapter.Kind(),
		State:    obs.State,
		Step:     obs.Step,
		Total:    obs.Total,
		Label:    obs.Label,
		Reason:   obs.Reason,
		Counts:   obs.Counts,
		Observed: obs.Observed,
		Evidence: evidenceOf(obs),
	}
}

// observe is the long-poll mode: wait for a change, but never longer than
// Observe, then return the current event. The caller decides whether the event
// is worth repeating; this mode never invents a "nothing new" sentence.
func (r *Runner) observe(ctx context.Context, now func() time.Time, sleep func(context.Context, time.Duration) bool, interval time.Duration) (Event, error) {
	start := now()
	var base Event
	first := true
	for {
		obs, err := r.Adapter.Poll(ctx)
		if err != nil {
			return r.pollFailed(base, 0, now, err)
		}
		ev := r.event(obs, now())
		if first {
			base = ev
			first = false
		}
		// A task that already ended must not make the caller wait out the
		// window: there is nothing left to observe.
		if ev.Terminal() || ev.State == StateNeedsYou {
			ev.Seq = 1
			ev.Say = Copy(ev)
			if r.Emit != nil {
				if err := r.Emit(ev); err != nil {
					return ev, err
				}
			}
			return ev, nil
		}
		if materialChange(base, ev) || now().Sub(start) >= r.Observe {
			ev.Seq = 1
			ev.Say = Copy(ev)
			if r.Emit != nil {
				if err := r.Emit(ev); err != nil {
					return ev, err
				}
			}
			return ev, nil
		}
		if !sleep(ctx, interval) {
			return base, ctx.Err()
		}
	}
}
