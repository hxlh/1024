/*
 * @Date: 2023-10-24 16:01:32
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-26 05:24:57
 * @FilePath: /1024/server/src/controller/video_controller.go
 */
package controller

import (
	"dev1024/src/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	service *service.VideoService
}

func NewVideoController(service *service.VideoService) *VideoController {
	return &VideoController{
		service: service,
	}
}

func (vc *VideoController) FetchNext(c *gin.Context) {
	c.String(http.StatusOK, "Hello, Gin!")
}
