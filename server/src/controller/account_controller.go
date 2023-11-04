/*
 * @Date: 2023-10-27 07:56:30
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-03 12:46:18
 * @FilePath: /1024/server/src/controller/account_controller.go
 */
package controller

import (
	"context"
	"dev1024/src/entities"
	"dev1024/src/log"
	"dev1024/src/service"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type AccountController struct {
	service service.AccountService
	logger  log.Logger
}

func NewAccountController(accountService service.AccountService, logger log.Logger) *AccountController {
	return &AccountController{
		service: accountService,
		logger:  logger,
	}
}

func (t *AccountController) getBodyToJson(c *gin.Context, code *int, resp *entities.RespMsg, obj any) error {
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

func (t *AccountController) Login(c *gin.Context) {
	ctxValue, _ := c.Get("ctx")
	ctx := ctxValue.(context.Context)
	code := http.StatusOK
	resp := entities.RespMsg{
		Status: "ok",
	}
	defer func() {
		c.JSON(code, &resp)
	}()

	req := entities.LoginReq{}

	if err := t.getBodyToJson(c, &code, &resp, &req); err != nil {
		code = http.StatusBadRequest
		resp.Status = "error"
		resp.Data = "Request parse error"
		t.logger.Error(errWithStack(err))
		return
	}

	account, token, err := t.service.Login(ctx, req.UserName, req.Pwd)
	if err != nil {
		code = http.StatusBadRequest
		resp.Status = "error"
		resp.Data = err.Error()
		t.logger.Error(errWithStack(err))
		return
	}

	loginResp := entities.LoginResp{
		Uid:      account.Uid,
		Username: account.Username,
		NickName: account.NickName,
		Avatar:   account.Avatar,
		Token:    token,
	}
	resp.Data = loginResp
}

func (t *AccountController) Register(c *gin.Context) {
	ctxValue, _ := c.Get("ctx")
	ctx := ctxValue.(context.Context)
	code := http.StatusOK
	resp := entities.RespMsg{
		Status: "ok",
	}
	defer func() {
		c.JSON(code, &resp)
	}()

	req := entities.RegisterReq{}

	if err := t.getBodyToJson(c, &code, &resp, &req); err != nil {
		t.logger.Error(errWithStack(err))
		return
	}

	err := t.service.Register(ctx, &entities.Account{
		Username: req.UserName,
		NickName: req.NickName,
		Pwd:      req.Pwd,
	})
	if err != nil {
		code = http.StatusBadRequest
		resp.Data = "Registration failed."
		t.logger.Error(errWithStack(err))
		return
	}
}

func (t *AccountController) LoginAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxValue, _ := c.Get("ctx")
		ctx := ctxValue.(context.Context)
		code := http.StatusBadRequest
		resp := entities.RespMsg{
			Status: "error",
			Data:   "Authentication failure",
		}
		token, ok := c.Request.Header["Authorization"]
		if !ok{
			c.JSON(code, &resp)
			c.Abort()
			return
		}
		claim, err := t.service.JwtAuth(ctx, token[0])
		if err != nil {
			c.JSON(code, &resp)
			c.Abort()
			t.logger.Error(errWithStack(err))
			return
		}
		c.Set("claim", claim)
		c.Next()
	}
}
