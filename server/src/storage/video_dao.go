/*
 * @Date: 2023-10-24 16:29:27
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-02 15:12:41
 * @FilePath: /1024-dev/1024/server/src/storage/video_dao.go
 */
package storage

import (
	"context"
	"dev1024/src/entities"
)

type VideoDao interface {
	Create(ctx context.Context, upLoadVideoReq *entities.VideoInfo) (uint64, error)
	GetVKeyByID(ctx context.Context, vid uint64) (string, error)
	GetUploadCompleteByID(ctx context.Context, vid uint64) (uint64, int, error)
	SearchVideo(ctx context.Context, key string, offset int, size int) (map[string]interface{}, error)
	GetBy(ctx context.Context, keys []string, values []any, by string, byValue any) error
	UpdateBy(ctx context.Context, keys []string, values []any, by string, byValue any) error
	AddToElasticsearch(ctx context.Context, videoinfo *entities.VideoInfo) error
}

type VideoDaoBase struct {
	instance VideoDao
}

// UpdateBy implements VideoDao.
func (t*VideoDaoBase) UpdateBy(ctx context.Context, keys []string, values []any, by string, byValue any) error {
	return t.instance.UpdateBy(ctx,keys,values,by,byValue)
}

// AddToElasticsearch implements VideoDao.
func (t *VideoDaoBase) AddToElasticsearch(ctx context.Context, videoinfo *entities.VideoInfo) error {
	return t.instance.AddToElasticsearch(ctx, videoinfo)
}

// GetBy implements VideoDao.
func (t *VideoDaoBase) GetBy(ctx context.Context, keys []string, values []any, by string, byValue any) error {
	return t.instance.GetBy(ctx, keys, values, by, byValue)
}

// SearchVideo implements VideoDao.
func (t *VideoDaoBase) SearchVideo(ctx context.Context, key string, offset int, size int) (map[string]interface{}, error) {
	return t.instance.SearchVideo(ctx, key, offset, size)
}

// GetUploadCompleteByID implements VideoDao.
func (t *VideoDaoBase) GetUploadCompleteByID(ctx context.Context, vid uint64) (uint64, int, error) {
	return t.instance.GetUploadCompleteByID(ctx, vid)
}

// GetVKeyByID implements VideoDao.
func (t *VideoDaoBase) GetVKeyByID(ctx context.Context, vid uint64) (string, error) {
	return t.instance.GetVKeyByID(ctx, vid)
}

func (t *VideoDaoBase) Create(ctx context.Context, videoInfo *entities.VideoInfo) (uint64, error) {
	return t.instance.Create(ctx, videoInfo)
}

func NewVideoDao(instance VideoDao) VideoDao {
	return &VideoDaoBase{
		instance: instance,
	}
}

var _ VideoDao = (*VideoDaoBase)(nil)
