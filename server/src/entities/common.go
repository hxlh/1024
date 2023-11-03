/*
 * @Date: 2023-11-01 09:39:37
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-03 13:53:01
 * @FilePath: /1024/server/src/entities/common.go
 */
package entities

type RespMsg struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
