/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2023-10-24 13:44:28
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-29 13:08:42
 * @FilePath: /1024/server/src/service/video_service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import "dev1024/src/entities"

type VideoService interface {
	FetchNextByVID(vid int64) (*entities.VideoInfo,string, error)
	// @return 
	// vid: 上传文件的key
	// token: 上传token
	ApplyVideoUpToken(video *entities.VideoInfo) (string, error)
}
