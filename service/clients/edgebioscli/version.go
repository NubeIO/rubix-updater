package edgebioscli

import (
	"errors"
	"fmt"
	"github.com/NubeIO/rubix-assist/service/clients/edgebioscli/ebmodel"
	"github.com/NubeIO/rubix-assist/service/clients/helpers/nresty"
)

func (inst *BiosClient) GetRubixEdgeVersion() (*ebmodel.Version, error) {
	installLocation := fmt.Sprintf("/data/rubix-service/apps/install/%s", rubixEdgeName)
	url := fmt.Sprintf("/api/files/list?path=%s", installLocation)
	resp, err := nresty.FormatRestyResponse(inst.Rest.R().
		SetResult(&[]string{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	versions := resp.Result().(*[]string)
	if versions != nil && len(*versions) > 0 {
		return &ebmodel.Version{Version: (*versions)[0]}, nil
	}
	return nil, errors.New("doesn't found the installation file")
}
