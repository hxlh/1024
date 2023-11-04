/*
 * @Date: 2023-10-24 16:29:27
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-28 09:52:47
 * @FilePath: /1024/server/src/storage/video_dao.go
 */
package storage

import (
	"dev1024/src/entities"
)

type VideoDao interface {
	Save(videoInfo *entities.VideoInfo) error
	// 获取比vid大的n个视频的信息
	GetNextNByVid(vid int64, n int) ([]*entities.VideoInfo, error)
}

type VideoDaoBase struct {
	instance VideoDao
}

func NewVideoDao(instance VideoDao) VideoDao {
	return &VideoDaoBase{
		instance: instance,
	}
}

func (t *VideoDaoBase) GetNextNByVid(vid int64, n int) ([]*entities.VideoInfo, error) {
	return t.instance.GetNextNByVid(vid, n)
}

func (t *VideoDaoBase) Save(videoInfo *entities.VideoInfo) error {
	return t.instance.Save(videoInfo)
}

var _ VideoDao = (*VideoDaoBase)(nil)
