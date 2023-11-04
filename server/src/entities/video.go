/*
 * @Date: 2023-10-25 05:54:39
 * @LastEditors: hxlh
 * @LastEditTime: 2023-11-03 08:32:17
 * @FilePath: /1024/server/src/entities/video.go
 */
package entities

type VideoInfo struct {
	Vid            uint64 `json:"vid"`             // 视频id
	UpLoader       uint64 `json:"uploader"`        // 上传者id
	VKey           string `json:"vkey"`            // 视频对象存储的key
	Thumbnail      string `json:"thumbnail"`       // 封面对象存储的key
	Subtitled      string `json:"subtitled"`       // 视频标题
	Likes          int64  `json:"likes"`           // 该视频点赞数
	Tags           string `json:"tags"`            // 视频的tags
	UpLoadTime     int64  `json:"upload_time"`     // 上传的时间,毫秒时间戳
	UpLoadComplete int    `json:"upload_complete"` // 上传是否完成，因为由客户端上传，服务端不到何时上传完成，需要有一个回调通知
}

type ApplyVideoUpTokenResp struct {
	UploadKey string `json:"upload_key"`
	Token     string `json:"token"`
}

type UpLoadVideoReq struct {
	UpLoader  uint64 `json:"uploader"`
	Subtitled string `json:"subtitled"`
	Tags      string `json:"tags"`
}

type UpLoadVideoResp struct {
	Vid   uint64 `json:"vid"`
	VKey  string `json:"vkey"`
	Token string `json:"token"`
}

type GetVideoByIDResp struct {
	Url string `json:"url"`
}

type UpLoadVideoCallBackReq struct {
	Vid uint64
}

type SearchVideoReq struct {
	Key    string `json:"key"`
	Offset int    `json:"offset"` //分页
}

type SearchVideoRespInfo struct {
	Vid                uint64 `json:"vid"`                 // 视频id
	UpLoaderNickname   string `json:"uploader_nickname"`   // 上传者的别名
	UpLoaderUsername   string `json:"uploader_username"`   // 上传者的别名
	VideoUrl           string `json:"video"`               // 视频对象存储的直链
	ThumbnailUrl       string `json:"thumbnail"`           // 封面对象存储的直链
	Subtitled          string `json:"subtitled"`           // 视频标题
	HighLightSubtitled string `json:"highlight_subtitled"` // 设置了高亮的视频标题
	Likes              int64  `json:"likes"`               // 该视频点赞数
	Tags               string `json:"tags"`                // 视频的tags
	UpLoadTime         int64  `json:"upload_time"`         // 上传的时间，毫秒时间戳
}

type SearchVideoResp struct {
	Info []*SearchVideoRespInfo `json:"info"` //结果
}

type LikeVideoReq struct {
	Vid uint64 `json:"vid"`
	Uid uint64 `json:"uid"`
}

type LikeVideoResp struct {
}

type CancelLikeVideoReq struct {
	Vid uint64 `json:"vid"`
	Uid uint64 `json:"uid"`
}

type CancelLikeVideoResp struct {
}

type TagLikes struct {
	Tag   string
	Likes int
}

type RecommendedVideoReq struct {
	Uid uint64 `json:"uid"`
}

type RecommendedVideoResp struct {
	Info []*SearchVideoRespInfo `json:"info"` //结果
}
