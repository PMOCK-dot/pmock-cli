package progress

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type scriptedAdapter struct {
	kind string
	obs  []Observation
	i    int
	err  error
}

func (s *scriptedAdapter) Kind() string { return s.kind }

func (s *scriptedAdapter) Poll(ctx context.Context) (Observation, error) {
	if s.err != nil {
		return Observation{}, s.err
	}
	if s.i >= len(s.obs) {
		return s.obs[len(s.obs)-1], nil
	}
	o := s.obs[s.i]
	s.i++
	return o, nil
}

// fakeClock makes the runner's heartbeat, ETA and timeout deterministic.
type fakeClock struct {
	t      time.Time
	max    int
	cancel context.CancelFunc
	n      int
}

func (c *fakeClock) Now() time.Time { return c.t }

func (c *fakeClock) Sleep(ctx context.Context, d time.Duration) bool {
	c.n++
	if c.max > 0 && c.n >= c.max {
		if c.cancel != nil {
			c.cancel()
		}
		return false
	}
	c.t = c.t.Add(d)
	return true
}

func newRunner(ad Adapter, clock *fakeClock, emit func(Event) error) *Runner {
	return &Runner{
		Adapter:   ad,
		Interval:  3 * time.Second,
		Heartbeat: -1,
		Timeout:   time.Minute,
		Emit:      emit,
		Now:       clock.Now,
		Sleep:     clock.Sleep,
	}
}

func TestRunnerEmitsOnlyOnChange(t *testing.T) {
	ad := &scriptedAdapter{kind: "rpa", obs: []Observation{
		{State: StateRunning},
		{State: StateRunning},
		{State: StateRunning, Step: 2, Total: 3, Label: "检查资料"},
		{State: StateDone},
	}}
	clock := &fakeClock{t: time.Unix(0, 0)}
	var events []Event
	r := newRunner(ad, clock, func(e Event) error { events = append(events, e); return nil })

	last, err := r.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 3 {
		t.Fatalf("got %d events, want 3: %+v", len(events), events)
	}
	if events[0].State != StateRunning || events[0].Seq != 1 {
		t.Fatalf("first = %+v", events[0])
	}
	if events[1].Step != 2 || events[1].Evidence != EvidenceStep {
		t.Fatalf("second = %+v", events[1])
	}
	if events[1].Say != "正在做第 2 步（共 3 步）：检查资料，预计还要 3 秒。" {
		t.Fatalf("say = %q", events[1].Say)
	}
	if !last.Terminal() || last.State != StateDone {
		t.Fatalf("last = %+v", last)
	}
}

func TestRunnerHeartbeatSaysNothingNew(t *testing.T) {
	ad := &scriptedAdapter{kind: "rpa", obs: []Observation{{State: StateRunning}}}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	clock := &fakeClock{t: time.Unix(0, 0), max: 6, cancel: cancel}
	var events []Event
	r := newRunner(ad, clock, func(e Event) error { events = append(events, e); return nil })
	r.Heartbeat = 10 * time.Second

	if _, err := r.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if len(events) < 2 {
		t.Fatalf("expected a heartbeat after silence, got %d events", len(events))
	}
	hb := events[1]
	if hb.Evidence != EvidenceNone || hb.Step != 0 || hb.ETA != nil {
		t.Fatalf("heartbeat kept specific claims: %+v", hb)
	}
	if hb.Say != "还在进行，暂时没有新的进度。" {
		t.Fatalf("say = %q", hb.Say)
	}
}

func TestRunnerTimeoutEndsWithAFailedEvent(t *testing.T) {
	ad := &scriptedAdapter{kind: "rpa", obs: []Observation{{State: StateRunning}}}
	clock := &fakeClock{t: time.Unix(0, 0)}
	var events []Event
	r := newRunner(ad, clock, func(e Event) error { events = append(events, e); return nil })
	r.Timeout = 10 * time.Second

	last, err := r.Run(context.Background())
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("err = %v, want ErrTimeout", err)
	}
	if last.State != StateFailed || !strings.Contains(last.Reason, "没跑完") {
		t.Fatalf("last = %+v", last)
	}
	if !strings.Contains(last.Say, "没做成") {
		t.Fatalf("say = %q", last.Say)
	}
	if len(events) != 2 {
		t.Fatalf("got %d events, want the opening event plus the timeout", len(events))
	}
}

func TestRunnerStopsAtNeedsYou(t *testing.T) {
	ad := &scriptedAdapter{kind: "rpa", obs: []Observation{
		{State: StateRunning},
		{State: StateNeedsYou, Reason: "账号登录已失效"},
	}}
	clock := &fakeClock{t: time.Unix(0, 0)}
	var events []Event
	r := newRunner(ad, clock, func(e Event) error { events = append(events, e); return nil })

	last, err := r.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if last.State != StateNeedsYou {
		t.Fatalf("last = %+v", last)
	}
	if clock.n != 1 {
		t.Fatalf("runner kept polling after needs_you: slept %d times", clock.n)
	}
}

func TestRunnerEstimatesRemainingTime(t *testing.T) {
	ad := &scriptedAdapter{kind: "rpa", obs: []Observation{
		{State: StateRunning, Step: 1, Total: 4, Label: "第一步"},
		{State: StateRunning, Step: 2, Total: 4, Label: "第二步"},
		{State: StateDone},
	}}
	clock := &fakeClock{t: time.Unix(0, 0)}
	var events []Event
	r := newRunner(ad, clock, func(e Event) error { events = append(events, e); return nil })

	if _, err := r.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if events[1].ETA == nil {
		t.Fatalf("second event has no estimate: %+v", events[1])
	}
	if got := *events[1].ETA; got != 3 {
		t.Fatalf("eta = %d, want 3 (one step took 3s, two steps left)", got)
	}
}

func TestRunnerOnceReturnsTheFirstSnapshot(t *testing.T) {
	ad := &scriptedAdapter{kind: "phone", obs: []Observation{{State: StateRunning}}}
	clock := &fakeClock{t: time.Unix(0, 0)}
	r := newRunner(ad, clock, func(Event) error { return nil })
	r.Once = true

	if _, err := r.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if clock.n != 0 {
		t.Fatalf("--once slept %d times", clock.n)
	}
}
