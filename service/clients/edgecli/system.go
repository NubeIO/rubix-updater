package edgecli

import (
	"fmt"
	"github.com/NubeIO/rubix-assist/model"
	"github.com/NubeIO/rubix-assist/service/clients/helpers/nresty"
	"github.com/NubeIO/rubix-registry-go/rubixregistry"
)

// EdgeProductInfo get edge product info
func (inst *Client) EdgeProductInfo() (*model.Product, error) {
	url := fmt.Sprintf("/api/system/product")
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&model.Product{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Product), nil
}

// EdgeGetDeviceInfo get edge device info
func (inst *Client) EdgeGetDeviceInfo() (*rubixregistry.DeviceInfo, error) {
	url := fmt.Sprintf("/api/system/device")
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&rubixregistry.DeviceInfo{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*rubixregistry.DeviceInfo), nil
}

// Ping ping an edge device
func (inst *Client) Ping() (*model.Message, error) {
	url := fmt.Sprintf("/api/system/ping")
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&model.Message{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Message), nil
}
