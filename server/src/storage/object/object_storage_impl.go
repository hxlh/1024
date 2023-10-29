/*
 * @Date: 2023-10-26 05:46:07
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-26 13:03:24
 * @FilePath: /1024/server/src/storage/object/object_storage_impl.go
 */
package object

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiNiuObjectStorage struct {
	mac           *auth.Credentials
	bucketManager *storage.BucketManager
	domain        string
	bucket        string
}

func (q *QiNiuObjectStorage) GetUpToken(expire uint64) string {
	putPolicy := storage.PutPolicy{
		Expires: expire,
		Scope:   q.bucket,
	}
	upToken := putPolicy.UploadToken(q.mac)
	return upToken
}

func NewQiNiuObjectStorage(accessKey, secretKey, domain, bucket string) *QiNiuObjectStorage {
	q := &QiNiuObjectStorage{
		domain: domain,
	}

	q.mac = qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	// cfg.Region=&storage.ZoneHuabei
	q.bucketManager = storage.NewBucketManager(q.mac, &cfg)
	return q
}

func (q *QiNiuObjectStorage) Load(key string, deadline int64) string {
	return storage.MakePrivateURL(q.mac, q.domain, key, deadline)
}

var _ ObjectStorage = (*QiNiuObjectStorage)(nil)
