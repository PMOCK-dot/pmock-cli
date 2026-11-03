package progress

import (
	"strings"
	"testing"
)

func intp(v int) *int { return &v }

func TestCopyStates(t *testing.T) {
	tests := []struct {
		name string
		ev   Event
		want string
	}{
		{
			name: "running with steps and an estimate",
			ev: Event{State: StateRunning, Evidence: EvidenceStep, Step: 2, Total: 3,
				Label: "检查账号页面", ETA: intp(40)},
			want: "正在做第 2 步（共 3 步）：检查账号页面，预计还要 40 秒。",
		},
		{
			name: "running with a step but no estimate",
			ev:   Event{State: StateRunning, Evidence: EvidenceStep, Step: 1, Total: 3, Label: "打开应用"},
			want: "正在做第 1 步（共 3 步）：打开应用。",
		},
		{
			name: "running with only a status label",
			ev:   Event{State: StateRunning, Evidence: EvidenceStatus, Label: "开机中"},
			want: "正在开机中。",
		},
		{
			name: "running with nothing at all",
			ev:   Event{State: StateRunning, Evidence: EvidenceStatus},
			want: "正在处理。",
		},
		{
			name: "heartbeat",
			ev:   Event{State: StateRunning, Evidence: EvidenceNone},
			want: "还在进行，暂时没有新的进度。",
		},
		{
			name: "done with a structured result",
			ev:   Event{State: StateDone, Evidence: EvidenceStructured, Counts: &Counts{Checked: 5, Issues: 1}},
			want: "完成了：共检查 5 项，1 项需要关注。",
		},
		{
			name: "done without a result",
			ev:   Event{State: StateDone, Evidence: EvidenceStatus},
			want: "完成了。",
		},
		{
			name: "failed with a reason",
			ev:   Event{State: StateFailed, Evidence: EvidenceStep, Reason: "在「检查资料」这一步失败了"},
			want: "没做成：在「检查资料」这一步失败了。想让我重试吗？",
		},
		{
			name: "failed without a reason",
			ev:   Event{State: StateFailed, Evidence: EvidenceStatus},
			want: "这次没做成。想让我重试吗？",
		},
		{
			name: "needs you with a reason",
			ev:   Event{State: StateNeedsYou, Evidence: EvidenceStatus, Reason: "账号登录已失效"},
			want: "需要你先处理：账号登录已失效。处理好告诉我，我接着做。",
		},
		{
			name: "needs you without a reason",
			ev:   Event{State: StateNeedsYou, Evidence: EvidenceStatus},
			want: "需要你确认一下，我才能继续。",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Copy(tc.ev); got != tc.want {
				t.Fatalf("Copy() = %q, want %q", got, tc.want)
			}
		})
	}
}

// The whole point of the contract: a step number may only appear in the copy
// when the evidence supports it. An adapter that guesses a step must not leak
// that guess into a customer sentence.
func TestCopyNeverInventsSteps(t *testing.T) {
	ev := Event{State: StateRunning, Evidence: EvidenceStatus, Step: 3, Total: 5, Label: "检查资料"}
	got := Copy(ev)
	if strings.Contains(got, "第 3 步") || strings.Contains(got, "共 5 步") {
		t.Fatalf("copy leaked an unsupported step: %q", got)
	}
	if strings.Contains(got, "3") || strings.Contains(got, "5") {
		t.Fatalf("copy leaked unsupported numbers: %q", got)
	}
}

// A heartbeat must never carry numbers either, even if the event still has
// stale fields attached.
func TestCopyHeartbeatDropsNumbers(t *testing.T) {
	ev := Event{State: StateRunning, Evidence: EvidenceNone, Step: 2, Total: 3, ETA: intp(30)}
	got := Copy(ev)
	if got != "还在进行，暂时没有新的进度。" {
		t.Fatalf("heartbeat copy = %q", got)
	}
}

func TestHuman(t *testing.T) {
	cases := map[int]string{
		0:    "0 秒",
		5:    "5 秒",
		59:   "59 秒",
		60:   "1 分钟",
		89:   "1 分钟",
		90:   "2 分钟",
		600:  "10 分钟",
		3600: "1 小时",
		3900: "1 小时 5 分钟",
	}
	for in, want := range cases {
		if got := Human(in); got != want {
			t.Errorf("Human(%d) = %q, want %q", in, got, want)
		}
	}
}
