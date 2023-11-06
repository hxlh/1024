/*
 * @Date: 2023-10-25 05:59:02
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-26 05:24:43
 * @FilePath: /1024/server/src/storage/video_dao_impl.go
 */
package storage

import "dev1024/src/entities"

type VideoDaoImpl struct {
}

var _ VideoDao = (*VideoDaoImpl)(nil)

func (vd *VideoDaoImpl) GetNextInfo(vid int) *entities.VideoInfo {
	return nil
}
