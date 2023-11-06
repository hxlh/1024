/*
 * @Date: 2023-10-24 03:35:04
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-26 05:31:36
 * @FilePath: /1024/server/src/main.go
 */

package main

import (
	"dev1024/src/config"
	"dev1024/src/controller"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// controller instance
	vc := controller.NewVideoController(nil)

	r.GET("/hello", vc.FetchNext)

	r.Static("/static/", "../static/")
	r.Run(fmt.Sprintf(":%v", config.GetConfig().Server.Port))
}
