/*
 * @Date: 2023-10-24 16:01:32
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-29 09:10:57
 * @FilePath: /1024/server/src/controller/video_controller.go
 */
package controller

import (
	"dev1024/src/log"
	"dev1024/src/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (t *VideoController) FetchNextByVID(c *gin.Context) {
	code := http.StatusOK
	resp := RespMsg{
		Status: "ok",
	}
	defer func ()  {
		c.JSON(code,&resp)
	}()
	vidStr := c.Query("vid")
	// 请求第一个视频
	if vidStr == "" {
		code=http.StatusBadRequest
		resp.Status="error"
		resp.Data="Parameter error"
		t.logger.Warn("vidStr is nil")
		return
	}
	vid, err := strconv.ParseInt(vidStr, 10, 64)
	if err != nil {
		code=http.StatusBadRequest
		resp.Status="error"
		resp.Data="Parameter error"
		t.logger.Error(errWithStack(err))
		return
	}
	video, url, err := t.service.FetchNextByVID(vid)
	if err != nil {
		code=http.StatusBadRequest
		resp.Status="error"
		resp.Data="Parameter error"
		t.logger.Error(errWithStack(err))
		return
	}

	video.CDN = url
	resp.Data=gin.H{
		"video": video,
	}
}

func (t *VideoController) ApplyVideoUpToken(c *gin.Context) {
	code := http.StatusOK
	resp := RespMsg{
		Status: "ok",
	}
	defer func() {
		c.JSON(code, &resp)
	}()
	token, err := t.service.ApplyVideoUpToken()
	if err != nil {
		resp.Status = "error"
		t.logger.Error(errWithStack(err))
		return
	}
	resp.Data = gin.H{
		"token": token,
	}
	return
}
