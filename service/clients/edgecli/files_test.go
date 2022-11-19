package edgecli

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/rubix-assist/model"
	pprint "github.com/NubeIO/rubix-assist/pkg/helpers/print"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestClient_ListFiles(t *testing.T) {
	cli := New(&Client{})
	// apps, err := cli.ListFiles("/Users/raibnod/rubix/store/apps/rubix-wires/amd64/v2.7.4/NubeIO-wires-builds-5082d47")
	// apps, err := cli.ListFilesDetails("/Users/raibnod/rubix/store/apps/rubix-wires/amd64/v2.7.4/NubeIO-wires-builds-5082d47")
	apps, err := cli.MoveFile(
		"/Users/raibnod/rubix/store/apps/rubix-wires/amd64/v2.7.4/NubeIO-wires-builds-5082d47/rubix-wires",
		"/data/tmp/test")
	if err != nil {
		fmt.Println(err)
		return
	}
	pprint.PrintJSON(apps)
}

func TestClient_RenameFile(t *testing.T) {
	cli := New(&Client{})
	dir, err := fileutils.HomeDir()
	if err != nil {
		return
	}
	old := fmt.Sprintf("%s/test/test.txt", dir)
	newName := fmt.Sprintf("%s/test/test2.txt", dir)
	fmt.Println(old, newName)
	apps, err := cli.RenameFile(old, newName)
	if err != nil {
		fmt.Println(err)
		return
	}
	pprint.PrintJSON(apps)
}

func TestClient_CopyFile(t *testing.T) {
	cli := New(&Client{})
	dir, err := fileutils.HomeDir()
	if err != nil {
		return
	}
	old := fmt.Sprintf("%s/test/test2.txt", dir)
	newName := fmt.Sprintf("%s/test/test2/test2.txt", dir)
	fmt.Println(old, newName)
	apps, err := cli.CopyFile(old, newName)
	if err != nil {
		fmt.Println(err)
		return
	}
	pprint.PrintJSON(apps)
}

func TestClient_MoveFile(t *testing.T) {
	cli := New(&Client{})
	dir, err := fileutils.HomeDir()
	if err != nil {
		return
	}
	old := fmt.Sprintf("%s/test/test2.txt", dir)
	newName := fmt.Sprintf("%s/test/test2/test2.txt", dir)
	fmt.Println(old, newName)
	apps, err := cli.MoveFile(old, newName)
	if err != nil {
		fmt.Println(err)
		return
	}
	pprint.PrintJSON(apps)
}

func TestClient_DownloadFile(t *testing.T) {
	cli := New(&Client{})
	dir, err := fileutils.HomeDir()
	if err != nil {
		return
	}
	path := fmt.Sprintf("%s/test", dir)
	fileName := "test.txt"
	dest := "/home/aidan/test/test33.txt"
	message, err := cli.DownloadFile(path, fileName, dest)
	if err != nil {
		fmt.Println(err)
		return
	}
	pprint.PrintJSON(message)
}

func TestClient_DeleteAllFiles(t *testing.T) {
	cli := New(&Client{})
	dir, err := fileutils.HomeDir()
	if err != nil {
		return
	}
	path := fmt.Sprintf("%s/test", dir)
	message, err := cli.DeleteAllFiles(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	pprint.PrintJSON(message)
}

type testYml struct {
	Auth bool `json:"auth" yaml:"auth"`
}

func TestClient_ReadFile(t *testing.T) {
	cli := New(&Client{})
	data, err := cli.ReadFile("/data/flow-framework/config/.env")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
}

func TestClient_ReadFileToYml(t *testing.T) {
	cli := New(&Client{})
	message, err := cli.ReadFile("/data/flow-framework/config/config.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(message))
	data := testYml{}
	err = yaml.Unmarshal(message, &data)
	fmt.Println(err)
	pprint.PrintJSON(data)
}

type testJson struct {
	ImageVersion string `json:"image_version"`
	Arch         string `json:"arch"`
	Product      string `json:"product"`
}

func TestClient_ReadFileToJson(t *testing.T) {
	cli := New(&Client{})
	message, err := cli.ReadFile("/data/product.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	data := testJson{}
	err = json.Unmarshal(message, &data)
	fmt.Println(err)
	fmt.Println(data)
}

func TestClient_WriteFileJson(t *testing.T) {
	data := testJson{
		ImageVersion: "v1.2.3.4",
		Arch:         "amd64",
		Product:      "Server",
	}
	cli := New(&Client{})
	body := &model.WriteFile{
		FilePath: "/data/product.json",
		Body:     data,
	}
	message, err := cli.WriteFileJson(body)
	if err != nil {
		fmt.Println(err)
		return
	}
	pprint.PrintJSON(message)
}

func TestClient_WriteFileYml(t *testing.T) {
	data := testYml{
		Auth: false,
	}
	cli := New(&Client{})
	body := &model.WriteFile{
		FilePath: "/data/flow-framework/config/config.yml",
		Body:     data,
	}
	message, err := cli.WriteFileYml(body)
	if err != nil {
		fmt.Println(err)
		return
	}
	pprint.PrintJSON(message)
}

func TestClient_WriteFile(t *testing.T) {
	data := `
PORT=1313
SECRET_KEY=__SECRET_KEY__
`
	cli := New(&Client{})
	body := &model.WriteFile{
		FilePath:     "/data/rubix-wires/config/.env",
		BodyAsString: data,
	}
	message, err := cli.WriteString(body)
	if err != nil {
		fmt.Println(err)
		return
	}
	pprint.PrintJSON(message)
}
