/*
 * @Date: 2023-10-27 07:56:59
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-04 17:10:17
 * @FilePath: /1024/server/src/controller/common.go
 */
package controller

import "fmt"

// 读缓冲大小
const READ_BUFFER_SIZE = 20480

func errWithStack(err error) string {
	return fmt.Sprintf("%+v\n", err)
}
