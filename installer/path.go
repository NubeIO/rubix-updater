package installer

import (
	"github.com/NubeIO/lib-uuid/uuid"
	"github.com/NubeIO/rubix-assist/namings"
	"os"
	"path"
)

func (inst *Installer) GetAppDataPath(appName string) string {
	dataDirName := namings.GetDataDirNameFromAppName(appName)
	return path.Join(inst.DataDir, dataDirName) // /data/rubix-wires
}

func (inst *Installer) GetAppDataDataPath(appName string) string {
	dataDirName := namings.GetDataDirNameFromAppName(appName)
	return path.Join(inst.DataDir, dataDirName, "data") // /data/rubix-wires/data
}

func (inst *Installer) GetAppDataConfigPath(appName string) string {
	dataDirName := namings.GetDataDirNameFromAppName(appName)
	return path.Join(inst.DataDir, dataDirName, "config") // /data/rubix-wires/config
}

func (inst *Installer) GetAppInstallPath(appName string) string {
	repoName := namings.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsInstallDir, repoName) // /data/rubix-service/apps/install/wires-builds
}

func (inst *Installer) GetAppInstallPathWithVersion(appName, version string) string {
	repoName := namings.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsInstallDir, repoName, version) // /data/rubix-service/apps/install/wires-builds/v0.0.1
}

func (inst *Installer) GetAppDownloadPath(appName string) string {
	repoName := namings.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsDownloadDir, repoName) // /data/rubix-service/apps/download/wires-builds
}

func (inst *Installer) GetAppDownloadPathWithVersion(appName, version string) string {
	repoName := namings.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsDownloadDir, repoName, version) // /data/rubix-service/apps/download/wires-builds/v0.0.1
}

func (inst *Installer) GetEmptyNewTmpFolder() string {
	return path.Join(inst.TmpDir, uuid.ShortUUID("tmp")) // /data/tmp/tmp_45EA34EB
}

func (inst *Installer) MakeTmpDir() error {
	return os.MkdirAll(inst.TmpDir, os.FileMode(inst.FileMode)) // /data/tmp
}

func (inst *Installer) MakeTmpDirUpload() (string, error) {
	tmpDir := inst.GetEmptyNewTmpFolder()
	err := os.MkdirAll(tmpDir, os.FileMode(inst.FileMode)) // /data/tmp/tmp_45EA34EB
	return tmpDir, err
}

func (inst *Installer) GetAppPluginDownloadPath(appName string) string {
	repoName := namings.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsDownloadDir, repoName, "plugins") // /data/rubix-service/apps/download/flow-framework/plugins
}

func (inst *Installer) GetAppPluginInstallPath(appName string) string {
	return path.Join(inst.GetAppDataDataPath(appName), "plugins") // /data/flow-framework/data/plugins
}