package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (base *Controller) clearDir(id, path string) (result bool, err error) {
	c := base.newClient(id)
	defer c.Close()
	command := fmt.Sprintf("sudo rm %s/*", path)
	_, err = c.Run(command)
	if err != nil {
		return false, err
	}
	return true, err
}

func (base *Controller) ClearDir(ctx *gin.Context) {
	body := dirBody(ctx)
	id := ctx.Params.ByName("id")
	dir, err := base.clearDir(id, body.Path)
	if err != nil {
		reposeHandler(nil, err, ctx)
	} else {
		reposeHandler(dir, err, ctx)
	}
}