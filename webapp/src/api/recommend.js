import api from './api'

export function recommendVideo(uid){
  // 创建一个params对象，如果uid存在，则添加到params中
  const params = uid ? { uid: uid } : {};
  return api.get('/video/recommended', { params }).then(response => {
    if(response.data.status === 'ok'){
      return response.data;
    } else {
      // 可以处理错误情况或返回一个错误标志/消息
      throw new Error('获取推荐视频失败！');
    }
  }).catch(error => {
    console.error('获取推荐视频发送错误', error);
    // 可以根据需求处理错误，比如返回null或特定错误信息
    return null;
  });
}
