/*
 * @Date: 2023-10-25 05:33:57
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-05 17:33:32
 * @FilePath: /1024/server/src/service/video_service_impl.go
 */
package service

import (
	"context"
	"dev1024/src/entities"
	"dev1024/src/storage"
	"dev1024/src/storage/object"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
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
	// 2.查询该视频对应的tags
	// 3.更新Elasticsearch的user_tags表的tags
	// 4.更新数据库中的likes
	// 5.更新videoinfo索引中的likes
	var resp entities.CancelLikeVideoResp
	var tags string
	var likes int64

	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return resp, err
	}
	err = t.videoDao.GetBy(ctx, []string{"tags", "likes"}, []any{&tags, &likes}, "vid", req.Vid)
	defer tx.Commit()
	if err != nil {
		return resp, err
	}
	if tags != "" {
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
	}
	err = t.videoDao.UserLikesIndexDel(ctx, req.Vid, req.Uid)
	if err != nil {
		return resp, err
	}

	err = t.videoDao.UpdateBy(ctx, []string{"likes"}, []any{likes - 1}, "vid", req.Vid)
	if err != nil {
		return resp, err
	}

	err = t.videoDao.VideoIndexLikesInc(ctx, req.Vid, -1)
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
	// 4.更新数据库中的likes
	// 5.更新videoinfo索引中的likes
	var resp entities.LikeVideoResp
	err := t.videoDao.UserLikesIndexAdd(ctx, req.Vid, req.Uid, time.Now().UnixMilli())
	if err != nil {
		return resp, err
	}
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return resp, err
	}
	defer tx.Commit()
	var tags string
	var likes int64
	err = t.videoDao.GetBy(ctx, []string{"tags", "likes"}, []any{&tags, &likes}, "vid", req.Vid)
	if err != nil {
		return resp, err
	}

	if tags != "" {
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
	}

	err = t.videoDao.UpdateBy(ctx, []string{"likes"}, []any{likes + 1}, "vid", req.Vid)
	if err != nil {
		return resp, err
	}
	err = t.videoDao.VideoIndexLikesInc(ctx, req.Vid, 1)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// RecommendedVideo implements VideoService.
func (t *VideoServiceImpl) RecommendedVideo(ctx context.Context, req entities.RecommendedVideoReq) (entities.RecommendedVideoResp, error) {
	var resp entities.RecommendedVideoResp
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	// 获取最大最小时间戳
	videoinfoMinTimestamp, videoinfoMaxTimestamp, err := t.searchTimeByIndex(ctx, "videoinfo", "upload_time")
	if err != nil {
		return resp, err
	}
	userLikesMinTimestamp, userLikesMaxTimestamp, err := t.searchTimeByIndex(ctx, "user_likes", "like_time")
	if err != nil {
		return resp, err
	}

	videoinfoMaxTime := time.UnixMilli(videoinfoMaxTimestamp)
	videoinfoMinTime := time.UnixMilli(videoinfoMinTimestamp)
	userLikesMaxTime := time.UnixMilli(userLikesMaxTimestamp)
	userLikesMinTime := time.UnixMilli(userLikesMinTimestamp)

	// for strategy1,strategy2,strategy4,strategy5
	// 获取一共有几个月,几个15天
	var videoinfoDur time.Duration
	var videoinfoDurMonths int
	var videoinfoDur15Days int
	videoinfoDur = videoinfoMaxTime.Sub(videoinfoMinTime)
	videoinfoDurMonths = int(videoinfoDur.Hours() / 24 / 30)
	videoinfoDur15Days = min(int(videoinfoDur.Hours()/24/15), 6)

	// for strategy3
	var userLikesDur time.Duration
	var userLikesDurHours int
	if userLikesMinTimestamp != -1 {
		userLikesDur = userLikesMaxTime.Sub(userLikesMinTime)
		userLikesDurHours = min(int(userLikesDur.Hours()/3), (30*24)/3)
	}

	strategy1SelectedMonth := randomSelect(0, videoinfoDurMonths)
	strategy2Selected15Days := randomSelect(0, videoinfoDur15Days)
	strategy3Selected3Hours := randomSelect(0, userLikesDurHours)
	strategy4SelectedMonth := randomSelect(0, videoinfoDurMonths)
	strategy5Selected15Days := randomSelect(0, videoinfoDur15Days)

	strategy1End := videoinfoMaxTime.AddDate(0, -strategy1SelectedMonth, 0).UnixMilli()
	strategy1Start := videoinfoMaxTime.AddDate(0, -(strategy1SelectedMonth + 1), 0).UnixMilli()

	strategy2End := videoinfoMaxTime.AddDate(0, 0, -(strategy2Selected15Days * 15)).UnixMilli()
	strategy2Start := videoinfoMaxTime.AddDate(0, 0, -((strategy2Selected15Days + 1) * 15)).UnixMilli()

	strategy3End := userLikesMaxTime.AddDate(0, 0, -(strategy3Selected3Hours * 3)).UnixMilli()
	strategy3Start := userLikesMaxTime.AddDate(0, 0, -((strategy3Selected3Hours + 1) * 3)).UnixMilli()

	strategy4End := videoinfoMaxTime.AddDate(0, -strategy4SelectedMonth, 0).UnixMilli()
	strategy4Start := videoinfoMaxTime.AddDate(0, -(strategy4SelectedMonth + 1), 0).UnixMilli()

	strategy5End := videoinfoMaxTime.AddDate(0, 0, -(strategy5Selected15Days * 15)).UnixMilli()
	strategy5Start := videoinfoMaxTime.AddDate(0, 0, -((strategy5Selected15Days + 1) * 15)).UnixMilli()

	baseSize := 30
	strategy1Rate := 0.25
	strategy2Rate := 0.15
	strategy3Rate := 0.20
	strategy4Rate := 0.25
	strategy5Rate := 0.15

	res := make([]uint64, 0)
	// 是否登录都能采用
	r, _ := t.strategy1(ctx, strategy1Start, strategy1End, int(float64(baseSize)*strategy1Rate))
	if len(r) > 0 {
		res = append(res, r...)
	}
	r, _ = t.strategy2(ctx, strategy2Start, strategy2End, int(float64(baseSize)*strategy2Rate))
	if len(r) > 0 {
		res = append(res, r...)
	}
	if userLikesMinTimestamp != -1 {
		r, _ = t.strategy3(ctx, strategy3Start, strategy3End, int(float64(baseSize)*strategy3Rate))
		if len(r) > 0 {
			res = append(res, r...)
		}
	}
	// 已登录才能采用
	if req.Uid != 0 {
		r, err := t.strategy4(ctx, req.Uid, strategy4Start, strategy4End, int(float64(baseSize)*strategy4Rate))
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		if len(r) > 0 {
			res = append(res, r...)
		}
		r, err = t.strategy5(ctx, req.Uid, strategy5Start, strategy5End, int(float64(baseSize)*strategy5Rate))
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		if len(r) > 0 {
			res = append(res, r...)
		}
	}

	if len(res) <= 0 {
		return resp, nil
	}

	// 从数据库查询信息和生成直链
	ctx, tx, err := getTxFromCtx(ctx)
	if err != nil {
		return resp, err
	}
	defer tx.Commit()
	// build vid str
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%v", res[0]))
	for i := 1; i < len(res); i++ {
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%v", res[i]))
	}

	query := fmt.Sprintf(`
	SELECT 
	vid,username,nickname,vkey,thumbnail,subtitled,likes,tags,upload_time
	FROM video1024.video_info t1,account t2
	WHERE vid in (%v) AND t1.uploader =t2.uid 
	`, b.String())

	rows, err := tx.Query(query)
	if err != nil {
		return resp, err
	}
	for rows.Next() {
		info := &entities.SearchVideoRespInfo{}
		err = rows.Scan(&info.Vid, &info.UpLoaderUsername, &info.UpLoaderNickname,
			&info.VideoUrl, &info.ThumbnailUrl, &info.Subtitled,
			&info.Likes, &info.Tags, &info.UpLoadTime,
		)
		if err != nil {
			continue
		}
		resp.Info = append(resp.Info, info)
	}

	// 生成直链
	for i := 0; i < len(resp.Info); i++ {
		info := resp.Info[i]
		info.VideoUrl = t.objectStorage.Load(info.VideoUrl, time.Now().Add(VIDEO_CDN_DEAD_TIME).Unix())
		info.ThumbnailUrl = t.objectStorage.Load(info.ThumbnailUrl, time.Now().Add(VIDEO_CDN_DEAD_TIME).Unix())
	}

	// 已登录则查询点赞状态
	if req.Uid != 0 {
		for i := 0; i < len(resp.Info); i++ {
			info := resp.Info[i]
			isLike, err := t.videoDao.CheckUserLikes(ctx, info.Vid, req.Uid)
			if err != nil {
				continue
			}
			info.IsLike = isLike
		}
	}

	return resp, nil
}

func randomSelect(start int, end int) int {
	rand.Seed(time.Now().UnixMilli())
	return start + rand.Intn((end+1)-start)
}

// 找不到数据将返回 -1,-1,nil
func (t *VideoServiceImpl) searchTimeByIndex(ctx context.Context, index string, field string) (int64, int64, error) {
	var minTime int64
	var maxTime int64

	pool := ctx.Value("elastic.pool").(*sync.Pool)
	client := pool.Get().(*elasticsearch.Client)
	defer pool.Put(client)

	body := fmt.Sprintf(`
	{
		"aggs": {
			"min_time": {
				"min": {
					"field": "%v",
					"format": "#"
				}
			},
			"max_time": {
				"max": {
					"field": "%v",
					"format": "#"
				}
			}
		},
		"size": 0
	}
	`, field, field)
	resp, err := client.Search(client.Search.WithIndex(index), client.Search.WithBody(strings.NewReader(body)))
	if err != nil {
		return minTime, maxTime, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return minTime, maxTime, err
	}

	total := gjson.GetBytes(data, "hits.total.value").Int()
	if total == 0 {
		// 该索引没有数据
		return -1, -1, nil
	}

	minTimeStr := gjson.GetBytes(data, "aggregations.min_time.value_as_string").String()
	maxTimeStr := gjson.GetBytes(data, "aggregations.max_time.value_as_string").String()

	minTime, err = strconv.ParseInt(minTimeStr, 10, 64)
	if err != nil {
		return minTime, maxTime, errors.WithStack(err)
	}
	maxTime, err = strconv.ParseInt(maxTimeStr, 10, 64)
	if err != nil {
		return minTime, maxTime, errors.WithStack(err)
	}

	return minTime, maxTime, nil
}

func (t *VideoServiceImpl) strategy1(ctx context.Context, start, end int64, size int) ([]uint64, error) {
	pool := ctx.Value("elastic.pool").(*sync.Pool)
	client := pool.Get().(*elasticsearch.Client)
	defer pool.Put(client)

	res := make([]uint64, 0)

	query := fmt.Sprintf(`
	{
	  "_source": ["vid"],
	  "query": {
		"range": {
		  "upload_time": {
			"gte": %v,
			"lte": %v
		  }
		}
	  },
	  "sort": [
		{
		  "likes": {
			"order": "desc"
		  }
		}
	  ],
	  "size": %v
	}`, start, end, size)
	resp, err := client.Search(
		client.Search.WithIndex("videoinfo"), // 替换为你的索引名称
		client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	total := gjson.GetBytes(data, "hits.total.value").Int()
	if total == 0 {
		return res, nil
	}

	hitsArr := gjson.GetBytes(data, "hits.hits").Array()
	for i := 0; i < len(hitsArr); i++ {
		item := hitsArr[i]
		vid := item.Get("_source.vid").Uint()
		res = append(res, vid)
	}
	return res, nil
}

func (t *VideoServiceImpl) strategy2(ctx context.Context, start, end int64, size int) ([]uint64, error) {
	pool := ctx.Value("elastic.pool").(*sync.Pool)
	client := pool.Get().(*elasticsearch.Client)
	defer pool.Put(client)

	res := make([]uint64, 0)

	query := fmt.Sprintf(`
	{
	  "_source": ["vid"],
	  "query": {
		"range": {
		  "upload_time": {
			"gte": %v,
			"lte": %v
		  }
		}
	  },
	  "sort": [
		{
		  "likes": {
			"order": "desc"
		  }
		}
	  ],
	  "size": %v
	}`, start, end, size)

	resp, err := client.Search(
		client.Search.WithIndex("videoinfo"), // 替换为你的索引名称
		client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	total := gjson.GetBytes(data, "hits.total.value").Int()
	if total == 0 {
		return res, nil
	}

	hitsArr := gjson.GetBytes(data, "hits.hits").Array()
	for i := 0; i < len(hitsArr); i++ {
		item := hitsArr[i]
		vid := item.Get("_source.vid").Uint()
		res = append(res, vid)
	}
	return res, nil
}

// @return []vid,error
func (t *VideoServiceImpl) strategy3(ctx context.Context, start, end int64, size int) ([]uint64, error) {
	pool := ctx.Value("elastic.pool").(*sync.Pool)
	client := pool.Get().(*elasticsearch.Client)
	defer pool.Put(client)

	res := make([]uint64, 0)

	query := fmt.Sprintf(`
	{
		"query": {
			"range": {
				"like_time": {
					"gte": %v,
					"lte": %v
				}
			}
		},
		"aggs": {
			"vid_count": {
				"terms": {
					"field": "vid",
					"size": %v,
					"order": {
						"_count": "desc"
					}
				}
			}
		},
		"size": 0
	}
	`, start, end, size)

	resp, err := client.Search(
		client.Search.WithIndex("user_likes"), // 替换为你的索引名称
		client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	total := gjson.GetBytes(data, "hits.total.value").Int()
	if total == 0 {
		return res, nil
	}

	buckets := gjson.GetBytes(data, "aggregations.vid_count.buckets").Array()
	for i := 0; i < len(buckets); i++ {
		item := buckets[i]
		vid := item.Get("key").Uint()
		res = append(res, vid)
	}

	return res, nil
}

func (t *VideoServiceImpl) strategy4(ctx context.Context, uid uint64, start, end int64, size int) ([]uint64, error) {
	tags, err := t.getUserMostLikeTags(ctx, uid)
	if err != nil {
		return nil, err
	}

	res := make([]uint64, 0)

	if len(tags) <= 0 {
		return res, nil
	}

	// 构造tags
	b := strings.Builder{}
	b.WriteString(tags[0])
	for i := 1; i < len(tags); i++ {
		b.WriteString(",")
		b.WriteString(tags[i])
	}

	pool := ctx.Value("elastic.pool").(*sync.Pool)
	client := pool.Get().(*elasticsearch.Client)
	defer pool.Put(client)

	query := fmt.Sprintf(`
	{
		"_source": ["vid"],
		"query": {
			"bool": {
				"must": [
					{
						"match": {
							"tags": "%v"
						}
					},
					{
						"range": {
							"upload_time": {
								"gte": %v,
								"lte": %v
							}
						}
					}
				]
			}
		},
		"sort": [
			{
				"_score": {
					"order": "desc"
				}
			},
			{
				"likes": {
					"order": "desc"
				}
			}
		],
		"size": %v
	}
	`, b.String(), start, end, size)

	resp, err := client.Search(
		client.Search.WithIndex("videoinfo"), // 替换为你的索引名称
		client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	total := gjson.GetBytes(data, "hits.total.value").Int()
	if total == 0 {
		return res, nil
	}

	hitsArr := gjson.GetBytes(data, "hits.hits").Array()
	for i := 0; i < len(hitsArr); i++ {
		item := hitsArr[i]
		vid := item.Get("_source.vid").Uint()
		res = append(res, vid)
	}

	return res, nil
}

func (t *VideoServiceImpl) strategy5(ctx context.Context, uid uint64, start, end int64, size int) ([]uint64, error) {
	tags, err := t.getUserMostLikeTags(ctx, uid)
	if err != nil {
		return nil, err
	}
	res := make([]uint64, 0)

	if len(tags) <= 0 {
		return res, nil
	}

	// 构造tags
	b := strings.Builder{}
	b.WriteString(tags[0])
	for i := 1; i < len(tags); i++ {
		b.WriteString(",")
		b.WriteString(tags[i])
	}

	pool := ctx.Value("elastic.pool").(*sync.Pool)
	client := pool.Get().(*elasticsearch.Client)
	defer pool.Put(client)

	query := fmt.Sprintf(`
	{
		"_source": ["vid"],
		"query": {
			"bool": {
				"must": [
					{
						"match": {
							"tags": "%v"
						}
					},
					{
						"range": {
							"upload_time": {
								"gte": %v,
								"lte": %v
							}
						}
					}
				]
			}
		},
		"sort": [
			{
				"_score": {
					"order": "desc"
				}
			},
			{
				"likes": {
					"order": "asc"
				}
			}
		],
		"size": %v
	}
	`, b.String(), start, end, size)

	resp, err := client.Search(
		client.Search.WithIndex("videoinfo"), // 替换为你的索引名称
		client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	total := gjson.GetBytes(data, "hits.total.value").Int()
	if total == 0 {
		return res, nil
	}

	hitsArr := gjson.GetBytes(data, "hits.hits").Array()
	for i := 0; i < len(hitsArr); i++ {
		item := hitsArr[i]
		vid := item.Get("_source.vid").Uint()
		res = append(res, vid)
	}

	return res, nil
}

func (t *VideoServiceImpl) getUserMostLikeTags(ctx context.Context, uid uint64) ([]string, error) {
	pool := ctx.Value("elastic.pool").(*sync.Pool)
	client := pool.Get().(*elasticsearch.Client)
	defer pool.Put(client)

	res := make([]string, 0)

	query := fmt.Sprintf(`
	{
		"query":{
			"term":{
				"uid":%v
			}
		},
		"_source":[
			"tags"
		]
	}`, uid)

	resp, err := client.Search(
		client.Search.WithIndex("user_tags"), // 替换为你的索引名称
		client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	total := gjson.GetBytes(data, "hits.total.value").Int()
	if total == 0 {
		return res, nil
	}

	arr := make([]struct {
		Key string
		Val int
	}, 0)

	tags := gjson.GetBytes(data, "hits.hits").Array()[0].Get("_source.tags")
	tags.ForEach(func(key, value gjson.Result) bool {
		var kv struct {
			Key string
			Val int
		}
		kv.Key = key.String()
		kv.Val = int(value.Int())
		arr = append(arr, kv)
		return true
	})

	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Val > arr[i].Val
	})

	for i := 0; i < 5; i++ {
		if i >= len(arr) {
			break
		}
		res = append(res, arr[i].Key)
	}

	return res, nil
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
	if !ok {
		return res, nil
	}
	total := hits["total"].(map[string]interface{})
	hitNum := int(total["value"].(float64))
	if hitNum == 0 {
		return res, nil
	}

	uids := make([]uint64, 0)
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
		uids = append(uids, uploaderID)

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

	for i := 0; i < len(res.Info); i++ {
		info := res.Info[i]
		isLike, err := t.videoDao.CheckUserLikes(ctx, info.Vid, uids[i])
		if err != nil {
			return res, err
		}
		info.IsLike = isLike
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
