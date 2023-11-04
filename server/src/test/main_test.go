/*
 * @Date: 2023-10-24 15:32:30
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-02 07:17:20
 * @FilePath: /1024-dev/1024/server/src/test/main_test.go
 */
package test

import (
	"database/sql"
	"dev1024/src/config"
	"dev1024/src/storage"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/go-mysql-org/go-mysql/driver"
)

func TestTemp(t *testing.T) {
	config.ConfigInit("../server.yml")
	video1024 := config.GetConfig().DataBases["mysql.video1024"]
	username := video1024["username"]
	password := video1024["password"]
	name := video1024["name"]
	host := video1024["host"]
	port := video1024["port"]
	connStr := fmt.Sprintf("%v:%v@%v:%v/%v", username, password, host, port, name)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	err = storage.UpdateTableBy(tx, "video1024.video_info", []string{"uploader", "vkey"}, []any{1222, "1222.vdd"}, "vid", 1)
	if err != nil {
		panic(err)
	}

	tx.Commit()
}

func TestUrlSafeBase64(t *testing.T) {
	res := base64.URLEncoding.EncodeToString([]byte("video1024:11111111.jpg"))
	fmt.Println(res)
}

// func TestCdn(t *testing.T) {
// 	objectStorage := object.NewObjectStorage(object.NewQiNiuObjectStorage(
// 		config.GetConfig().ObjectStorage.AccessKey,
// 		config.GetConfig().ObjectStorage.SecretKey,
// 		config.GetConfig().ObjectStorage.Domain, config.GetConfig().ObjectStorage.Bucket,
// 	),
// 	)
// 	url := objectStorage.Load("hxlh/video.mp4", time.Now().Add(time.Minute*30).Unix())
// 	fmt.Println(url)
// }

// func TestDB(t *testing.T) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			var err error
// 			if v, ok := r.(error); !ok {
// 				// 使用 pkg/errors 包捕获 panic
// 				err = errors.New("Recovered: " + fmt.Sprint(r))
// 			} else {
// 				err = v
// 			}
// 			// 打印堆栈跟踪
// 			stackTrace := errors.WithStack(err)
// 			fmt.Printf("%+v\n", stackTrace)
// 		}
// 	}()
// 	config.ConfigInit("../server.yml")
// 	video1024 := config.GetConfig().DataBases["mysql.video1024"]
// 	username := video1024["username"]
// 	password := video1024["password"]
// 	name := video1024["name"]
// 	host := video1024["host"]
// 	port := video1024["port"]
// 	connStr := fmt.Sprintf("%v:%v@%v:%v/%v", username, password, host, port, name)
// 	db, err := sql.Open("mysql", connStr)

// 	baseDB := &storage.BaseDB{
// 		DB: db,
// 	}

// 	var videoDao storage.VideoDao = storage.NewVideoDao(storage.NewVideoDaoImpl(baseDB))
// 	ans, err := videoDao.GetNextNByVid(-110, 20)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, v := range ans {
// 		fmt.Println(v)
// 	}

// 	info := entities.VideoInfo{
// 		UpLoader:  1111,
// 		CDN:       "cdn",
// 		Subtitled: "test",
// 		Likes:     100,
// 		Tags:      "运动",
// 	}
// 	err = videoDao.Save(&info)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer db.Close()
// }

// type MyCustomClaims struct {
// 	jwt.RegisteredClaims
// 	Uid       int64  `json:"uid"`
// 	Username  string `json:"username"`
// 	LoginTime int64  `json:"deadline"`
// }

// func TestJwt(t *testing.T) {
// 	claim := MyCustomClaims{
// 		Uid:       111111,
// 		Username:  "dasfdsf",
// 		LoginTime: time.Now().Unix(),
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(1 * time.Second)},
// 		},
// 	}
// 	key := []byte("mykey")

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
// 	tokenSigned, err := token.SignedString(key)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(tokenSigned)

// 	time.Sleep(time.Second * 1)
// 	parseClaim := MyCustomClaims{}
// 	_, err = jwt.ParseWithClaims(tokenSigned, &parseClaim, func(t *jwt.Token) (interface{}, error) {
// 		return key, nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(parseClaim.Username)
// }

// func TestRedis(t *testing.T) {
// 	ctx := context.Background()
// 	config.ConfigInit("../server.yml")
// 	redisCfg := config.GetConfig().DataBases["redis"]
// 	addr := fmt.Sprintf("%v:%v", redisCfg["host"], redisCfg["port"])
// 	rc := redis.NewClient(&redis.Options{
// 		Addr: addr,
// 		DB:   0,
// 	})
// 	defer rc.Close()

// 	baseDB := storage.BaseDB{
// 		RC: rc,
// 	}

// 	accountDao := storage.NewAccountDao(storage.NewAccountDaoImpl(&baseDB, ctx))

// 	accountDao.DelJwtSignedKey()
// 	accountDao.SetJwtSignedKey("videodasfdsfsd1213213213")
// 	ans, err := accountDao.GetJwtSignedKey()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(ans)
// }

// func TestGenJwtSignKey(t *testing.T) {
// 	s := service.NewAccountServiceImpl(nil)
// 	for i := 0; i < 10; i++ {
// 		fmt.Println(s.GenJwtSignKey())
// 	}
// }

// func TestRegisterAccount(t *testing.T) {
// 	config.ConfigInit("../server.yml")
// 	video1024 := config.GetConfig().DataBases["mysql.video1024"]
// 	username := video1024["username"]
// 	password := video1024["password"]
// 	name := video1024["name"]
// 	host := video1024["host"]
// 	port := video1024["port"]
// 	connStr := fmt.Sprintf("%v:%v@%v:%v/%v", username, password, host, port, name)
// 	db, err := sql.Open("mysql", connStr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	baseDB := &storage.BaseDB{
// 		DB: db,
// 	}

// 	accountDao := storage.NewAccountDao(storage.NewAccountDaoImpl(baseDB, context.Background()))
// 	for i := 0; i < 10; i++ {
// 		account := entities.Account{
// 			Username: fmt.Sprintf("hxlh%v", i),
// 			NickName: "deepcat",
// 			Pwd:      fmt.Sprintf("pwd%v", i),
// 		}
// 		uid, err := accountDao.Save(&account)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(uid)
// 	}

// 	for i := 0; i < 10; i++ {
// 		username := fmt.Sprintf("hxlh%v", i)
// 		uid, pwd, err := accountDao.GetUidAndPwdByUsername(username)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(uid, pwd)
// 	}
// }

func TestElasticSearch(t *testing.T) {
	config.ConfigInit("../server.yml")
	elasticCfg := config.GetConfig().DataBases["elastic"]
	addr := fmt.Sprintf("https://%v:%v", elasticCfg["host"], elasticCfg["port"])
	username := elasticCfg["username"].(string)
	password := elasticCfg["password"].(string)

	crtFile, err := os.Open("../http_ca.crt")
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
	elasticClient.Create("user", "1", strings.NewReader(""))
	elasticClient.Indices.PutMapping([]string{"user"}, strings.NewReader(
		`
		{
			"properties":{
				"title":{
					"type":"text",
					"index":true,
					"analyzer":"ik_max_word",
					"search_analyzer":"ik_smart"
				},
				"category":{
					"type":"keyword",
					"index":true
				},
				"price":{
					"type":"long",
					"index":true
				}
			}
		}
		`,
	))
}
