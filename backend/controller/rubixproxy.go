package controller

import (
	"errors"
	"fmt"
	"github.com/NubeIO/rubix-updater/model"
	"github.com/NubeIO/rubix-updater/service/rubixmodel"
	"github.com/NubeIO/rubix-updater/utils/rest"
	"github.com/gin-gonic/gin"
	"net/url"
	"strings"
	"time"
)

func getMethod(method string) (out string) {
	out = rest.GET
	switch method {
	case "GET":
		out = rest.GET
	case "PATCH":
		out = rest.PATCH
	case "DELETE":
		out = rest.DELETE
	case "POST":
		out = rest.POST
	}
	return out
}



func (base *Controller) resolveHost(ctx *gin.Context) (host *model.Host, useID bool, err error) {
	idName, useID := useHostNameOrID(ctx)
	host, err = base.GetHostByName(idName, useID)
	return host, useID, err
}

func urlProxyPath(u string, nonProxyReq bool) (clean string) {
	_url := fmt.Sprintf("http://%s", u)
	p, err := url.Parse(_url)
	if err != nil {
		return ""
	}
	var parts []string
	if nonProxyReq {
		parts = strings.SplitAfter(p.String(), "api")
	} else {
		parts = strings.SplitAfter(p.String(), "proxy")
	}

	if len(parts) >= 1 {
		return parts[1]
	} else {
		return ""
	}
}


type proxyOptions struct {
	ctx *gin.Context
	refreshToken bool
	reqOpt rest.ReqOpt
	NonProxyReq bool
}

type proxyReturn struct {
	Token, Method, RequestURL string
	Body interface{}
}

func tokenTimeDiffMin(t time.Time, timeDiff float64) bool {
	t1 := time.Now()
	if t1.Sub(t).Minutes() > timeDiff {
		return true
	} else {
		return false
	}
}


func (base *Controller) buildProxyReq(proxyOptions proxyOptions) (s *rest.Service, options *rest.ReqOpt, rtn proxyReturn, err error) {
	ctx := proxyOptions.ctx
	host, _, err := base.resolveHost(ctx)
	if err != nil {
		return nil, nil, rtn, err
	}
	method := getMethod(ctx.Request.Method)
	rtn.Method = method
	ru := urlProxyPath(ctx.Request.URL.String(), proxyOptions.NonProxyReq)
	rtn.RequestURL = ru
	body, err := bodyAsJSON(ctx)
	rtn.Body = body
	ip := fmt.Sprintf("http://%s:%d", host.IP, host.RubixPort)

	fmt.Println("IP:", ip)

	s = &rest.Service{
		BaseUri: ip,
	}
	token := host.RubixToken
	fmt.Println("UPDATE HOST TOKEN:", "ID", host.ID, host.RubixToken)
	fmt.Println("tokenTimeDiffMin", tokenTimeDiffMin(host.RubixTokenLastUpdate, 15))
	if token == "" || proxyOptions.refreshToken || tokenTimeDiffMin(host.RubixTokenLastUpdate, 15) {
		options = &rest.ReqOpt{
			Timeout:          2 * time.Second,
			RetryCount:       2,
			RetryWaitTime:    2 * time.Second,
			RetryMaxWaitTime: 0,
			Json:             map[string]interface{}{"username": host.RubixUsername, "password": host.RubixPassword},
		}
		req := s.Do(rest.POST, "/api/users/login", options)
		fmt.Println("REQ GET TOKEN:", req.AsString())
		res := new(rubixmodel.TokenResponse)
		err = req.ToInterface(&res)
		if err != nil {
			return nil, nil, rtn, err
		}
		if res.AccessToken == "" {
			return nil, nil, rtn, errors.New("ERROR: failed to get token")
		}
		token = res.AccessToken
		rtn.Token = token
		var h model.Host

		h.RubixToken = token
		h.RubixTokenLastUpdate = time.Now()
		fmt.Println("UPDATE HOST TOKEN:", "ID", host.ID, h.RubixTokenLastUpdate)
		_, err := base.DBUpdateHost(host.ID, &h)
		if err != nil {
			fmt.Println("ERROR: failed to update host token in db", err)
			return nil, nil, rtn, errors.New("ERROR: failed to update host token in db")
		}
	}
	return s, options, rtn, nil
}


func (base *Controller) RubixProxyRequest(ctx *gin.Context) {
	po := proxyOptions{
		ctx: ctx,
		refreshToken: true,
	}
	proxyReq, opt, rtn, err := base.buildProxyReq(po)
	if err != nil {
		reposeHandler(nil, err, ctx)
	} else {
		opt = &rest.ReqOpt{
			Timeout:          500 * time.Second,
			RetryCount:       0,
			RetryWaitTime:    0 * time.Second,
			RetryMaxWaitTime: 0,
			Headers:          map[string]interface{}{"Authorization": rtn.Token},
			Json:             rtn.Body,
		}
		req := proxyReq.Do(rtn.Method, rtn.RequestURL, opt)
		json, err := req.AsJson()
		fmt.Println(req.Err)
		fmt.Println(req.StatusCode)
		if err != nil {
			reposeHandler(nil, err, ctx)
		} else {
			reposeHandler(json, err, ctx)
		}
	}
}
