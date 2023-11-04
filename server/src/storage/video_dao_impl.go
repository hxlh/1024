/*
 * @Date: 2023-10-25 05:59:02
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-03 14:00:34
 * @FilePath: /1024/server/src/storage/video_dao_impl.go
 */
package storage

import (
	"context"
	"database/sql"
	"dev1024/src/entities"
	"dev1024/src/storage/object"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type VideoDaoImpl struct {
	objectStorage object.ObjectStorage
}

// UserTagsIndexUpdate implements VideoDao.
func (*VideoDaoImpl) UserTagsIndexAdd(ctx context.Context, uid uint64) error {
	pool := ctx.Value("elastic.pool").(*sync.Pool)
	client := pool.Get().(*elasticsearch.Client)
	defer pool.Put(client)

	body := fmt.Sprintf(`
	{
		"uid": %v,
		"tags": {}
	}
	`,
		uid,
	)

	resp, err := client.Index("user_tags", strings.NewReader(body), client.Index.WithDocumentID(fmt.Sprintf("%v", uid)))
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	return nil
}

// UserTagsIndexAdd implements VideoDao.
func (t *VideoDaoImpl) UserTagsIndexUpdate(ctx context.Context, uid uint64, tags []entities.TagLikes) error {
	pool := ctx.Value("elastic.pool").(*sync.Pool)
	client := pool.Get().(*elasticsearch.Client)
	defer pool.Put(client)

	builder := strings.Builder{}
	builder.WriteString("\"")
	builder.WriteString("tags." + tags[0].Tag)
	builder.WriteString("\"")
	for i := 1; i < len(tags); i++ {
		builder.WriteString(",")
		builder.WriteString("\"")
		builder.WriteString("tags." + tags[i].Tag)
		builder.WriteString("\"")
	}

	// 查询原来的值，再递增
	body := fmt.Sprintf(`
	{
		"query": {
			"term":{
				"uid":%v
			}
		},
		"_source": [
			%v
		]
	}
	`, uid, builder.String())

	resp, err := client.Search(client.Search.WithIndex("user_tags"), client.Search.WithBody(strings.NewReader(body)))
	if err != nil {
		return errors.WithStack(err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return errors.WithStack(err)
	}
	resp.Body.Close()

	res := gjson.GetBytes(data, "hits")
	num := res.Get("total.value").Int()
	if num == 0 {
		// 没有这个用户文档
		err = t.UserTagsIndexAdd(ctx, uid)
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		source := res.Get("hits").Array()[0].Get("_source")
		tagsSource := source.Get("tags")
		for i := 0; i < len(tags); i++ {
			tmp := tagsSource.Get(tags[i].Tag)
			if !tmp.Exists() {
				continue
			}
			tags[i].Likes += int(tmp.Int())
		}
	}

	// 写回
	builder = strings.Builder{}
	builder.WriteString("\"")
	builder.WriteString(tags[0].Tag)
	builder.WriteString("\"")
	builder.WriteString(":")
	builder.WriteString(fmt.Sprintf("%v", tags[0].Likes))
	for i := 1; i < len(tags); i++ {
		builder.WriteString(",")
		builder.WriteString("\"")
		builder.WriteString(tags[i].Tag)
		builder.WriteString("\"")
		builder.WriteString(":")
		builder.WriteString(fmt.Sprintf("%v", tags[i].Likes))
	}
	body = fmt.Sprintf(`
	{
		"doc":{
			"tags":{
				%v
			}
		}
	}
	`, builder.String())

	resp, err = client.Update("user_tags", fmt.Sprintf("%v", uid), strings.NewReader(body))
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	return nil
}

// UserTagsIndexDel implements VideoDao.
func (t *VideoDaoImpl) UserTagsIndexDel(ctx context.Context, uid uint64, tags []entities.TagLikes) error {
	return t.UserTagsIndexUpdate(ctx, uid, tags)
}

// UserLikesIndexAdd implements VideoDao.
func (t *VideoDaoImpl) UserLikesIndexAdd(ctx context.Context, vid uint64, uid uint64, likeTime int64) error {
	pool := ctx.Value("elastic.pool").(*sync.Pool)
	client := pool.Get().(*elasticsearch.Client)
	defer pool.Put(client)

	body := fmt.Sprintf(`
	{
		"vid": %v,
		"uid": %v,
		"like_time": %v
	}
	`,
		vid,
		uid,
		time.Now().Unix(),
	)

	resp, err := client.Index("user_likes", strings.NewReader(body))
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()
	return nil
}

// UserLikesIndexDel implements VideoDao.
func (*VideoDaoImpl) UserLikesIndexDel(ctx context.Context, vid uint64, uid uint64) error {
	pool := ctx.Value("elastic.pool").(*sync.Pool)
	client := pool.Get().(*elasticsearch.Client)
	defer pool.Put(client)

	body := fmt.Sprintf(`
	{
		"query": {
		  "bool": {
			"must": [
			  { "term": { "vid": %v }},
			  { "term": { "uid": %v }}
			]
		  }
		}
	}
	`,
		vid,
		uid,
	)
	resp, err := client.DeleteByQuery([]string{"user_likes"}, strings.NewReader(body))
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()
	return nil
}

// UpdateBy implements VideoDao.
func (*VideoDaoImpl) UpdateBy(ctx context.Context, keys []string, values []any, by string, byValue any) error {
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return err
	}
	err = UpdateTableBy(tx, "video1024.video_info", keys, values, by, byValue)

	return errors.WithStack(err)
}

// VideoIndexAdd implements VideoDao.
func (t *VideoDaoImpl) VideoIndexAdd(ctx context.Context, videoinfo *entities.VideoInfo) error {
	pool := ctx.Value("elastic.pool").(*sync.Pool)
	client := pool.Get().(*elasticsearch.Client)
	defer pool.Put(client)

	body := fmt.Sprintf(`
	{
		"vid": %v,
		"uploader": %v,
		"subtitled": "%v",
		"tags": "%v",
		"likes": %v,
		"upload_time": %v
	}
	`, videoinfo.Vid,
		videoinfo.UpLoader,
		videoinfo.Subtitled,
		videoinfo.Tags,
		videoinfo.Likes,
		videoinfo.UpLoadTime)
	resp, err := client.Index("videoinfo", strings.NewReader(body))
	if err != nil {
		return errors.WithStack(err)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetBy implements VideoDao.
func (*VideoDaoImpl) GetBy(ctx context.Context, keys []string, values []any, by string, byValue any) error {
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return err
	}
	return errors.WithStack(SelectTableBy(tx, "video1024.video_info", keys, values, by, byValue))
}

// VideoIndexSearch implements VideoDao.
func (*VideoDaoImpl) VideoIndexSearch(ctx context.Context, key string, offset int, size int) (map[string]interface{}, error) {
	pool := ctx.Value("elastic.pool").(*sync.Pool)
	client := pool.Get().(*elasticsearch.Client)
	defer pool.Put(client)

	body := fmt.Sprintf(`
	{
		"query": {
		  "bool": {
			"should": [
			  { "match": { "subtitled": "%v" } },
			  { "match": { "tags": "%v" } }
			]
		  }
		},
		"sort": [
		  { "likes": "desc" },
		  { "upload_time": "desc" }
		],
		"from": %v,   
		"size": %v,
		"highlight": {
			"fields": {
			  "subtitled": {}
			}
		}
	}
	`, key, key, offset, size)

	resp, err := client.Search(
		client.Search.WithIndex("videoinfo"),
		client.Search.WithBody(strings.NewReader(body)),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer resp.Body.Close()
	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	err = json.Unmarshal(buff, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateComplete implements VideoDao.
func (t *VideoDaoImpl) UpdateComplete(ctx context.Context, vid uint64, complete int) error {
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return err
	}
	return t.updateComplete(tx, vid, complete)
}

func (t *VideoDaoImpl) updateComplete(tx *sql.Tx, vid uint64, complete int) error {
	err := UpdateTableBy(tx, "video1024.video_info", []string{"upload_complete"}, []any{complete}, "vid", vid)
	return errors.WithStack(err)
}

// GetUploadCompleteByID implements VideoDao.
func (t *VideoDaoImpl) GetUploadCompleteByID(ctx context.Context, vid uint64) (uint64, int, error) {
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return 0, 0, err
	}
	return t.getUploadCompleteByID(tx, vid)
}

func (t *VideoDaoImpl) getUploadCompleteByID(tx *sql.Tx, vid uint64) (uint64, int, error) {
	var uploader uint64
	var uploadComplete int
	err := SelectTableBy(tx, "video1024.video_info", []string{"uploader", "upload_complete"}, []any{&uploader, &uploadComplete}, "vid", vid)
	if err != nil {
		return 0, 0, errors.WithStack(err)
	}
	return uploader, uploadComplete, nil
}

// GetVKeyByID implements VideoDao.
func (t *VideoDaoImpl) GetVKeyByID(ctx context.Context, vid uint64) (string, error) {
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return "", err
	}
	return t.getVKeyByID(tx, vid)
}

func (t *VideoDaoImpl) getVKeyByID(tx *sql.Tx, vid uint64) (string, error) {
	var vkey sql.NullString
	err := SelectTableBy(tx, "video1024.video_info", []string{"vkey"}, []any{&vkey}, "vid", vid)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if !vkey.Valid {
		return "", errors.New("vkey not exist")
	}
	return vkey.String, nil
}

func (t *VideoDaoImpl) Create(ctx context.Context, videoInfo *entities.VideoInfo) (uint64, error) {
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return 0, err
	}
	vid, err := t.create(tx, videoInfo)
	return vid, err
}

func (t *VideoDaoImpl) create(tx *sql.Tx, videoInfo *entities.VideoInfo) (uint64, error) {
	vid, err := InsertTableWith(tx, "video1024.video_info", []string{
		"uploader", "subtitled", "likes", "tags", "upload_complete", "upload_time",
	}, []any{
		videoInfo.UpLoader,
		videoInfo.Subtitled,
		videoInfo.Likes,
		videoInfo.Tags,
		videoInfo.UpLoadComplete,
		time.Now().UnixMilli(),
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}
	videoInfo.Vid = uint64(vid)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return videoInfo.Vid, nil
}

// UpdateVKey implements VideoDao.
func (t *VideoDaoImpl) UpdateVKey(ctx context.Context, videoInfo *entities.VideoInfo) error {
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return err
	}
	err = t.updateVKey(tx, videoInfo.Vid, videoInfo.VKey)
	return err
}

func (t *VideoDaoImpl) updateVKey(tx *sql.Tx, vid uint64, vkey string) error {
	stmt, err := tx.Prepare("UPDATE video1024.video_info SET vkey=? WHERE vid = ?")
	if err != nil {
		return errors.WithStack(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(vid, vkey)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func NewVideoDaoImpl() *VideoDaoImpl {
	return &VideoDaoImpl{}
}

var _ VideoDao = (*VideoDaoImpl)(nil)
