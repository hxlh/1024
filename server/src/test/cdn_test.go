/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2023-10-24 15:32:30
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-26 05:24:47
 * @FilePath: /1024/server/src/test/cdn_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package test

import (
	"dev1024/src/config"
	"fmt"
	"testing"

	_ "github.com/go-mysql-org/go-mysql/driver"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func TestCdn(t *testing.T) {
	mac := qbox.NewMac(config.GetConfig().CDN.AccessKey, config.GetConfig().CDN.SecretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	// cfg.Region=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(mac, &cfg)
	bucket := "video1024"
	key := "hxlh/video.mp4"
	fileInfo, err := bucketManager.Stat(bucket, key)
	if err != nil {
		fmt.Println(err)
		return
	}
	bucketManager.Get(bucket, key, &storage.GetObjectInput{})
	fmt.Println(fileInfo.String())
	// 可以解析文件的PutTime
	fmt.Println(storage.ParsePutTime(fileInfo.PutTime))
}

func TestDB(t *testing.T) {
	
}
