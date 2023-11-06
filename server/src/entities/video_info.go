/*
 * @Date: 2023-10-25 05:54:39
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-25 05:57:51
 * @FilePath: /1024/src/server/entities/video_info.go
 */
package entities

type VideoInfo struct{
	Vid int64 `json:"vid"`
	UpLoader string `json:"uploader"`
	CDN string `json:"cdn"`
	Subtitled string `json:"subtitled"`
	Likes int64 `json:"likes"`
	Tags int64 `json:"tags"`
}