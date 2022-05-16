package controller

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/networking/ssh"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/utilities/git"
	"github.com/NubeIO/nubeio-rubix-lib-rest-go/pkg/rest"

	dbase "github.com/NubeIO/rubix-assist/database"
	"github.com/NubeIO/rubix-assist/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/melbahja/goph"
	"gopkg.in/olahol/melody.v1"
)

type Controller struct {
	//DB  *gorm.DB
	SSH  *goph.Client
	WS   *melody.Melody //web socket
	DB   *dbase.DB
	Rest *rest.Service
}

////publishMSG send websocket message
//func (base *Controller) publishMSG(in TMsg) ([]byte, error) {
//	jmsg := map[string]interface{}{
//		"topic":    in.Topic,
//		"msg":      in.Message,
//		"is_error": in.IsError,
//	}
//	b, err := json.Marshal(jmsg)
//	if err != nil {
//		//panic(err)
//	}
//	if in.IsError {
//		log.Errorf("ERROR: publish websocket message: %v\n", in.Message)
//	} else {
//		log.Infof("INFO: publish websocket message: %v\n", in.Message)
//	}
//	err = base.WS.Broadcast(b)
//	if err != nil {
//		return nil, err
//	}
//	return b, nil
//}

func (base *Controller) resolveHost(ctx *gin.Context) (host *model.Host, useID bool, err error) {
	idName, useID := useHostNameOrID(ctx)
	host, err = base.DB.GetHostByName(idName, useID)
	return host, useID, err
}

func getGitBody(ctx *gin.Context) (dto *git.Git, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func bodyAsJSON(ctx *gin.Context) (interface{}, error) {
	var body interface{} //get the body and put it into an interface
	err = ctx.ShouldBindJSON(&body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func useHostNameOrID(ctx *gin.Context) (idName string, useID bool) {
	hostID := resolveHeaderHostID(ctx)
	hostName := resolveHeaderHostName(ctx)
	if hostID != "" {
		return hostID, true
	} else if hostName != "" {
		return hostName, false
	} else {
		return "", false
	}
}

func resolveHeaderHostID(ctx *gin.Context) string {
	return ctx.GetHeader("host_uuid")
}

func resolveHeaderHostName(ctx *gin.Context) string {
	return ctx.GetHeader("host_name")
}

func resolveHeaderGitToken(ctx *gin.Context) string {
	return ctx.GetHeader("git_token")
}

func reposeHandler(body interface{}, err error, ctx *gin.Context) {
	if err != nil {
		if err == nil {
			ctx.JSON(404, Message{Message: "unknown error"})
		} else {
			if body != nil {
				ctx.JSON(404, body)
			} else {
				ctx.JSON(404, Message{Message: err.Error()})
			}

		}
	} else {
		ctx.JSON(200, body)
	}
}

//hostCopy copy same types from this host to the host needed for ssh.Host
func (base *Controller) hostCopy(host *model.Host) (ssh.Host, error) {
	h := new(ssh.Host)
	err = copier.Copy(&h, &host)
	if err != nil {
		fmt.Println(err)
		return *h, err
	} else {
		return *h, err
	}
}