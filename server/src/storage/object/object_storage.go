/*
 * @Date: 2023-10-26 05:44:54
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-02 07:26:12
 * @FilePath: /1024-dev/1024/server/src/storage/object/object_storage.go
 */
package object

type ObjectStorage interface {
	// 获取对象的下载url
	Load(key string, deadline int64) string
	// 获取上传Token，expire 单位为秒
	GetUploadToken(vkey string,expire uint64) string
}

type ObjectStorageBase struct {
	instance ObjectStorage
}

// GetUploadToken implements ObjectStorage.
func (s *ObjectStorageBase) GetUploadToken(vkey string,expire uint64) string {
	return s.instance.GetUploadToken(vkey,expire)
}

func NewObjectStorage(instance ObjectStorage) ObjectStorage {
	return &ObjectStorageBase{
		instance: instance,
	}
}

func (s *ObjectStorageBase) Load(key string, deadline int64) string {
	return s.instance.Load(key, deadline)
}

var _ ObjectStorage = (*ObjectStorageBase)(nil)
