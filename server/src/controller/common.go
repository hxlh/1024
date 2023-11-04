/*
 * @Date: 2023-10-27 07:56:59
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-29 09:09:15
 * @FilePath: /1024/server/src/controller/common.go
 */
package controller

import "fmt"

type RespMsg struct {
	Status string      `json:"status "`
	Data   interface{} `json:"data"`
}

// 读缓冲大小
const READ_BUFFER_SIZE = 2048

func errWithStack(err error) string {
	return fmt.Sprintf("%+v\n", err)
}
