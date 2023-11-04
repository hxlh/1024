const baseURL = 'http://47.236.105.41:8080';
import axios from 'axios'

axios.defaults.baseURL = baseURL;
axios.defaults.withCredentials = true;


export function getUploadTokenAndKey(uploaderId, subtitled, tags){
  const headers = {
    'Authorization': localStorage.getItem('token')
  };
  // axios.post返回的是一个promise，所以可以直接return这个调用
  return axios.post('/video/upload', {
    uploader: uploaderId,
    subtitled: subtitled,
    tags: tags
  },{headers}) // 注意这里的括号是结束axios.post的参数列表
    .then(response => { // 处理响应
      console.log(response.data)
      if (response.data.status === 'ok') {
        console.log("response.data.data:",response.data.data)
        console.log('上传凭证为:', response.data.data.token);
        console.log('上传key为:', response.data.data.vkey);
        localStorage.setItem('vid',response.data.data.vid)
        return response; // 最好返回具体的数据以便后续操作
      } else {
        console.error('获取上传凭证失败:', response.data.message);
        throw new Error(response.data.message); // 抛出错误
      }
    })
    .catch(error => {
      console.error('获取上传凭证出现错误：', error);
      throw error; // 向上抛出错误
    });
}


export function callBackUpload(vid){
  return axios.post('/video/upload_callback',
    {
      Vid:vid,
    },
    {
      headers:{
        "Authorization":localStorage.getItem('token')
      }
    })
}

