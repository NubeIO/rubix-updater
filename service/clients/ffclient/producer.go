package ffclient

import (
	"fmt"
	"github.com/NubeIO/lib-uuid/uuid"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-assist/service/clients/ffclient/nresty"
)

// AddProducer an object
func (inst *FlowClient) AddProducer(body model.Producer) (*model.Producer, error) {
	name := uuid.ShortUUID()
	name = fmt.Sprintf("sub_name_%s", name)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Producer{}).
		SetBody(body).
		Post("/api/producers"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Producer), nil
}

func (inst *FlowClient) GetProducers(streamUUID *string) (*[]model.Producer, error) {
	req := inst.client.R().
		SetResult(&[]model.Producer{})
	if streamUUID != nil {
		req.SetQueryParam("stream_uuid", *streamUUID)
	}
	resp, err := nresty.FormatRestyResponse(req.Get("/api/producers"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*[]model.Producer), nil
}

func (inst *FlowClient) GetProducer(uuid string) (*model.Producer, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Producer{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/producers/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Producer), nil
}

// EditProducer edit an object
func (inst *FlowClient) EditProducer(uuid string, body model.Producer) (*model.Producer, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Producer{}).
		SetBody(body).
		SetPathParams(map[string]string{"uuid": uuid}).
		Patch("/api/producers/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Producer), nil
}