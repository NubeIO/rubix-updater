package controller

import (
	"errors"
	"fmt"
	"github.com/NubeIO/rubix-assist/amodel"
	"github.com/NubeIO/rubix-assist/cligetter"
	"github.com/NubeIO/rubix-assist/pkg/config"
	"github.com/NubeIO/rubix-assist/pkg/interfaces"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type Snapshots struct {
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

func getBodyCreateSnapshot(c *gin.Context) (dto *interfaces.CreateSnapshot, err error) {
	err = c.ShouldBindJSON(&dto)
	return dto, err
}

func getBodyRestoreSnapshot(c *gin.Context) (dto *interfaces.RestoreSnapshot, err error) {
	err = c.ShouldBindJSON(&dto)
	return dto, err
}

func (inst *Controller) GetSnapshots(c *gin.Context) {
	host, err := inst.resolveHost(c)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	cli := cligetter.GetEdgeBiosClient(host)
	arch, err := cli.GetArch()
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	snapshots, err := inst.getSnapshots(arch.Arch)
	responseHandler(snapshots, err, c)
}

func (inst *Controller) getSnapshots(arch string) ([]Snapshots, error) {
	_path := config.Config.GetAbsSnapShotDir()
	fileInfo, err := os.Stat(_path)
	dirContent := make([]Snapshots, 0)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(_path)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			fileParts := strings.Split(file.Name(), "_")
			archParts := fileParts[len(fileParts)-1]
			archFromSnapshot := strings.Split(archParts, ".")[0]
			if archFromSnapshot == arch {
				dirContent = append(dirContent, Snapshots{
					Name:      file.Name(),
					Size:      file.Size(),
					CreatedAt: file.ModTime(),
				})
			}
		}
	} else {
		return nil, errors.New("it needs to be a directory, found a file")
	}
	return dirContent, nil
}

func (inst *Controller) DeleteSnapshot(c *gin.Context) {
	file := c.Query("file")
	if file == "" {
		responseHandler(nil, errors.New("file can not be empty"), c)
		return
	}
	err := os.Remove(path.Join(config.Config.GetAbsSnapShotDir(), file))
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	responseHandler(amodel.Message{Message: fmt.Sprintf("deleted file: %s", file)}, err, c)
}

func (inst *Controller) CreateSnapshot(c *gin.Context) {
	host, err := inst.resolveHost(c)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	body, _ := getBodyCreateSnapshot(c)
	createLog, err := inst.DB.CreateSnapshotCreateLog(&amodel.SnapshotCreateLog{UUID: "", HostUUID: host.UUID, Msg: "",
		Status: amodel.Creating, Description: body.Description, CreatedAt: time.Now()})
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	go func() {
		cli := cligetter.GetEdgeClient(host)
		snapshot, filename, err := cli.CreateSnapshot()
		if err == nil {
			err = os.WriteFile(path.Join(config.Config.GetAbsSnapShotDir(), filename), snapshot,
				os.FileMode(inst.FileMode))
		}
		createLog.Status = amodel.Created
		createLog.Msg = filename
		if err != nil {
			createLog.Status = amodel.CreateFailed
			createLog.Msg = err.Error()
		}
		_, _ = inst.DB.UpdateSnapshotCreateLog(createLog.UUID, createLog)
	}()
	responseHandler(amodel.Message{Message: "create snapshot process has submitted"}, nil, c)
}

func (inst *Controller) RestoreSnapshot(c *gin.Context) {
	body, _ := getBodyRestoreSnapshot(c)
	if body.File == "" {
		responseHandler(nil, errors.New("file can not be empty"), c)
		return
	}
	host, err := inst.resolveHost(c)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	restoreLog, err := inst.DB.CreateSnapshotRestoreLog(&amodel.SnapshotRestoreLog{UUID: "", HostUUID: host.UUID,
		Msg: "", Status: amodel.Restoring, Description: body.Description, CreatedAt: time.Now()})
	go func() {
		cli := cligetter.GetEdgeClient(host)
		reader, err := os.Open(path.Join(config.Config.GetAbsSnapShotDir(), body.File))
		if err == nil {
			err = cli.RestoreSnapshot(body.File, reader)
		}
		restoreLog.Status = amodel.Restored
		restoreLog.Msg = body.File
		if err != nil {
			restoreLog.Status = amodel.RestoreFailed
			restoreLog.Msg = err.Error()
		}
		_, _ = inst.DB.UpdateSnapshotRestoreLog(restoreLog.UUID, restoreLog)
	}()
	responseHandler(amodel.Message{Message: "restore snapshot process has submitted"}, nil, c)
}
