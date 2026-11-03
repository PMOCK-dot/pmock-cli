package progress

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/haidaisiganbei/pmock-cli/internal/client"
)

type fakeDoer struct {
	data  string
	err   error
	calls int
	req   client.Request
}

func (f *fakeDoer) Do(ctx context.Context, req client.Request) (*client.Envelope, error) {
	f.calls++
	f.req = req
	if f.err != nil {
		return nil, f.err
	}
	return &client.Envelope{Code: 0, Data: json.RawMessage(f.data)}, nil
}

func TestRPARunning(t *testing.T) {
	api := &fakeDoer{data: `{"status":1,"taskList":[{"status":1}]}`}
	ad := &RPAAdapter{API: api, TaskID: 42}
	obs, err := ad.Poll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if obs.State != StateRunning {
		t.Fatalf("state = %s", obs.State)
	}
	if want := "/open/api/rpa/task/getInfo/42"; api.req.Path != want {
		t.Fatalf("path = %s, want %s", api.req.Path, want)
	}
}

func TestRPAPending(t *testing.T) {
	api := &fakeDoer{data: `{"status":0}`}
	obs, err := (&RPAAdapter{API: api, TaskID: 1}).Poll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if obs.State != StateRunning || obs.Label != "排队中" {
		t.Fatalf("obs = %+v", obs)
	}
}

func TestRPADoneUsesLastNode(t *testing.T) {
	api := &fakeDoer{data: `{"status":2,"taskList":[{"status":2,"nodeList":[
		{"idx":0,"name":"打开应用"},{"idx":1,"name":"检查资料"}]}]}`}
	obs, err := (&RPAAdapter{API: api, TaskID: 1}).Poll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if obs.State != StateDone {
		t.Fatalf("state = %s", obs.State)
	}
	if obs.Step != 2 || obs.Total != 2 || obs.Label != "检查资料" {
		t.Fatalf("obs = %+v", obs)
	}
}

func TestRPAFailedNamesTheFailingNode(t *testing.T) {
	api := &fakeDoer{data: `{"status":3,"taskList":[{"status":3,"failIdx":1,"nodeList":[
		{"idx":0,"name":"打开应用","description":""},
		{"idx":1,"name":"检查资料","description":"页面没有加载出来"}]}]}`}
	obs, err := (&RPAAdapter{API: api, TaskID: 1}).Poll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if obs.State != StateFailed {
		t.Fatalf("state = %s", obs.State)
	}
	if obs.Reason != "页面没有加载出来" {
		t.Fatalf("reason = %q", obs.Reason)
	}
	if obs.Step != 2 || obs.Total != 2 || obs.Label != "检查资料" {
		t.Fatalf("obs = %+v", obs)
	}
}

func TestRPAFailedFallsBackToRemarkThenScreenshot(t *testing.T) {
	api := &fakeDoer{data: `{"status":4,"taskList":[{"status":3,"remark":"任务创建失败"}]}`}
	obs, _ := (&RPAAdapter{API: api, TaskID: 1}).Poll(context.Background())
	if obs.Reason != "任务创建失败" {
		t.Fatalf("reason = %q", obs.Reason)
	}

	api = &fakeDoer{data: `{"status":4,"taskList":[{"status":3,"errorImgUrl":"https://x/y.png"}]}`}
	obs, _ = (&RPAAdapter{API: api, TaskID: 1}).Poll(context.Background())
	if obs.Reason == "" {
		t.Fatalf("expected a fallback reason when only a screenshot exists")
	}
}

func TestPhonePowerOnNeedsThePhoneToBeUp(t *testing.T) {
	api := &fakeDoer{data: `{"taskStatus":2,"taskStatusName":"开机","phoneStatus":3,"phoneStatusName":"启动中"}`}
	obs, err := (&PhoneAdapter{API: api, PhoneID: 7, OperateType: "powerOn"}).Poll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if obs.State != StateRunning {
		t.Fatalf("taskStatus=2 with the phone still booting must stay running, got %s", obs.State)
	}
	if obs.Label != "开机" {
		t.Fatalf("label = %q", obs.Label)
	}

	api = &fakeDoer{data: `{"taskStatus":2,"phoneStatus":1}`}
	obs, _ = (&PhoneAdapter{API: api, PhoneID: 7, OperateType: "powerOn"}).Poll(context.Background())
	if obs.State != StateDone {
		t.Fatalf("state = %s", obs.State)
	}
}

func TestPhonePowerOffAndFailure(t *testing.T) {
	api := &fakeDoer{data: `{"taskStatus":2,"phoneStatus":4}`}
	obs, _ := (&PhoneAdapter{API: api, PhoneID: 7, OperateType: "powerOff"}).Poll(context.Background())
	if obs.State != StateDone {
		t.Fatalf("state = %s", obs.State)
	}

	api = &fakeDoer{data: `{"taskStatus":3,"taskStatusName":"关机失败","phoneStatus":1}`}
	obs, _ = (&PhoneAdapter{API: api, PhoneID: 7, OperateType: "powerOff"}).Poll(context.Background())
	if obs.State != StateFailed || obs.Reason != "关机失败" {
		t.Fatalf("obs = %+v", obs)
	}
}

func TestPhoneQueryCarriesOperateType(t *testing.T) {
	api := &fakeDoer{data: `{}`}
	_, _ = (&PhoneAdapter{API: api, PhoneID: 7, OperateType: "restart"}).Poll(context.Background())
	if got := api.req.Query.Get("phoneId"); got != "7" {
		t.Fatalf("phoneId = %q", got)
	}
	if got := api.req.Query.Get("operateType"); got != "restart" {
		t.Fatalf("operateType = %q", got)
	}
}

func TestActionStatusDialect(t *testing.T) {
	cases := []struct {
		body string
		want State
	}{
		{`{"taskStatus":1}`, StateRunning},
		{`{"taskStatus":2}`, StateDone},
		{`{"taskStatus":3}`, StateFailed},
		{`{"taskStatus":-1}`, StateFailed},
		{`{}`, StateRunning},
	}
	for _, tc := range cases {
		api := &fakeDoer{data: tc.body}
		obs, err := (&ActionAdapter{API: api, Code: "ABC", TaskID: 9}).Poll(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if obs.State != tc.want {
			t.Errorf("%s -> %s, want %s", tc.body, obs.State, tc.want)
		}
	}
}

func TestPollPropagatesAPIError(t *testing.T) {
	apiErr := &client.APIError{Code: 10123, Msg: "access key rejected"}
	api := &fakeDoer{err: apiErr}
	_, err := (&RPAAdapter{API: api, TaskID: 1}).Poll(context.Background())
	if !errors.Is(err, apiErr) {
		t.Fatalf("err = %v", err)
	}
}
