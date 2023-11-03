/*
 * @Date: 2023-10-25 05:33:57
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-03 13:09:53
 * @FilePath: /1024/server/src/service/video_service_impl.go
 */
package service

import (
	"context"
	"dev1024/src/entities"
	"dev1024/src/storage"
	"dev1024/src/storage/object"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const VIDEO_UPTOKEN_EXPIRE_TIME = 3600
const VIDEO_CDN_DEAD_TIME = time.Minute * 30

type VideoServiceImpl struct {
	accountDao    storage.AccountDao
	videoDao      storage.VideoDao
	objectStorage object.ObjectStorage
}

// CancelLikeVideo implements VideoService.
func (t *VideoServiceImpl) CancelLikeVideo(ctx context.Context, req entities.CancelLikeVideoReq) (entities.CancelLikeVideoResp, error) {
	// 1.查询该视频对应的tags
	// 2.更新Elasticsearch的user_tags表的tags
	var resp entities.CancelLikeVideoResp
	var tags string

	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return resp, err
	}
	err = t.videoDao.GetBy(ctx, []string{"tags"}, []any{&tags}, "vid", req.Vid)
	defer tx.Commit()
	if err != nil {
		return resp, err
	}
	if tags == "" {
		return resp, nil
	}
	tagArr := strings.Split(tags, ",")
	tagLikes := make([]entities.TagLikes, 0)
	for i := 0; i < len(tagArr); i++ {
		tagLikes = append(tagLikes, entities.TagLikes{
			Tag:   tagArr[i],
			Likes: -1,
		})
	}

	err = t.videoDao.UserTagsIndexDel(ctx, req.Uid, tagLikes)
	if err != nil {
		return resp, err
	}
	err = t.videoDao.UserLikesIndexDel(ctx, req.Vid, req.Uid)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// LikeVideo implements VideoService.
func (t *VideoServiceImpl) LikeVideo(ctx context.Context, req entities.LikeVideoReq) (entities.LikeVideoResp, error) {
	// 1.添加到Elasticsearch的user_like表
	// 2.查询该视频对应的tags
	// 3.更新Elasticsearch的user_tags表的tags
	var resp entities.LikeVideoResp
	err := t.videoDao.UserLikesIndexAdd(ctx, req.Vid, req.Uid, time.Now().Unix())
	if err != nil {
		return resp, err
	}
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return resp, err
	}
	defer tx.Commit()
	var tags string
	err = t.videoDao.GetBy(ctx, []string{"tags"}, []any{&tags}, "vid", req.Vid)
	if err != nil {
		return resp, err
	}
	if tags == "" {
		return resp, nil
	}
	tagArr := strings.Split(tags, ",")
	tagLikes := make([]entities.TagLikes, 0)
	for i := 0; i < len(tagArr); i++ {
		tagLikes = append(tagLikes, entities.TagLikes{
			Tag:   tagArr[i],
			Likes: 1,
		})
	}

	err = t.videoDao.UserTagsIndexUpdate(ctx, req.Uid, tagLikes)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RecommendedVideo implements VideoService.
func (*VideoServiceImpl) RecommendedVideo(ctx context.Context, req entities.RecommendedVideoReq) (entities.RecommendedVideoResp, error) {
	panic("unimplemented")
}

// SearchVideo implements VideoService.
func (t *VideoServiceImpl) SearchVideo(ctx context.Context, key string, offset int, size int) (entities.SearchVideoResp, error) {
	res := entities.SearchVideoResp{
		Info: make([]*entities.SearchVideoRespInfo, 0),
	}
	resp, err := t.videoDao.VideoIndexSearch(ctx, key, offset, size)
	if err != nil {
		return res, err
	}

	hits, ok := resp["hits"].(map[string]interface{})
	fmt.Println(resp)
	panic("SearchVideo not finished")
	if !ok {
		return res, nil
	}
	total := hits["total"].(map[string]interface{})
	hitNum := int(total["value"].(float64))
	if hitNum == 0 {
		return res, nil
	}
	data := hits["hits"].([]interface{})
	for i := 0; i < len(data); i++ {
		info := &entities.SearchVideoRespInfo{}
		item := data[i].(map[string]interface{})
		source := item["_source"].(map[string]interface{})
		highlight, ok := item["highlight"].(map[string]interface{})
		if ok {
			info.HighLightSubtitled = highlight["subtitled"].([]interface{})[0].(string)
		}
		// 解析数据
		info.Vid = uint64(source["vid"].(float64))
		uploaderID := uint64(source["uploader"].(float64))
		info.Subtitled = source["subtitled"].(string)
		info.Tags = source["tags"].(string)
		info.Likes = int64(source["likes"].(float64))
		info.UpLoadTime = int64(source["upload_time"].(float64))
		// 从数据库中查找用户相关信息
		ctx, tx, err := getTxFromCtx(ctx)
		if err != nil {
			return res, err
		}
		account, err := t.accountDao.GetAccountByID(ctx, uploaderID)
		if err != nil {
			return res, commitOrRollback(err, tx)
		}
		info.UpLoaderUsername = account.Username
		info.UpLoaderNickname = account.NickName

		var vkey, imgKey string
		// 从数据库中查找视频相关信息
		err = t.videoDao.GetBy(ctx, []string{"vkey", "thumbnail"}, []any{&vkey, &imgKey}, "vid", info.Vid)
		err = commitOrRollback(err, tx)
		if err != nil {
			return res, err
		}

		// 从对象存储生成Url直链
		info.VideoUrl = t.objectStorage.Load(vkey, time.Now().Add(VIDEO_CDN_DEAD_TIME).Unix())
		info.ThumbnailUrl = t.objectStorage.Load(imgKey, time.Now().Add(VIDEO_CDN_DEAD_TIME).Unix())

		res.Info = append(res.Info, info)
	}

	return res, nil
}

// UpLoadVideoCallBack implements VideoService.
// UpLoadVideoCallBack implements VideoService.
func (t *VideoServiceImpl) UpLoadVideoCallBack(ctx context.Context, uploaderIn uint64, vid uint64) (finalErr error) {
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return err
	}
	defer func(finalErr *error) {
		*finalErr = commitOrRollback(*finalErr, tx)
	}(&finalErr)

	// 查出对应视频是否已完成上传
	uploader, complete, err := t.videoDao.GetUploadCompleteByID(ctx, vid)
	if err != nil {
		finalErr = commitOrRollback(err, tx)
		return
	}
	if complete == 1 {
		finalErr = errors.New("Video has been uploaded.")
		return
	}
	if uploader != uploaderIn {
		finalErr = errors.New("The requester is not the same as the uploader")
		return
	}
	err = t.videoDao.UpdateBy(ctx, []string{"upload_complete"}, []any{1}, "vid", vid)
	if err != nil {
		finalErr = err
		return
	}
	finalErr = nil
	return
}

func (t *VideoServiceImpl) AddToElasticsearch(ctx context.Context, vid uint64) error {
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return err
	}
	videoinfo := entities.VideoInfo{}
	err = t.videoDao.GetBy(ctx, []string{
		"vid", "uploader", "subtitled", "tags", "likes", "upload_time",
	}, []any{
		&videoinfo.Vid, &videoinfo.UpLoader, &videoinfo.Subtitled, &videoinfo.Tags, &videoinfo.Likes, &videoinfo.UpLoadTime,
	}, "vid", vid)
	if err != nil {
		return err
	}
	tx.Commit()

	err = t.videoDao.VideoIndexAdd(ctx, &videoinfo)
	if err != nil {
		return err
	}
	return nil
}

// UpLoadVideo implements VideoService.
func (t *VideoServiceImpl) UpLoadVideo(ctx context.Context, upLoadVideoReq *entities.UpLoadVideoReq) (entities.UpLoadVideoResp, error) {
	var resp entities.UpLoadVideoResp
	// uploader,subtitled,likes,tags,upload_complete
	videoInfo := entities.VideoInfo{
		UpLoader:       upLoadVideoReq.UpLoader,
		Subtitled:      upLoadVideoReq.Subtitled,
		Likes:          0,
		Tags:           upLoadVideoReq.Tags,
		UpLoadComplete: 0,
	}

	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return resp, err
	}
	vid, err := t.videoDao.Create(ctx, &videoInfo)
	if err != nil {
		return resp, commitOrRollback(err, tx)
	}

	vkey := strconv.FormatUint(vid, 10) + ".mp4"
	imgKey := strconv.FormatUint(vid, 10) + ".jpg"

	token := t.GetUpLoadToken(ctx, imgKey, uint64(VIDEO_UPTOKEN_EXPIRE_TIME))

	videoInfo.VKey = vkey
	err = t.videoDao.UpdateBy(ctx, []string{"vkey", "thumbnail"}, []any{vkey, imgKey}, "vid", vid)
	err = commitOrRollback(err, tx)
	if err != nil {
		return resp, err
	}

	resp.VKey = videoInfo.VKey
	resp.Vid = vid
	resp.Token = token
	return resp, nil
}

func (t *VideoServiceImpl) GetUpLoadToken(ctx context.Context, vkey string, expire uint64) string {
	return t.objectStorage.GetUploadToken(vkey, expire)
}

// GetVideoByID implements VideoService.
func (t *VideoServiceImpl) GetVideoByID(ctx context.Context, vid uint64) (string, error) {
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return "", err
	}
	vkey, err := t.videoDao.GetVKeyByID(ctx, vid)
	err = commitOrRollback(err, tx)
	if err != nil {
		return "", err
	}
	ctx, objectStorage, err := getObjectStorageFromCtx(ctx)
	url := objectStorage.Load(vkey, time.Now().Add(VIDEO_CDN_DEAD_TIME).Unix())
	return url, nil
}

func NewVideoServiceImpl(accountDao storage.AccountDao, videoDao storage.VideoDao, objectStorage object.ObjectStorage) *VideoServiceImpl {
	return &VideoServiceImpl{
		videoDao:      videoDao,
		accountDao:    accountDao,
		objectStorage: objectStorage,
	}
}

var _ VideoService = (*VideoServiceImpl)(nil)
