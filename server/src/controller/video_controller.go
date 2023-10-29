/*
 * @Date: 2023-10-24 16:01:32
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-29 13:40:05
 * @FilePath: /1024/server/src/controller/video_controller.go
 */
package controller

import (
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

func (t *VideoController) FetchNextByVID(c *gin.Context) {
	code := http.StatusOK
	resp := RespMsg{
		Status: "ok",
	}
	defer func() {
		c.JSON(code, &resp)
	}()
	vidStr := c.Query("vid")
	// 请求第一个视频
	if vidStr == "" {
		code = http.StatusBadRequest
		resp.Status = "error"
		resp.Data = "Parameter error"
		t.logger.Warn("vidStr is nil")
		return
	}
	vid, err := strconv.ParseInt(vidStr, 10, 64)
	if err != nil {
		code = http.StatusBadRequest
		resp.Status = "error"
		resp.Data = "Parameter error"
		t.logger.Error(errWithStack(err))
		return
	}
	video, url, err := t.service.FetchNextByVID(vid)
	if err != nil {
		code = http.StatusBadRequest
		resp.Status = "error"
		resp.Data = "Parameter error"
		t.logger.Error(errWithStack(err))
		return
	}

	video.CDN = url
	resp.Data = gin.H{
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

	videoInfo := entities.VideoInfo{}
	if err := t.getBodyToJson(c, &code, &resp, &videoInfo); err != nil {
		t.logger.Error(errWithStack(err))
		return
	}

	claimObj, _ := c.Get("claim")
	claim := claimObj.(*service.MyCustomClaims)
	videoInfo.UpLoader = claim.Uid

	token, err := t.service.ApplyVideoUpToken(&videoInfo)
	if err != nil {
		code = http.StatusBadRequest
		resp.Status = "error"
		resp.Data = "token request failed"
		t.logger.Error(errWithStack(err))
		return
	}

	resp.Data = entities.ApplyVideoUpTokenResp{
		Vid:   videoInfo.Vid,
		Token: token,
	}
	return
}

func (t *VideoController) getBodyToJson(c *gin.Context, code *int, resp *RespMsg, obj any) error {
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
