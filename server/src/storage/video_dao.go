/*
 * @Date: 2023-10-24 16:29:27
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-26 05:24:28
 * @FilePath: /1024/server/src/storage/video_dao.go
 */
package storage

import (
	"dev1024/src/entities"

	_ "github.com/go-sql-driver/mysql"
)

type VideoDao interface {
	// 获取比vid大的下一个视频的信息
	GetNextInfo(vid int) *entities.VideoInfo
}
