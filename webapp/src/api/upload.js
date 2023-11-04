
import api from './api'
export function getUploadTokenAndKey(uploaderId, subtitled, tags){
  return api.post('/video/upload', {
    uploader: uploaderId,
    subtitled: subtitled,
    tags: tags
  })
    .then(response=>{
      console.log(response.data)
      if (response.data.status === "ok"){
        console.log("response.data.data:",response.data.data)
        console.log('上传凭证为:', response.data.data.token);
        console.log('上传key为:', response.data.data.vkey);
        localStorage.setItem('vid',response.data.data.vid)
        return response;
      }else{
        console.error('获取上传凭证失败:', response.data.message);
        throw new Error(response.data.message);
      }
    }).catch(error=>{
      console.log('获取上传凭证出现错误：', error);
        throw error;
    })
}



export function callBackUpload(vid){
  return api.post('/video/upload_callback',
    {
      Vid:vid,
    },
    {
      headers:{
        "Authorization":localStorage.getItem('token')
      }
    })
}

