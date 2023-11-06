import api from './api'

export function likeVideo(vid,uid) {
  const requestBody = {
    vid:vid,
    uid:uid
  }
  return api
    .post(`/video/like`,requestBody)
    .then(response=>{
    console.log("like中的",response.data)
    return response.data
  })
}

export function unLikeVideo(vid,uid){
  const requestBody = {
    vid:vid,
    uid:uid
  }
  return api
    .post(`/video/cancel_like`,requestBody)
    .then(response=>{
      console.log("unlike中的",response.data)
      return response.data
  })

}
