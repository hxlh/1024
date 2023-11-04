/*
 * @Date: 2023-10-24 16:29:27
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-04 16:06:42
 * @FilePath: /1024/server/src/storage/video_dao.go
 */
package storage

import (
	"context"
	"database/sql"
	"dev1024/src/entities"
)

type VideoDao interface {
	Create(ctx context.Context, upLoadVideoReq *entities.VideoInfo) (uint64, error)

	GetVKeyByID(ctx context.Context, vid uint64) (string, error)

	GetUploadCompleteByID(ctx context.Context, vid uint64) (uint64, int, error)

	VideoIndexSearch(ctx context.Context, key string, offset int, size int) (map[string]interface{}, error)

	GetBy(ctx context.Context, keys []string, values []any, by string, byValue any) error

	UpdateBy(ctx context.Context, keys []string, values []any, by string, byValue any) error

	VideoIndexAdd(ctx context.Context, videoinfo *entities.VideoInfo) error

	VideoIndexLikesInc(ctx context.Context, vid uint64, inc int) error

	UserLikesIndexAdd(ctx context.Context, vid uint64, uid uint64, likeTime int64) error

	UserLikesIndexDel(ctx context.Context, vid uint64, uid uint64) error

	UserTagsIndexAdd(ctx context.Context, uid uint64) error

	UserTagsIndexUpdate(ctx context.Context, uid uint64, tags []entities.TagLikes) error

	UserTagsIndexDel(ctx context.Context, uid uint64, tags []entities.TagLikes) error

	SelectWhereIn(ctx context.Context, table string, keys []string, in string, values []any) (*sql.Rows, error)
}

type VideoDaoBase struct {
	instance VideoDao
}

// SelectWhereIn implements VideoDao.
func (t*VideoDaoBase) SelectWhereIn(ctx context.Context, table string, keys []string, in string, values []any) (*sql.Rows, error) {
	return t.instance.SelectWhereIn(ctx,table,keys,in,values)
}

// VideoIndexLikesInc implements VideoDao.
func (t *VideoDaoBase) VideoIndexLikesInc(ctx context.Context, vid uint64, inc int) error {
	return t.instance.VideoIndexLikesInc(ctx, vid, inc)
}

// UserTagsIndexAdd implements VideoDao.
func (t *VideoDaoBase) UserTagsIndexAdd(ctx context.Context, uid uint64) error {
	return t.instance.UserTagsIndexAdd(ctx, uid)
}

// UserTagsIndexDel implements VideoDao.
func (t *VideoDaoBase) UserTagsIndexDel(ctx context.Context, uid uint64, tags []entities.TagLikes) error {
	return t.instance.UserTagsIndexDel(ctx, uid, tags)
}

// UserTagsIndexUpdate implements VideoDao.
func (t *VideoDaoBase) UserTagsIndexUpdate(ctx context.Context, uid uint64, tags []entities.TagLikes) error {
	return t.instance.UserTagsIndexUpdate(ctx, uid, tags)
}

// UserLikesIndexAdd implements VideoDao.
func (t *VideoDaoBase) UserLikesIndexAdd(ctx context.Context, vid uint64, uid uint64, likeTime int64) error {
	return t.instance.UserLikesIndexAdd(ctx, vid, uid, likeTime)
}

// UserLikesIndexDel implements VideoDao.
func (t *VideoDaoBase) UserLikesIndexDel(ctx context.Context, vid uint64, uid uint64) error {
	return t.instance.UserLikesIndexDel(ctx, vid, uid)
}

// UpdateBy implements VideoDao.
func (t *VideoDaoBase) UpdateBy(ctx context.Context, keys []string, values []any, by string, byValue any) error {
	return t.instance.UpdateBy(ctx, keys, values, by, byValue)
}

// VideoIndexAdd implements VideoDao.
func (t *VideoDaoBase) VideoIndexAdd(ctx context.Context, videoinfo *entities.VideoInfo) error {
	return t.instance.VideoIndexAdd(ctx, videoinfo)
}

// GetBy implements VideoDao.
func (t *VideoDaoBase) GetBy(ctx context.Context, keys []string, values []any, by string, byValue any) error {
	return t.instance.GetBy(ctx, keys, values, by, byValue)
}

// VideoIndexSearch implements VideoDao.
func (t *VideoDaoBase) VideoIndexSearch(ctx context.Context, key string, offset int, size int) (map[string]interface{}, error) {
	return t.instance.VideoIndexSearch(ctx, key, offset, size)
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
