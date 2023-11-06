/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2023-10-24 13:44:28
 * @LastEditors: hxlh
 * @LastEditTime: 2023-10-25 05:36:15
 * @FilePath: /1024/src/server/service/video_service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

type VideoService interface {
	FetchNext(vid int64) (string, error)
}
