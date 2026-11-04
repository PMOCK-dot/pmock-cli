package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/haidaisiganbei/pmock-cli/internal/build"
	"github.com/haidaisiganbei/pmock-cli/internal/client"
	"github.com/haidaisiganbei/pmock-cli/internal/exitcode"
	"github.com/haidaisiganbei/pmock-cli/internal/progress"
)

const watchLong = `watch turns PMOCK's per-capability task status into one progress stream.

Each capability reports progress in its own dialect: RPA uses 0..4, phone
lifecycle operations return a taskStatus plus a phoneStatus, UI actions return
1/2/3 in practice. watch normalizes all of them and emits one JSON object per
line, with a sentence a customer can read.

STATES
  running    the work is moving; the customer waits
  needs_you  the work is blocked on a decision only the customer can make
  done       there is a result to read
  failed     there is no usable result

EVENTS
  Unchanged polls are silent. After 30 seconds without news the stream says so
  explicitly instead of going quiet. A timeout emits a failed event before
  exiting, so the stream never stops without an explanation.

  "say" is the canonical customer sentence. Its specificity never exceeds the
  evidence the backend actually gave ("evidence"): step numbers appear only when
  a step was really reported. Use it verbatim, or use it as the floor.

EXIT CODES
  0  done, needs_you, or --once returned a snapshot
  1  the task ended in failed
  2  usage error
  3  timeout, network or API failure`

func newWatchCmd(g *globals) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Follow a long-running task and print one progress event per change",
		Long:  watchLong,
	}
	cmd.AddCommand(
		newWatchRPACmd(g),
		newWatchPhoneCmd(g),
		newWatchActionCmd(g),
		newWatchLogCmd(),
	)
	return cmd
}

// watchOpts are the flags shared by every watch target. The behaviour they
// configure lives in internal/progress.Runner.
type watchOpts struct {
	interval  time.Duration
	heartbeat time.Duration
	timeout   time.Duration
	once      bool
	observe   time.Duration
	human     bool
	record    bool
	stateDir  string
}

func (o *watchOpts) bind(c *cobra.Command) {
	c.Flags().DurationVar(&o.interval, "interval", 3*time.Second, "polling interval")
	c.Flags().DurationVar(&o.heartbeat, "heartbeat", 30*time.Second,
		"emit a 'nothing new' event after this much silence (negative disables)")
	c.Flags().DurationVar(&o.timeout, "deadline", 20*time.Minute,
		"give up after this long and emit a failed event")
	c.Flags().BoolVar(&o.once, "once", false, "poll once and exit")
	c.Flags().DurationVar(&o.observe, "observe", 0,
		"long poll: wait up to this long for a change, return the current event either way (0 = stream until done)")
	c.Flags().BoolVar(&o.human, "human", false, "print the canonical sentence instead of a JSON line")
	c.Flags().BoolVar(&o.record, "record", false, "append every event to the local ledger")
	c.Flags().StringVar(&o.stateDir, "state-dir", "",
		"ledger directory (default $PMOCK_CLI_WATCH_DIR or ~/.pmock-cli/watch)")
}

func (o *watchOpts) run(cmd *cobra.Command, key string, ad progress.Adapter) error {
	led, err := o.openLedger(key)
	if err != nil {
		return err
	}
	if led != nil {
		defer led.Close()
	}

	observe := o.observe
	if o.once {
		observe = 0
	}
	out := cmd.OutOrStdout()
	r := &progress.Runner{
		Adapter:   ad,
		Interval:  o.interval,
		Heartbeat: o.heartbeat,
		Timeout:   o.timeout,
		Once:      o.once,
		Observe:   observe,
		Emit: func(ev progress.Event) error {
			if err := writeWatchEvent(out, ev, o.human); err != nil {
				return err
			}
			if led == nil {
				return nil
			}
			return led.Write(ev)
		},
	}

	last, runErr := r.Run(cmd.Context())
	if runErr != nil {
		var coded *exitcode.Error
		if errors.As(runErr, &coded) {
			return runErr
		}
		switch {
		case errors.Is(runErr, progress.ErrTimeout):
			return exitcode.Runtimef("watch timed out; the last event explains why")
		case errors.Is(runErr, context.Canceled), errors.Is(runErr, context.DeadlineExceeded):
			return nil
		}
		var apiErr *client.APIError
		if errors.As(runErr, &apiErr) {
			return exitcode.Businessf("%v", runErr)
		}
		return exitcode.Runtimef("%v", runErr)
	}

	if last.State == progress.StateFailed {
		// The failed event on stdout already carries the explanation; keep
		// stderr clean so streaming consumers do not see a second report.
		return exitcode.Businessf("the task ended in failed")
	}
	return nil
}

func writeWatchEvent(w io.Writer, ev progress.Event, human bool) error {
	if human {
		_, err := fmt.Fprintf(w, "%s %s\n", ev.At.Format("15:04:05"), ev.Say)
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(ev); err != nil {
		return exitcode.Runtimef("failed to write event: %v", err)
	}
	return nil
}

func newWatchRPACmd(g *globals) *cobra.Command {
	var taskID int64
	o := &watchOpts{}
	c := &cobra.Command{
		Use:   "rpa",
		Short: "Follow an RPA task (the only target that can report named steps)",
		Long: `Follow an RPA task by the id returned from "rpa task add".

RPA reports status 0 待执行, 1 进行中, 2 成功, 3 部分失败, 4 失败. When the task
carries a node list, the failing node names the step in the customer sentence;
without a node list the sentence stays generic.`,
		Example: `  pmock-cli watch rpa --task-id 2117537701360017409
  pmock-cli watch rpa --task-id 2117537701360017409 --human --record`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if taskID <= 0 {
				return exitcode.Usagef("--task-id is required")
			}
			api, err := watchClient(g)
			if err != nil {
				return err
			}
			return o.run(cmd, fmt.Sprintf("rpa-%d", taskID), &progress.RPAAdapter{API: api, TaskID: taskID})
		},
	}
	c.Flags().Int64Var(&taskID, "task-id", 0, "RPA task id")
	o.bind(c)
	return c
}

func newWatchPhoneCmd(g *globals) *cobra.Command {
	var (
		phoneID int64
		op      string
	)
	o := &watchOpts{}
	c := &cobra.Command{
		Use:   "phone",
		Short: "Follow a phone lifecycle operation (start, stop, restart, reset, renew)",
		Long: `Follow a phone operation by phone id and operate type.

An operation only counts as done when the phone itself reached the expected
state; until then the stream keeps saying running rather than reporting success
too early.`,
		Example: `  pmock-cli watch phone --id 123 --op powerOn`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if phoneID <= 0 {
				return exitcode.Usagef("--id is required")
			}
			switch op {
			case "powerOn", "powerOff", "restart", "reset", "renew":
			default:
				return exitcode.Usagef("--op must be one of powerOn, powerOff, restart, reset, renew")
			}
			api, err := watchClient(g)
			if err != nil {
				return err
			}
			return o.run(cmd, fmt.Sprintf("phone-%d", phoneID),
				&progress.PhoneAdapter{API: api, PhoneID: phoneID, OperateType: op})
		},
	}
	c.Flags().Int64Var(&phoneID, "id", 0, "phone id")
	c.Flags().StringVar(&op, "op", "", "operate type: powerOn, powerOff, restart, reset, renew")
	o.bind(c)
	return c
}

func newWatchActionCmd(g *globals) *cobra.Command {
	var (
		code   string
		taskID int64
	)
	o := &watchOpts{}
	c := &cobra.Command{
		Use:   "action",
		Short: "Follow one UI action task (click, swipe, input, ...)",
		Long: `Follow a UI action by phone code and action task id.

Actions usually finish in well under a second. This target exists for the rare
long chain, and it exists so that every capability goes through the same four
states instead of inventing a fifth.`,
		Example: `  pmock-cli watch action --code ABC123 --task-id 456`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if code == "" {
				return exitcode.Usagef("--code is required")
			}
			if taskID <= 0 {
				return exitcode.Usagef("--task-id is required")
			}
			api, err := watchClient(g)
			if err != nil {
				return err
			}
			return o.run(cmd, fmt.Sprintf("action-%d", taskID),
				&progress.ActionAdapter{API: api, Code: code, TaskID: taskID})
		},
	}
	c.Flags().StringVar(&code, "code", "", "phone code (the six-digit number shown by `phone list`)")
	c.Flags().Int64Var(&taskID, "task-id", 0, "action task id")
	o.bind(c)
	return c
}

func watchClient(g *globals) (*client.Client, error) {
	cfg, err := g.config()
	if err != nil {
		return nil, err
	}
	c := client.New(cfg.BaseURL, cfg.AccessKey, cfg.SecretKey, cfg.Timeout)
	c.UserAgent = "pmock-cli/" + build.Version
	return c, nil
}

// ledger is the local record of emitted events. It exists so a customer (or an
// agent) can come back later and replay what was said, instead of re-running
// the work to find out how it ended.
type ledger struct {
	f *os.File
}

func (o *watchOpts) ledgerPath(key string) (string, error) {
	dir := o.stateDir
	if dir == "" {
		dir = os.Getenv("PMOCK_CLI_WATCH_DIR")
	}
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", exitcode.Runtimef("cannot locate the home directory: %v", err)
		}
		dir = filepath.Join(home, ".pmock-cli", "watch")
	}
	return filepath.Join(dir, key+".jsonl"), nil
}

func (o *watchOpts) openLedger(key string) (*ledger, error) {
	if !o.record {
		return nil, nil
	}
	path, err := o.ledgerPath(key)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, exitcode.Runtimef("cannot create the ledger directory: %v", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, exitcode.Runtimef("cannot open the ledger: %v", err)
	}
	return &ledger{f: f}, nil
}

func (l *ledger) Write(ev progress.Event) error {
	b, err := json.Marshal(ev)
	if err != nil {
		return exitcode.Runtimef("cannot encode the ledger entry: %v", err)
	}
	if _, err := l.f.Write(append(b, '\n')); err != nil {
		return exitcode.Runtimef("cannot write the ledger: %v", err)
	}
	return nil
}

// Close is deliberately best-effort: a failed fsync must not turn a completed
// task into a failed command.
func (l *ledger) Close() { _ = l.f.Close() }

func newWatchLogCmd() *cobra.Command {
	var (
		kind     string
		taskID   int64
		stateDir string
	)
	c := &cobra.Command{
		Use:   "log",
		Short: "Replay the recorded event stream of a watched task",
		Long: `Replay what was already recorded by "watch --record".

The same event shape is written, so the replay can be piped through the same
consumer as a live stream.`,
		Example: `  pmock-cli watch log --kind rpa --task-id 2117537701360017409`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if taskID <= 0 {
				return exitcode.Usagef("--task-id is required")
			}
			switch kind {
			case "rpa", "phone", "action":
			default:
				return exitcode.Usagef("--kind must be one of rpa, phone, action")
			}
			o := &watchOpts{stateDir: stateDir}
			path, err := o.ledgerPath(fmt.Sprintf("%s-%d", kind, taskID))
			if err != nil {
				return err
			}
			f, err := os.Open(path)
			if err != nil {
				if os.IsNotExist(err) {
					return exitcode.Usagef("no ledger for %s %d (was it watched with --record?)", kind, taskID)
				}
				return exitcode.Runtimef("cannot open the ledger: %v", err)
			}
			defer f.Close()
			if _, err := io.Copy(cmd.OutOrStdout(), f); err != nil {
				return exitcode.Runtimef("cannot read the ledger: %v", err)
			}
			return nil
		},
	}
	c.Flags().StringVar(&kind, "kind", "", "capability: rpa, phone or action")
	c.Flags().Int64Var(&taskID, "task-id", 0, "task id")
	c.Flags().StringVar(&stateDir, "state-dir", "",
		"ledger directory (default $PMOCK_CLI_WATCH_DIR or ~/.pmock-cli/watch)")
	return c
}
