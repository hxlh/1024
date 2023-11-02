/*
 * @Date: 2023-10-28 11:34:17
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-01 16:34:55
 * @FilePath: /1024/server/src/entities/account.go
 */
package entities

import "time"

type Account struct {
	Uid          uint64    `json:"uid"`
	Username     string    `json:"username"`
	NickName     string    `json:"nickname"`
	Pwd          string    `json:"pwd"`
	RegisterTime time.Time `json:"register_time"`
	Avatar       string    `json:"avatar"`
}

type LoginReq struct {
	UserName string `json:"username"`
	Pwd      string `json:"pwd"`
}

type LoginResp struct {
	Uid          uint64    `json:"uid"`
	Username     string    `json:"username"`
	NickName     string    `json:"nickname"`
	Avatar       string    `json:"avatar"`
}

type RegisterReq struct {
	UserName string `json:"username"`
	Pwd      string `json:"pwd"`
	NickName string `json:"nickname"`
}

type RegisterResp struct {
	RespMsg
}
