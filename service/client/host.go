package client

import (
	"fmt"
	"github.com/NubeIO/rubix-assist/pkg/model"
)

func (inst *Client) GetHosts() (data []model.Host, response *Response) {
	path := fmt.Sprintf(Paths.Hosts.Path)
	response = &Response{}
	resp, err := inst.Rest.R().
		SetResult(&[]model.Host{}).
		Get(path)
	return *resp.Result().(*[]model.Host), response.buildResponse(resp, err)
}

func (inst *Client) GetHost(uuid string) (data *model.Host, response *Response) {
	path := fmt.Sprintf("%s/%s", Paths.Hosts.Path, uuid)
	response = &Response{}
	resp, err := inst.Rest.R().
		SetResult(&model.Host{}).
		Get(path)
	return resp.Result().(*model.Host), response.buildResponse(resp, err)
}

func (inst *Client) AddHost(body *model.Host) (data *model.Host, response *Response) {
	path := fmt.Sprintf(Paths.Hosts.Path)
	response = &Response{}
	resp, err := inst.Rest.R().
		SetBody(body).
		SetResult(&model.Host{}).
		Post(path)
	return resp.Result().(*model.Host), response.buildResponse(resp, err)
}

func (inst *Client) UpdateHost(uuid string, body *model.Host) (data *model.Host, response *Response) {
	path := fmt.Sprintf("%s/%s", Paths.Hosts.Path, uuid)
	response = &Response{}
	resp, err := inst.Rest.R().
		SetBody(body).
		SetResult(&model.Host{}).
		Patch(path)
	return resp.Result().(*model.Host), response.buildResponse(resp, err)
}

func (inst *Client) DeleteHost(uuid string) (response *Response) {
	path := fmt.Sprintf("%s/%s", Paths.Hosts.Path, uuid)
	response = &Response{}
	resp, err := inst.Rest.R().
		Delete(path)
	return response.buildResponse(resp, err)
}
func (inst *Client) GetHostSchema() (data *model.HostSchema, response *Response) {
	path := fmt.Sprintf("%s/%s", Paths.Hosts.Path, "schema")
	response = &Response{}
	resp, err := inst.Rest.R().
		SetResult(&model.HostSchema{}).
		Get(path)
	return resp.Result().(*model.HostSchema), response.buildResponse(resp, err)
}
