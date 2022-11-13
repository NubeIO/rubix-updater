package systemctl

import (
	"fmt"
	"github.com/NubeIO/lib-rubix-installer/installer"
	"testing"
)

func TestStore_generateServiceFile(t *testing.T) {
	tmpDir, absoluteServiceFileName, err := GenerateServiceFile(&ServiceFile{
		Name:                        "rubix-edge",
		Version:                     "v0.6.0",
		ExecStart:                   "app -p 1661 -r /data -a rubix-edge -d data -c config --prod server",
		AttachWorkingDirOnExecStart: true,
	}, installer.New(&installer.App{}))
	fmt.Println(tmpDir, absoluteServiceFileName, err)
	if err != nil {
		return
	}
}
