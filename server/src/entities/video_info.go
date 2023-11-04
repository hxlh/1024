/*
 * @Date: 2023-10-25 05:54:39
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-28 11:34:29
 * @FilePath: /1024/server/src/entities/video_info.go
 */
package entities

type VideoInfo struct {
	Vid       uint64  `json:"vid"`
	UpLoader  int64  `json:"uploader"`
	CDN       string `json:"cdn"`
	Subtitled string `json:"subtitled"`
	Likes     int64  `json:"likes"`
	Tags      string `json:"tags"`
}
