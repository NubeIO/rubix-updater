package controller

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/rest/v1/rest"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/system/dirs"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/utilities/git"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (base *Controller) InstallBios(ctx *gin.Context) {
	//mk download dir and or clear
	//download and unzip build into bios
	//unzip bios
	//run install
	fmt.Println(rest.DELETE)
	body, err := getGitBody(ctx)
	token := resolveHeaderGitToken(ctx)
	host, _, err := base.resolveHost(ctx)
	if err != nil {
		reposeHandler(nil, err, ctx)
		return
	}
	path := fmt.Sprintf("/home/%s/rubix-bios-install", host.Username)
	_host, _ := base.hostCopy(host)
	_dirs := dirs.Dirs{
		Host:          _host,
		Name:          path,
		CheckIfExists: true,
		IfExistsClear: true,
	}

	//MAKE DIR if not existing and also clear dir
	log.Println("mk dir ", "try and make dir")
	_, err = _dirs.MKDir()
	if err != nil {
		fmt.Println("mk dir", "mk dir fail")
		reposeHandler(nil, err, ctx)
		return
	}
	log.Println("mk dir ", "mk dir pass")

	g := body
	g.Token = token
	g.DownloadPath = path
	_dirs.Host.CommandOpts.CMD = g.BuildCURL(git.CurlReleaseDownload)
	log.Println("download ", _dirs.Host.CommandOpts.CMD)
	//DOWNLOAD BUILD
	log.Println("download ", "try and download bios")
	_, download, err := _dirs.Host.RunCommand()
	if err != nil {
		reposeHandler(nil, err, ctx)
		return
	}
	log.Println("download ", download)

	//UNZIP BUILD
	_dirs.Host.CommandOpts.CMD = "unzip -o " + g.DownloadPath + "/" + g.FolderName + " -d " + g.DownloadPath
	log.Println("UNZIP BUILD ", _dirs.Host.CommandOpts.CMD)
	_, unzip, err := _dirs.Host.RunCommand()
	if err != nil {
		log.Errorln("UNZIP BUILD FAIL ", _dirs.Host.CommandOpts.CMD)
		reposeHandler(nil, err, ctx)
		return
	}
	log.Println("unzip ", unzip)

	_dirs.Host.CommandOpts.CMD = "rm /data/rubix-service/config/app.json"
	//rm /data/rubix-service/config/apps.json
	_, deleteConfigFile, _ := _dirs.Host.RunCommand()
	log.Println("deleteConfigFile ", deleteConfigFile, _dirs.Host.CommandOpts.CMD)

	_dirs.Host.CommandOpts.CMD = "rm /data/rubix-service/config/apps.json"
	_, deleteConfigFile, _ = _dirs.Host.RunCommand()
	log.Println("deleteConfigFile ", deleteConfigFile, _dirs.Host.CommandOpts.CMD)

	//INSTALL BUILD
	_dirs.Host.CommandOpts.CMD = fmt.Sprintf("cd %s; sudo ./rubix-bios -p 1615 -g /data/rubix-bios -d data -c config -a apps --prod --install --auth --device-type %s --token %s", g.DownloadPath, g.Target, token)
	log.Println("cmd ", _dirs.Host.CommandOpts.CMD)
	_, install, err := _dirs.Host.RunCommand()
	if err != nil {
		reposeHandler(nil, err, ctx)
		return
	}
	log.Println("install ", install)
	reposeHandler("installed", err, ctx)
	return
}

//func (base *Controller) RubixServiceUpdate(ctx *gin.Context) {
//	host, _, _ := base.resolveHost(ctx)
//	//get token
//	r := &nrest.ReqType{
//		BaseUri: host.IP,
//		Port:    host.BiosPort,
//		Path:    "/api/users/login",
//		Method:  nrest.POST,
//	}
//	opt := &nrest.ReqOpt{
//		Timeout:          2 * time.Second,
//		RetryCount:       2,
//		RetryWaitTime:    2 * time.Second,
//		RetryMaxWaitTime: 0,
//		Json:             map[string]interface{}{"username": host.RubixUsername, "password": host.RubixPassword},
//	}
//
//	req, status, err := nrest.DoHTTPReq(r, opt)
//	res := new(rubix.TokenResponse)
//	err = req.ToInterface(&res)
//	if err != nil {
//		reposeHandler(nil, errors.New("failed to get bios token"), ctx)
//	}
//	token := res.AccessToken
//	log.Info("GET bios token status:", status, "GET bios token string:", token)
//
//	opt = &nrest.ReqOpt{
//		Timeout:          500 * time.Second,
//		RetryCount:       0,
//		RetryWaitTime:    0 * time.Second,
//		RetryMaxWaitTime: 0,
//		Headers:          map[string]interface{}{"Authorization": token},
//	}
//
//	r = &nrest.ReqType{
//		BaseUri: host.IP,
//		Port:    host.BiosPort,
//		Path:    "/api/service/upgrade_and_check",
//		Method:  nrest.PUT,
//	}
//
//	req, status, err = nrest.DoHTTPReq(r, opt)
//	log.Info("bios: get rubix-service installed version status:", status)
//	log.Info("bios: get rubix-service installed version status:", req.AsString())
//	if err != nil {
//		reposeHandler(req.AsJsonNoErr(), err, ctx)
//		return
//	} else {
//		reposeHandler(req.AsJsonNoErr(), err, ctx)
//	}
//}

//func (base *Controller) RubixServiceCheck(ctx *gin.Context) {
//	host, _, _ := base.resolveHost(ctx)
//	//get token
//	r := &nrest.ReqType{
//		BaseUri: host.IP,
//		Port:    host.BiosPort,
//		Path:    "/api/users/login",
//		Method:  nrest.POST,
//	}
//	opt := &nrest.ReqOpt{
//		Timeout:          2 * time.Second,
//		RetryCount:       2,
//		RetryWaitTime:    2 * time.Second,
//		RetryMaxWaitTime: 0,
//		Json:             map[string]interface{}{"username": host.RubixUsername, "password": host.RubixPassword},
//	}
//
//	req, status, err := nrest.DoHTTPReq(r, opt)
//	res := new(rubix.TokenResponse)
//	err = req.ToInterface(&res)
//	if err != nil {
//		reposeHandler(nil, errors.New("failed to get bios token"), ctx)
//	}
//	token := res.AccessToken
//	log.Info("GET bios token status:", status, "GET bios token string:", token)
//
//	opt = &nrest.ReqOpt{
//		Timeout:          500 * time.Second,
//		RetryCount:       0,
//		RetryWaitTime:    0 * time.Second,
//		RetryMaxWaitTime: 0,
//		Headers:          map[string]interface{}{"Authorization": token},
//	}
//
//	r = &nrest.ReqType{
//		BaseUri: host.IP,
//		Port:    host.BiosPort,
//		Path:    "/api/service/update_check",
//	}
//
//	req, status, err = nrest.DoHTTPReq(r, opt)
//	log.Info("bios: get rubix-service installed version status:", status)
//	if err != nil {
//		reposeHandler(nil, err, ctx)
//		return
//	} else {
//		reposeHandler(req.AsJsonNoErr(), err, ctx)
//	}
//}