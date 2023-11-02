/*
 * @Date: 2023-10-24 03:35:04
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-02 13:05:56
 * @FilePath: /1024-dev/1024/server/src/main.go
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
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	_ "github.com/go-mysql-org/go-mysql/driver"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const CONFIG_PATH = "server.yml"

func CreateDB() *sql.DB {
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
	return db
}

func CreateRedis() *redis.Client {
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
	return redisClient
}

func CreateObjectStorage() object.ObjectStorage {
	var objectStorage object.ObjectStorage = object.NewQiNiuObjectStorage(
		config.GetConfig().ObjectStorage.AccessKey,
		config.GetConfig().ObjectStorage.SecretKey,
		config.GetConfig().ObjectStorage.Domain,
		config.GetConfig().ObjectStorage.Bucket,
	)
	return objectStorage
}

func CreateElasticSearchClient() *elasticsearch.Client {
	elasticCfg := config.GetConfig().DataBases["elastic"]
	addr := fmt.Sprintf("https://%v:%v", elasticCfg["host"], elasticCfg["port"])
	username := elasticCfg["username"].(string)
	password := elasticCfg["password"].(string)

	crtFile, err := os.Open("http_ca.crt")
	if err != nil {
		panic(err)
	}
	defer crtFile.Close()
	crt, err := io.ReadAll(crtFile)
	if err != nil {
		panic(err)
	}
	elasticClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{addr},
		Username:  username,
		Password:  password,
		CACert:    crt,
	})
	if err != nil {
		panic(err)
	}
	_, err = elasticClient.Info()
	if err != nil {
		panic(err)
	}
	return elasticClient
}

func main() {
	config.ConfigInit(CONFIG_PATH)

	db := CreateDB()

	redisClient := CreateRedis()

	objectStorage := CreateObjectStorage()

	elasticPool := &sync.Pool{
		New: func() interface{} {
			// 如果池中没有对象可用，这个函数会创建一个新的对象
			return CreateElasticSearchClient()
		},
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "db", db)
	ctx = context.WithValue(ctx, "redis", redisClient)
	ctx = context.WithValue(ctx, "object_storage", objectStorage)
	ctx = context.WithValue(ctx, "elastic.pool", elasticPool)

	// 日志
	// logger, _ := zap.NewProduction()
	logger, _ := zap.NewDevelopment()
	logger = logger.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	defer logger.Sync()
	sugar := logger.Sugar()

	accountDao := storage.NewAccountDao(storage.NewAccountDaoImpl())
	videoDao := storage.NewVideoDao(storage.NewVideoDaoImpl())

	accountService := service.NewAccountServiceImpl(accountDao)
	videoService := service.NewVideoServiceImpl(accountDao, videoDao, objectStorage)

	videoController := controller.NewVideoController(videoService, sugar)
	accountController := controller.NewAccountController(accountService, sugar)

	r := gin.Default()
	r.Use(SetCrossDomain, func(c *gin.Context) {
		c.Set("ctx", ctx)
	})

	accountGroup := r.Group("/account")
	accountGroup.POST("/register", accountController.Register)
	accountGroup.POST("/login", accountController.Login)

	videoGroup := r.Group("/video")
	videoGroup.Use(accountController.LoginAuthMiddleware())
	videoGroup.GET("/url", videoController.GetVideoByID)
	videoGroup.POST("/upload", videoController.UpLoadVideo)
	videoGroup.POST("/upload_callback", videoController.UpLoadVideoCallBack)
	videoGroup.POST("/search", videoController.SearchVideo)

	r.Static("/static/", "static/")
	r.Run(fmt.Sprintf(":%v", config.GetConfig().Server.Port))
}

func SetCrossDomain(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()
}
