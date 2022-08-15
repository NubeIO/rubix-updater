package appstore

import (
	"errors"
	"fmt"
	fileutils "github.com/NubeIO/lib-dirs/dirs"
	"github.com/NubeIO/lib-rubix-installer/installer"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/sergeymakinen/go-systemdconf/v2"
	"github.com/sergeymakinen/go-systemdconf/v2/unit"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type EdgeApp struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	Product           string `json:"product"`
	Arch              string `json:"arch"`
	ServiceDependency string `json:"service_dependency"` // nodejs
}

// AddUploadEdgeApp
// upload the build
func (inst *Store) AddUploadEdgeApp(hostUUID, hostName string, app *EdgeApp) (*installer.AppResponse, error) {
	appName := app.Name
	version := app.Version
	archType := app.Arch
	productType := app.Product
	if appName == "" {
		return nil, errors.New("upload app to edge app name can not be empty")
	}
	if version == "" {
		return nil, errors.New("upload app to edge  app version can not be empty")
	}
	if productType == "" {
		return nil, errors.New("upload app to edge  product type can not be empty, try RubixCompute, RubixComputeIO, RubixCompute5, Server, Edge28, Nuc")
	}
	if archType == "" {
		return nil, errors.New("upload app to edge arch type can not be empty, try armv7 amd64")
	}
	var fileName string
	path := inst.getAppStorePathAndVersion(appName, version)
	fileNames, err := inst.App.GetBuildZipNames(path)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to find zip build err:%s", err.Error()))
		return nil, err
	}
	if len(fileNames) > 0 {
		fileName = fileNames[0].ZipName
	} else {
		err := errors.New(fmt.Sprintf("no zip builds found in path:%s", path))
		return nil, err
	}
	fileAndPath := filePath(fmt.Sprintf("%s/%s", path, fileName))
	reader, err := os.Open(fileAndPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error open build for app:%s fileName:%s  err:%s", appName, fileName, err.Error()))
	}
	client, err := inst.getClient(hostUUID, hostName)
	if err != nil {
		return nil, err
	}
	return client.UploadApp(appName, version, productType, archType, fileName, reader)
}

func (inst *Store) setServiceName(appName string) string {
	return fmt.Sprintf("nubeio-%s", appName)
}

func (inst *Store) setServiceFileName(appName string) string {
	return fmt.Sprintf("nubeio-%s.service", appName)
}

func (inst *Store) setServiceWorkingDir(appName, appVersion string) string {
	return inst.App.GetAppInstallPathAndVersion(appName, appVersion)
}

func (inst *Store) setServiceExecStart(appName, appVersion, AppSpecficExecStart string) string {
	workingDir := inst.App.GetAppInstallPathAndVersion(appName, appVersion)
	return fmt.Sprintf("%s/%s", workingDir, AppSpecficExecStart)
}

func (inst *Store) checkServiceExecStart(service, appName, appVersion string) error {
	if strings.Contains(service, inst.App.GetAppInstallPathAndVersion(appName, appVersion)) {
		return nil
	}
	return errors.New(fmt.Sprintf("ExecStart command is not matching appName:%sappName & appVersion:%s", appName, appVersion))
}

type ServiceFile struct {
	Name                    string   `json:"name"`
	Version                 string   `json:"version"`
	ServiceDependency       string   `json:"service_dependency"` // nodejs
	ServiceDescription      string   `json:"service_description"`
	RunAsUser               string   `json:"run_as_user"`
	ServiceWorkingDirectory string   `json:"service_working_directory"` // /data/rubix-service/apps/install/flow-framework/v0.6.1/
	AppSpecficExecStart     string   `json:"app_specfic_exec_start"`    // WORKING-DIR/app -p 1660 -g /data/flow-framework -d data -prod
	CustomServiceExecStart  string   `json:"custom_service_exec_start"` // npm run prod:start --prod --datadir /data/rubix-wires/data --envFile /data/rubix-wires/config/.env
	EnvironmentVars         []string `json:"environment_vars"`          // Environment="g=/data/bacnet-server-c"
}

// InstallEdgeService this assumes that the service file and app already exists on the edge device
func (inst *Store) InstallEdgeService(hostUUID, hostName string, body *installer.Install) (*installer.InstallResp, error) {
	client, err := inst.getClient(hostUUID, hostName)
	if err != nil {
		return nil, err
	}
	return client.InstallService(body)
}

func (inst *Store) GenerateUploadEdgeService(hostUUID, hostName string, app *ServiceFile) (*installer.UploadResponse, error) {
	resp, err := inst.generateUploadEdgeService(hostUUID, hostName, app)
	if err != nil {
		log.Errorf("generate service file hostUUID:%s, hostName:%s appName:%s", hostUUID, hostName, app.Name)
		log.Errorf("generate service file main err:%s", err.Error())
	}
	return resp, err
}

func (inst *Store) generateServiceFile(app *ServiceFile) (tmpDir, serviceFile, fileAndPath string, err error) {
	tmpFilePath, err := inst.App.MakeTmpDirUpload()
	if err != nil {
		return "", "", "", err
	}
	if app.Name == "" {
		return "", "", "", errors.New("app name can not be empty, try flow-framework")
	}
	serviceName := inst.setServiceName(app.Name)
	serviceFileName := inst.setServiceFileName(app.Name)
	appVersion := app.Version
	if appVersion == "" {
		return "", "", "", errors.New("app version can not be empty, try v0.6.0")
	}
	if err = checkVersion(appVersion); err != nil {
		return "", "", "", err
	}
	appName := app.Name
	if appVersion == "" {
		return "", "", "", errors.New("app build name can not be empty, try wires-builds")
	}
	workingDirectory := app.ServiceWorkingDirectory
	if workingDirectory == "" {
		workingDirectory = inst.setServiceWorkingDir(appName, appVersion)
	}
	log.Infof("generate service working dir: %s", workingDirectory)
	user := app.RunAsUser
	if user == "" {
		user = "root"
	}
	execCmd := app.AppSpecficExecStart
	if app.CustomServiceExecStart != "" { // example use would be in wires
		execCmd = app.CustomServiceExecStart
	} else {
		if execCmd == "" {
			return "", "", "", errors.New("app service ExecStart cant not be empty")
		}
		execCmd = inst.setServiceExecStart(app.Name, appVersion, execCmd)
		if err := inst.checkServiceExecStart(execCmd, appName, appVersion); err != nil {
			return "", "", "", err
		}
	}

	log.Infof("generate service execCmd: %s", execCmd)
	description := app.ServiceDescription
	if description == "" {
		description = fmt.Sprintf("NubeIO %s", app.Name)
	}

	var env systemdconf.Value
	for _, s := range app.EnvironmentVars {
		env = append(env, s)
	}
	service := unit.ServiceFile{
		Unit: unit.UnitSection{ // [Unit]
			Description: systemdconf.Value{description},
			After:       systemdconf.Value{"network.target"},
		},
		Service: unit.ServiceSection{ // [Service]
			ExecStartPre: nil,
			Type:         systemdconf.Value{"simple"},
			ExecOptions: unit.ExecOptions{
				User:             systemdconf.Value{user},
				WorkingDirectory: systemdconf.Value{workingDirectory},
				Environment:      env,
				StandardOutput:   systemdconf.Value{"syslog"},
				StandardError:    systemdconf.Value{"syslog"},
				SyslogIdentifier: systemdconf.Value{app.Name},
			},
			ExecStart: systemdconf.Value{
				execCmd,
			},
			Restart:    systemdconf.Value{"always"},
			RestartSec: systemdconf.Value{"10"},
		},
		Install: unit.InstallSection{ // [Install]
			WantedBy: systemdconf.Value{"multi-user.target"},
		},
	}

	b, _ := systemdconf.Marshal(service)
	fmt.Println(serviceFileName)
	fmt.Println(string(b))

	servicePath := fmt.Sprintf("%s/%s", tmpFilePath, serviceFileName)
	file := fileutils.New()
	err = file.WriteFile(servicePath, string(b), os.FileMode(FilePerm))
	if err != nil {
		log.Errorf("write service file error %s", err.Error())
	}
	log.Infof("generate service file name:%s", serviceName)
	//fileAndPath := filePath(fmt.Sprintf("%s/%s", tmpFilePath, serviceFileName))
	log.Infof("generate service file path:%s", servicePath)
	return tmpFilePath, serviceFileName, servicePath, nil

}

// GenerateUploadEdgeService this will generate and upload the service file to the edge device
func (inst *Store) generateUploadEdgeService(hostUUID, hostName string, app *ServiceFile) (*installer.UploadResponse, error) {
	tmpDir, serviceFile, fileAndPath, err := inst.generateServiceFile(app)
	if err != nil {
		return nil, err
	}
	reader, err := os.Open(fileAndPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error open service-file:%s err:%s", serviceFile, err.Error()))
	}

	client, err := inst.getClient(hostUUID, hostName)
	if err != nil {
		return nil, err
	}
	err = fileutils.New().RmRF(tmpDir)
	if err != nil {
		log.Errorf("assist: delete tmp dir after generating service file%s", fileAndPath)
	}
	return client.UploadServiceFile(app.Name, app.Version, serviceFile, reader)

}

//// GenerateUploadEdgeService this will generate and upload the service file to the edge device
//func (inst *Store) generateUploadEdgeService(hostUUID, hostName string, app *ServiceFile) (*installer.UploadResponse, error) {
//	tmpFilePath, err := inst.App.MakeTmpDirUpload()
//	if err != nil {
//		return nil, err
//	}
//	if app.Name == "" {
//		return nil, errors.New("app name can not be empty, try flow-framework")
//	}
//	serviceName := inst.setServiceName(app.Name)
//	serviceFileName := inst.setServiceFileName(app.Name)
//	appVersion := app.Version
//	if appVersion == "" {
//		return nil, errors.New("app version can not be empty, try v0.6.0")
//	}
//	if err = checkVersion(appVersion); err != nil {
//		return nil, err
//	}
//	appName := app.Name
//	if appVersion == "" {
//		return nil, errors.New("app build name can not be empty, try wires-builds")
//	}
//	workingDirectory := app.ServiceWorkingDirectory
//	if workingDirectory == "" {
//		workingDirectory = inst.setServiceWorkingDir(appName, appVersion)
//	}
//	log.Infof("generate service working dir: %s", workingDirectory)
//	user := app.RunAsUser
//	if user == "" {
//		user = "root"
//	}
//	execCmd := app.AppSpecficExecStart
//	if app.CustomServiceExecStart != "" { // example use would be in wires
//		execCmd = app.CustomServiceExecStart
//	} else {
//		if execCmd == "" {
//			return nil, errors.New("app service ExecStart cant not be empty")
//		}
//		execCmd = inst.setServiceExecStart(app.Name, appVersion, execCmd)
//		if err := inst.checkServiceExecStart(execCmd, appName, appVersion); err != nil {
//			return nil, err
//		}
//	}
//
//	log.Infof("generate service execCmd: %s", execCmd)
//	description := app.ServiceDescription
//	bld := &builder.SystemDBuilder{
//		ServiceName:      serviceName,
//		Description:      description,
//		User:             user,
//		WorkingDirectory: workingDirectory,
//		ExecStart:        execCmd,
//		SyslogIdentifier: serviceName,
//		WriteFile: builder.WriteFile{
//			Write:    true,
//			FileName: serviceName,
//			Path:     tmpFilePath,
//		},
//	}
//	err = bld.Build(os.FileMode(inst.Perm))
//	if err != nil {
//		log.Errorf("generate service file name:%s, err:%s", serviceName, err.Error())
//		return nil, err
//	}
//
//	log.Infof("generate service file name:%s", serviceName)
//
//	fileAndPath := filePath(fmt.Sprintf("%s/%s", tmpFilePath, serviceFileName))
//	log.Infof("generate service file path:%s", fileAndPath)
//	reader, err := os.Open(fileAndPath)
//	if err != nil {
//		return nil, errors.New(fmt.Sprintf("error open file:%s err:%s", fileAndPath, err.Error()))
//	}
//
//	client, err := inst.getClient(hostUUID, hostName)
//	if err != nil {
//		return nil, err
//	}
//
//	return client.UploadServiceFile(app.Name, appVersion, serviceFileName, reader)
//}

func (inst *Store) EdgeUnInstallApp(hostUUID, hostName, appName string, deleteApp bool) (*installer.RemoveRes, error) {
	client, err := inst.getClient(hostUUID, hostName)
	if err != nil {
		return nil, err
	}
	return client.EdgeUnInstallApp(appName, deleteApp)
}

func (inst *Store) EdgeListApps(hostUUID, hostName string) ([]installer.Apps, error) {
	client, err := inst.getClient(hostUUID, hostName)
	if err != nil {
		return nil, err
	}
	return client.ListApps()
}

func (inst *Store) EdgeListAppsAndService(hostUUID, hostName string) ([]installer.InstalledServices, error) {
	client, err := inst.getClient(hostUUID, hostName)
	if err != nil {
		return nil, err
	}
	return client.ListAppsAndService()
}

func (inst *Store) EdgeListNubeServices(hostUUID, hostName string) ([]installer.InstalledServices, error) {
	client, err := inst.getClient(hostUUID, hostName)
	if err != nil {
		return nil, err
	}
	return client.ListAppsAndService()
}

func (inst *Store) EdgeCtlAction(hostUUID, hostName string, body *installer.CtlBody) (*systemctl.SystemResponse, error) {
	client, err := inst.getClient(hostUUID, hostName)
	if err != nil {
		return nil, err
	}
	return client.EdgeCtlAction(body)
}

func (inst *Store) EdgeCtlStatus(hostUUID, hostName string, body *installer.CtlBody) (*systemctl.SystemState, error) {
	client, err := inst.getClient(hostUUID, hostName)
	if err != nil {
		return nil, err
	}
	return client.EdgeCtlStatus(body)
}

func (inst *Store) EdgeServiceMassAction(hostUUID, hostName string, body *installer.CtlBody) ([]systemctl.MassSystemResponse, error) {
	client, err := inst.getClient(hostUUID, hostName)
	if err != nil {
		return nil, err
	}
	return client.EdgeServiceMassAction(body)
}

func (inst *Store) EdgeServiceMassStatus(hostUUID, hostName string, body *installer.CtlBody) ([]systemctl.SystemState, error) {
	client, err := inst.getClient(hostUUID, hostName)
	if err != nil {
		return nil, err
	}
	return client.EdgeServiceMassStatus(body)
}

func (inst *Store) EdgeProductInfo(hostUUID, hostName string) (*installer.Product, error) {
	client, err := inst.getClient(hostUUID, hostName)
	if err != nil {
		return nil, err
	}
	return client.EdgeProductInfo()
}
