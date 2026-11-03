// Package progress normalizes PMOCK's per-capability task dialects into one
// customer-facing progress contract.
//
// PMOCK reports progress in at least three different encodings: RPA tasks use
// 0..4, phone lifecycle operations return a taskStatus plus a phoneStatus, and
// UI actions document -1/0/1/2 on paper but return 1/2/3 in practice. None of
// them answers the only question a customer has: what should I do now?
//
// So this package answers exactly that question, with four states:
//
//	running    nothing to do, wait
//	needs_you  reply before the work can continue
//	done       read the result (it may carry items that need attention)
//	failed     decide whether to retry or change approach
//
// Everything else is progress attached to a state: which step, what it is
// called, how much longer, and how strong the evidence behind the claim is.
// The evidence level gates the wording in Copy, so an agent that repeats the
// canonical sentence verbatim cannot invent progress.
package progress

import (
	"fmt"
	"strings"
	"time"
)

// State is the customer-facing state of one handed-over piece of work.
//
// The set stays closed on purpose: a state earns its place only when it changes
// what the customer does next. Two states that lead to the same customer action
// are one state.
type State string

const (
	// StateRunning means the work is moving and the customer should wait.
	StateRunning State = "running"
	// StateNeedsYou means the work is blocked on a decision only the customer
	// can make: a price, a permission, a login, a choice between options.
	StateNeedsYou State = "needs_you"
	// StateDone means there is a result to read.
	StateDone State = "done"
	// StateFailed means there is no usable result.
	StateFailed State = "failed"
)

// Terminal reports whether the state ends the task.
func (s State) Terminal() bool { return s == StateDone || s == StateFailed }

// Evidence is how much the backend actually told us. It decides how specific
// the copy may be; it is never shown to the customer.
type Evidence string

const (
	// EvidenceNone means the last poll carried nothing new: a heartbeat.
	EvidenceNone Evidence = "none"
	// EvidenceStatus means only an opaque status code.
	EvidenceStatus Evidence = "status"
	// EvidenceStep means named business steps, or the failing step.
	EvidenceStep Evidence = "step"
	// EvidenceStructured means a structured result, such as a checklist.
	EvidenceStructured Evidence = "structured"
)

// Counts is a structured result summary, for example "checked 5, found 1".
//
// It is only set when the flow actually produced the numbers. Copy turns it
// into "共检查 5 项，1 项需要关注" and nothing else may guess at them.
type Counts struct {
	Checked int `json:"checked"`
	Issues  int `json:"issues"`
}

// Event is one line of the progress stream.
type Event struct {
	Seq      int            `json:"seq"`
	At       time.Time      `json:"at"`
	Kind     string         `json:"kind"`
	Task     string         `json:"task,omitempty"`
	State    State          `json:"state"`
	Step     int            `json:"step,omitempty"`
	Total    int            `json:"total,omitempty"`
	Label    string         `json:"label,omitempty"`
	ETA      *int           `json:"etaSeconds,omitempty"`
	Counts   *Counts        `json:"counts,omitempty"`
	Evidence Evidence       `json:"evidence"`
	Reason   string         `json:"reason,omitempty"`
	Say      string         `json:"say"`
	Observed map[string]any `json:"observed,omitempty"`
}

// Terminal reports whether this event ends the task.
func (e Event) Terminal() bool { return e.State.Terminal() }

// Observation is what one adapter made of a single poll.
//
// Adapters set State and, when they saw one, a step, a label, a reason or a
// structured result. They never format text: the wording lives in Copy so that
// every capability sounds the same.
type Observation struct {
	State    State
	Step     int
	Total    int
	Label    string
	Reason   string
	Counts   *Counts
	Observed map[string]any
}

// Copy renders the canonical customer-facing sentence for an event.
//
// The rule enforced here is the whole point of the package: the wording may
// only be as specific as the evidence. Step numbers and percentages appear only
// when the adapter saw a step; a bare status produces a bare sentence; a
// heartbeat says out loud that there is nothing new, instead of pretending.
func Copy(e Event) string {
	switch e.State {
	case StateDone:
		if e.Counts != nil {
			return fmt.Sprintf("完成了：共检查 %d 项，%d 项需要关注。", e.Counts.Checked, e.Counts.Issues)
		}
		return "完成了。"
	case StateFailed:
		reason := cleanReason(e.Reason)
		if reason == "" {
			return "这次没做成。想让我重试吗？"
		}
		return fmt.Sprintf("没做成：%s。想让我重试吗？", reason)
	case StateNeedsYou:
		reason := cleanReason(e.Reason)
		if reason == "" {
			return "需要你确认一下，我才能继续。"
		}
		return fmt.Sprintf("需要你先处理：%s。处理好告诉我，我接着做。", reason)
	}

	// running
	if e.Evidence == EvidenceNone {
		return "还在进行，暂时没有新的进度。"
	}
	name := cleanReason(e.Label)
	var b strings.Builder
	if e.Evidence == EvidenceStep && e.Step > 0 && e.Total > 0 {
		fmt.Fprintf(&b, "正在做第 %d 步（共 %d 步）", e.Step, e.Total)
		if name != "" {
			fmt.Fprintf(&b, "：%s", name)
		}
	} else if name != "" {
		fmt.Fprintf(&b, "正在%s", name)
	} else {
		b.WriteString("正在处理")
	}
	if e.ETA != nil && *e.ETA > 0 {
		fmt.Fprintf(&b, "，预计还要 %s", Human(*e.ETA))
	}
	b.WriteString("。")
	return b.String()
}

// cleanReason trims a backend-supplied string and drops values that read like
// machine output rather than a sentence a customer can act on.
func cleanReason(s string) string {
	s = strings.TrimSpace(s)
	if s == "" || s == "null" || s == "-" {
		return ""
	}
	return s
}

// Human renders a duration in seconds the way a person would say it.
func Human(seconds int) string {
	if seconds <= 0 {
		return "0 秒"
	}
	if seconds < 60 {
		return fmt.Sprintf("%d 秒", seconds)
	}
	minutes := (seconds + 30) / 60
	if minutes < 1 {
		minutes = 1
	}
	if minutes < 60 {
		return fmt.Sprintf("%d 分钟", minutes)
	}
	hours := minutes / 60
	rest := minutes % 60
	if rest == 0 {
		return fmt.Sprintf("%d 小时", hours)
	}
	return fmt.Sprintf("%d 小时 %d 分钟", hours, rest)
}
