/*
 * @Date: 2023-10-27 07:58:46
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-01 16:37:53
 * @FilePath: /1024/server/src/service/account_service.go
 */
package service

import (
	"context"
	"dev1024/src/entities"
)

type AccountService interface {
	JwtAuth(ctx context.Context, token string) (*MyCustomClaims, error)
	Login(ctx context.Context, username string, password string) (entities.Account,string, error)
	Register(ctx context.Context, account *entities.Account) error
	GenJwtSignKey(ctx context.Context) string
}
