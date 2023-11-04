/*
 * @Date: 2023-10-27 07:58:46
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-28 13:16:48
 * @FilePath: /1024/server/src/service/account_service.go
 */
package service

import "dev1024/src/entities"

type AccountService interface{
	JwtAuth(token string) error
	Login(username string,password string) (entities.LoginInfo,error)
	Register(account *entities.Account) error
	GenJwtSignKey() string
}
