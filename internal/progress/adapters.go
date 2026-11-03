package progress

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/haidaisiganbei/pmock-cli/internal/client"
)

// Doer is the slice of client.Client the adapters need. Tests substitute a
// scripted implementation.
type Doer interface {
	Do(ctx context.Context, req client.Request) (*client.Envelope, error)
}

// Adapter turns one PMOCK capability into observations.
//
// Adding a capability means adding an adapter, never a new state: the four
// states are the contract, the dialects are an implementation detail.
type Adapter interface {
	// Kind is the short capability name carried in every event ("rpa").
	Kind() string
	// Poll reads the current position of the task.
	Poll(ctx context.Context) (Observation, error)
}

type rpaNode struct {
	Idx         int    `json:"idx"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type rpaDetail struct {
	Status   *int      `json:"status"`
	FailIdx  *int      `json:"failIdx"`
	ErrorImg string    `json:"errorImgUrl"`
	Remark   string    `json:"remark"`
	Nodes    []rpaNode `json:"nodeList"`
}

type rpaInfo struct {
	Status   *int        `json:"status"`
	TaskList []rpaDetail `json:"taskList"`
}

// RPAAdapter follows an RPA task (0 待执行, 1 进行中, 2 成功, 3 部分失败, 4 失败).
//
// RPA is the only capability that can report named business steps, through the
// node list. When the node list is absent the adapter degrades to a bare
// status, and Copy degrades with it.
type RPAAdapter struct {
	API    Doer
	TaskID int64
}

// Kind implements Adapter.
func (a *RPAAdapter) Kind() string { return "rpa" }

// Poll implements Adapter.
func (a *RPAAdapter) Poll(ctx context.Context) (Observation, error) {
	var info rpaInfo
	req := client.Request{
		Method: "GET",
		Path:   fmt.Sprintf("/open/api/rpa/task/getInfo/%d", a.TaskID),
	}
	if err := doJSON(ctx, a.API, req, &info); err != nil {
		return Observation{}, err
	}

	obs := Observation{Observed: map[string]any{}}
	if info.Status != nil {
		obs.Observed["status"] = *info.Status
	}
	if len(info.TaskList) > 0 && info.TaskList[0].Status != nil {
		obs.Observed["detailStatus"] = *info.TaskList[0].Status
	}
	if info.Status == nil {
		obs.State = StateRunning
		return obs, nil
	}

	var detail *rpaDetail
	if len(info.TaskList) > 0 {
		detail = &info.TaskList[0]
	}

	switch *info.Status {
	case 0:
		// Accepted but not started: to the customer this is still "wait".
		obs.State = StateRunning
		obs.Label = "排队中"
	case 1:
		obs.State = StateRunning
	case 2:
		obs.State = StateDone
		if detail != nil && len(detail.Nodes) > 0 {
			obs.Total = len(detail.Nodes)
			obs.Step = obs.Total
			obs.Label = detail.Nodes[len(detail.Nodes)-1].Name
		}
	case 3, 4:
		obs.State = StateFailed
		if detail != nil {
			obs.Reason, obs.Step, obs.Total, obs.Label = rpaFailure(detail)
		}
	default:
		obs.State = StateRunning
	}
	return obs, nil
}

func rpaFailure(detail *rpaDetail) (reason string, step, total int, label string) {
	total = len(detail.Nodes)
	if detail.FailIdx != nil && total > 0 {
		idx := *detail.FailIdx
		step = idx + 1
		for _, n := range detail.Nodes {
			if n.Idx == idx {
				label, reason = n.Name, n.Description
				break
			}
		}
		if label == "" && idx >= 0 && idx < total {
			label, reason = detail.Nodes[idx].Name, detail.Nodes[idx].Description
		}
	}
	if reason == "" {
		reason = detail.Remark
	}
	if reason == "" && detail.ErrorImg != "" {
		reason = "流程执行时出错了，已留下当时的截图"
	}
	return reason, step, total, label
}

// PhoneAdapter follows a phone lifecycle operation: start, stop, restart, reset
// or one-click renew. The endpoint needs the phone id and the operate type.
//
// The API answers with two axes (taskStatus and phoneStatus); an operation is
// only settled when the phone itself reached the expected state. Until then the
// customer keeps seeing "running", not a premature success.
type PhoneAdapter struct {
	API         Doer
	PhoneID     int64
	OperateType string
}

// Kind implements Adapter.
func (a *PhoneAdapter) Kind() string { return "phone" }

// Poll implements Adapter.
func (a *PhoneAdapter) Poll(ctx context.Context) (Observation, error) {
	var st struct {
		TaskStatus      *int   `json:"taskStatus"`
		TaskStatusName  string `json:"taskStatusName"`
		PhoneStatus     *int   `json:"phoneStatus"`
		PhoneStatusName string `json:"phoneStatusName"`
	}
	req := client.Request{
		Method: "GET",
		Path:   "/open/api/phone/task/status",
		Query: url.Values{
			"phoneId":     {strconv.FormatInt(a.PhoneID, 10)},
			"operateType": {a.OperateType},
		},
	}
	if err := doJSON(ctx, a.API, req, &st); err != nil {
		return Observation{}, err
	}

	obs := Observation{Observed: map[string]any{"operateType": a.OperateType}}
	if st.TaskStatus != nil {
		obs.Observed["taskStatus"] = *st.TaskStatus
	}
	if st.PhoneStatus != nil {
		obs.Observed["phoneStatus"] = *st.PhoneStatus
	}

	label := st.TaskStatusName
	if label == "" {
		label = st.PhoneStatusName
	}
	switch {
	case st.TaskStatus == nil:
		obs.State = StateRunning
	case *st.TaskStatus == 2:
		if phoneSettled(a.OperateType, st.PhoneStatus) {
			obs.State = StateDone
		} else {
			obs.State = StateRunning
			obs.Label = label
		}
	case *st.TaskStatus == 3:
		obs.State = StateFailed
		obs.Reason = label
	default:
		obs.State = StateRunning
		obs.Label = label
	}
	return obs, nil
}

func phoneSettled(operateType string, phoneStatus *int) bool {
	switch operateType {
	case "powerOn":
		return phoneStatus != nil && *phoneStatus == 1
	case "powerOff":
		return phoneStatus != nil && *phoneStatus == 4
	default:
		// restart / reset / renew settle on the task status alone.
		return true
	}
}

// ActionAdapter follows one UI action (click, swipe, input, ...).
//
// The published enum is -1/0/1/2 but the endpoint returns 1/2/3 in practice,
// the same mismatch waitAction already compensates for. Actions are usually
// sub-second, so this adapter exists for the occasional long press chain rather
// than for everyday narration.
type ActionAdapter struct {
	API    Doer
	Code   string
	TaskID int64
}

// Kind implements Adapter.
func (a *ActionAdapter) Kind() string { return "action" }

// Poll implements Adapter.
func (a *ActionAdapter) Poll(ctx context.Context) (Observation, error) {
	var st struct {
		TaskStatus *int `json:"taskStatus"`
	}
	req := client.Request{
		Method: "POST",
		Path:   "/open/api/phone/action/task/status",
		Body:   map[string]any{"code": a.Code, "taskId": a.TaskID},
	}
	if err := doJSON(ctx, a.API, req, &st); err != nil {
		return Observation{}, err
	}
	obs := Observation{Observed: map[string]any{}}
	if st.TaskStatus == nil {
		obs.State = StateRunning
		return obs, nil
	}
	obs.Observed["taskStatus"] = *st.TaskStatus
	switch *st.TaskStatus {
	case 2:
		obs.State = StateDone
	case 3, -1:
		obs.State = StateFailed
	default:
		obs.State = StateRunning
	}
	return obs, nil
}

func doJSON(ctx context.Context, api Doer, req client.Request, out any) error {
	env, err := api.Do(ctx, req)
	if err != nil {
		return err
	}
	if env == nil || len(env.Data) == 0 || string(env.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(env.Data, out); err != nil {
		return fmt.Errorf("failed to parse %s response: %w", req.Path, err)
	}
	return nil
}
