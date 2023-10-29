package entities

import "time"

type Account struct {
	Uid          uint64     `json:"uid"`
	Username     string    `json:"username"`
	NickName     string    `json:"nickname"`
	Pwd          string    `json:"pwd"`
	RegisterTime time.Time `json:"register_time"`
	Avatar       string    `json:"avatar"`
}

type LoginInfo struct{
	Uid uint64 `json:"uid"`
	UserName string `json:"username"`
	Token string `json:"token"`
}

type LoginReq struct{
	UserName string `json:"username"`
	Pwd string `json:"pwd"`
}

type LoginResp struct{
	Uid uint64 `json:"uid"`
	UserName string `json:"username"`
}

type RegisterReq struct{
	UserName string `json:"username"`
	Pwd string `json:"pwd"`
	NickName string `json:"nickname"`
}