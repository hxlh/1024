// 引入你的api实例
import api from './api';

// 使用相对路径和api实例来发送post请求
export function searchVideo(requestBody) {
  return api.post('/video/search', requestBody)
    .then(response => {
      // 处理成功的响应
      console.log("response.data:",response.data);
      return response;

    })
    .catch(error => {
      // 处理错误
      console.error('Error during the API call', error);
    });
}

