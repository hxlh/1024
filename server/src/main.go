/*
 * @Date: 2023-10-24 03:35:04
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-29 14:07:36
 * @FilePath: /1024/server/src/main.go
 */

package main

import (
	"context"
	"database/sql"
	"dev1024/src/config"
	"dev1024/src/controller"
	"dev1024/src/service"
	"dev1024/src/storage"
	"dev1024/src/storage/object"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-mysql-org/go-mysql/driver"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const CONFIG_PATH = "server.yml"

func main() {
	config.ConfigInit(CONFIG_PATH)
	databaseCfg := config.GetConfig().DataBases["mysql.video1024"]
	username := databaseCfg["username"]
	password := databaseCfg["password"]
	name := databaseCfg["name"]
	host := databaseCfg["host"]
	port := databaseCfg["port"]
	connStr := fmt.Sprintf("%v:%v@%v:%v/%v", username, password, host, port, name)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(errors.WithStack(err))
	}

	redisCfg := config.GetConfig().DataBases["redis"]
	redisOpt := redis.Options{
		Addr:           fmt.Sprintf("%v:%v", redisCfg["host"], redisCfg["port"]),
		DB:             0,
		PoolSize:       100,
		MaxIdleConns:   50,
		MaxActiveConns: 100,
		MinIdleConns:   0,
	}
	redisClient := redis.NewClient(&redisOpt)

	baseDB := &storage.BaseDB{
		DB: db,
		RC: redisClient,
	}

	ctx := context.Background()

	// 日志
	// logger, _ := zap.NewProduction()
	logger, _ := zap.NewDevelopment()
	logger = logger.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	defer logger.Sync()
	sugar := logger.Sugar()

	var objectStorage object.ObjectStorage = object.NewQiNiuObjectStorage(
		config.GetConfig().ObjectStorage.AccessKey,
		config.GetConfig().ObjectStorage.SecretKey,
		config.GetConfig().ObjectStorage.Domain,
		config.GetConfig().ObjectStorage.Bucket,
	)

	videoDao := storage.NewVideoDao(storage.NewVideoDaoImpl(baseDB))
	videoService := service.NewVideoServiceImpl(videoDao, objectStorage)
	videoController := controller.NewVideoController(videoService, sugar)

	accountDao := storage.NewAccountDao(storage.NewAccountDaoImpl(baseDB, ctx))
	accountService := service.NewAccountServiceImpl(accountDao)
	accountController := controller.NewAccountController(accountService, sugar)

	r := gin.Default()

	accountGroup := r.Group("/account")
	accountGroup.POST("/register", accountController.Register)
	accountGroup.POST("/login", accountController.Login)

	videoGroup := r.Group("/video")
	videoGroup.Use(accountController.LoginAuthMiddleware())
	videoGroup.GET("/next", videoController.FetchNextByVID)
	videoGroup.GET("/uptoken", videoController.ApplyVideoUpToken)

	r.Static("/static/", "../static/")
	r.Run(fmt.Sprintf(":%v", config.GetConfig().Server.Port))
}
