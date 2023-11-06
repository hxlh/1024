/*
 * @Date: 2023-10-24 16:00:23
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-26 05:25:02
 * @FilePath: /1024/server/src/cdn/cdn.go
 */
package cdn

import (
	"dev1024/src/config"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

var mac *auth.Credentials
var bucketManager *storage.BucketManager

func init() {
	mac = qbox.NewMac(config.GetConfig().CDN.AccessKey, config.GetConfig().CDN.SecretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	// cfg.Region=&storage.ZoneHuabei
	bucketManager = storage.NewBucketManager(mac, &cfg)
}

func GenPrivateFileUrl(key string, deadLine int64) string {
	return storage.MakePrivateURL(mac, config.GetConfig().CDN.Domain, key, deadLine)
}

func GetMac() *auth.Credentials {
	return mac
}

func GetBucketManager() *storage.BucketManager {
	return bucketManager
}
