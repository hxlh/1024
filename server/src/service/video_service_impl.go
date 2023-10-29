/*
 * @Date: 2023-10-25 05:33:57
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-29 13:54:20
 * @FilePath: /1024/server/src/service/video_service_impl.go
 */
package service

import (
	"dev1024/src/entities"
	"dev1024/src/storage"
	"dev1024/src/storage/object"
	"time"

	"github.com/pkg/errors"
)

const VIDEO_UPTOKEN_EXPIRE_TIME = 3600

type VideoServiceImpl struct {
	videoDao      storage.VideoDao
	objectStorage object.ObjectStorage
}

// ApplyVideoUpToken implements VideoService.
func (t *VideoServiceImpl) ApplyVideoUpToken(video *entities.VideoInfo) (string, error) {
	err := t.videoDao.Save(video)
	if err != nil {
		return "", err
	}
	token := t.objectStorage.GetUpToken(uint64(VIDEO_UPTOKEN_EXPIRE_TIME))
	return token, nil
}

func (t *VideoServiceImpl) FetchNextByVID(vid int64) (*entities.VideoInfo, string, error) {
	videoInfos, err := t.videoDao.GetNextNByVid(vid, 1)
	if err != nil {
		return nil, "", err
	}
	if len(videoInfos) == 0 {
		return nil, "", errors.New("Couldn't find the right video")
	}

	video := videoInfos[0]
	deadline := time.Now().Add(time.Minute * 30).Unix()
	url := t.objectStorage.Load(video.CDN, deadline)
	return video, url, nil
}

func NewVideoServiceImpl(videoDao storage.VideoDao, objectStorage object.ObjectStorage) *VideoServiceImpl {
	return &VideoServiceImpl{
		videoDao:      videoDao,
		objectStorage: objectStorage,
	}
}

var _ VideoService = (*VideoServiceImpl)(nil)
