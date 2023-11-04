/*
 * @Date: 2023-10-27 07:56:30
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-29 09:08:12
 * @FilePath: /1024/server/src/controller/account_controller.go
 */
package controller

import (
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

func (t *AccountController) getBodyToJson(c *gin.Context, code *int, resp *RespMsg, obj any) error {
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
	code := http.StatusOK
	resp := RespMsg{
		Status: "ok",
	}
	defer func() {
		c.JSON(code, &resp)
	}()

	req := entities.LoginReq{}

	if err := t.getBodyToJson(c, &code, &resp, &req); err != nil {
		t.logger.Error(errWithStack(err))
		return
	}

	info, err := t.service.Login(req.UserName, req.Pwd)
	if err != nil {
		code = http.StatusBadRequest
		resp.Status = "error"
		resp.Data = err.Error()
		t.logger.Error(errWithStack(err))
		return
	}
	c.SetCookie("jwt_token", info.Token, 3600, "", "", false, true)

	loginResp := entities.LoginResp{
		Uid:      info.Uid,
		UserName: info.UserName,
	}
	resp.Data = loginResp
}

func (t *AccountController) Register(c *gin.Context) {
	code := http.StatusOK
	resp := RespMsg{
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

	err := t.service.Register(&entities.Account{
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
		code := http.StatusBadRequest
		resp := RespMsg{
			Status: "error",
			Data:   "Authentication failure",
		}
		token, err := c.Request.Cookie("jwt_token")
		if err != nil {
			c.JSON(code, &resp)
			c.Abort()
			t.logger.Error(errWithStack(err))
			return
		}
		err = t.service.JwtAuth(token.Value)
		if err != nil {
			c.JSON(code, &resp)
			c.Abort()
			t.logger.Error(errWithStack(err))
			return
		}

		c.Next()
	}
}
