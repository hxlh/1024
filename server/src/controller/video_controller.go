/*
 * @Date: 2023-10-24 16:01:32
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-02 14:16:40
 * @FilePath: /1024-dev/1024/server/src/controller/video_controller.go
 */
package controller

import (
	"context"
	"dev1024/src/entities"
	"dev1024/src/log"
	"dev1024/src/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type VideoController struct {
	service service.VideoService
	logger  log.Logger
}

func NewVideoController(service service.VideoService, logger log.Logger) *VideoController {
	return &VideoController{
		service: service,
		logger:  logger,
	}
}

func (t *VideoController) getBodyToJson(c *gin.Context, code *int, resp *entities.RespMsg, obj any) error {
	buf := make([]byte, READ_BUFFER_SIZE)
	n, err := c.Request.Body.Read(buf)
	if n <= 0 && err != nil {
		*code = http.StatusBadRequest
		resp.Status = "error"
		resp.Data = "Read request failed."
		return errors.WithStack(err)
	}

	err = json.Unmarshal(buf[:n], obj)
	if err != nil {
		*code = http.StatusBadRequest
		resp.Status = "error"
		resp.Data = "Failed to parse request."
		return errors.WithStack(err)
	}

	return nil
}

func (t *VideoController) UpLoadVideo(c *gin.Context) {
	ctxValue, _ := c.Get("ctx")
	ctx := ctxValue.(context.Context)
	code := http.StatusOK
	resp := entities.RespMsg{
		Status: "ok",
	}
	defer func() {
		c.JSON(code, &resp)
	}()

	upLoadVideoReq := entities.UpLoadVideoReq{}
	if err := t.getBodyToJson(c, &code, &resp, &upLoadVideoReq); err != nil {
		t.logger.Error(errWithStack(err))
		return
	}

	claimObj, _ := c.Get("claim")
	claim := claimObj.(*service.MyCustomClaims)
	upLoadVideoReq.UpLoader = claim.Uid

	upLoadVideoResp, err := t.service.UpLoadVideo(ctx, &upLoadVideoReq)
	if err != nil {
		code = http.StatusBadRequest
		resp.Status = "error"
		resp.Data = "token request failed"
		t.logger.Error(errWithStack(err))
		return
	}
	resp.Data = upLoadVideoResp
	return
}

func (t *VideoController) UpLoadVideoCallBack(c *gin.Context) {
	ctxValue, _ := c.Get("ctx")
	ctx := ctxValue.(context.Context)
	code := http.StatusOK
	resp := entities.RespMsg{
		Status: "ok",
	}
	defer func() {
		c.JSON(code, &resp)
	}()
	claimValue, _ := c.Get("claim")
	claim := claimValue.(*service.MyCustomClaims)

	req := entities.UpLoadVideoCallBackReq{}
	err := t.getBodyToJson(c, &code, &resp, &req)
	if err != nil {
		return
	}
	err = t.service.UpLoadVideoCallBack(ctx, claim.Uid, req.Vid)
	if err != nil {
		code = http.StatusInternalServerError
		resp.Status = "error"
		resp.Data = err.Error()
		return
	}
	err = t.service.AddToElasticsearch(ctx, req.Vid)
	if err != nil {
		code = http.StatusInternalServerError
		resp.Status = "error"
		resp.Data = err.Error()
		return
	}
	return
}

func (t *VideoController) GetVideoByID(c *gin.Context) {
	ctxValue, _ := c.Get("ctx")
	ctx := ctxValue.(context.Context)
	code := http.StatusOK
	resp := entities.RespMsg{
		Status: "ok",
	}
	defer func() {
		c.JSON(code, resp)
	}()

	vidStr := c.Query("vid")
	vid, err := strconv.ParseUint(vidStr, 10, 64)
	if err != nil {
		code = http.StatusBadRequest
		resp.Status = "error"
		resp.Data = "Parsing parameter failure"
		return
	}

	url, err := t.service.GetVideoByID(ctx, vid)
	if err != nil {
		code = http.StatusBadRequest
		resp.Status = "error"
		resp.Data = err.Error()
		return
	}

	resp.Data = entities.GetVideoByIDResp{
		Url: url,
	}
}

func (t *VideoController) SearchVideo(c *gin.Context) {
	ctxValue, _ := c.Get("ctx")
	ctx := ctxValue.(context.Context)
	code := http.StatusOK
	resp := entities.RespMsg{
		Status: "ok",
	}
	defer func() {
		c.JSON(code, resp)
	}()

	req := entities.SearchVideoReq{}
	err := t.getBodyToJson(c, &code, &resp, &req)
	if err != nil {
		return
	}
	res, err := t.service.SearchVideo(ctx, req.Key, req.Offset, 5)
	if err != nil {
		code = http.StatusBadRequest
		resp.Status = "error"
		resp.Data = err.Error()
		return
	}

	resp.Data = res
	return
}
