// api.js
import axios from 'axios';

const baseURL = 'http://47.236.105.41:8080';

// 创建一个axios实例
const api = axios.create({
  baseURL: baseURL,
  withCredentials: true
  //可以在这里添加其他配置，比如headers
});
//如果localStorage中存在token，则设置默认的Authorization header
const token = localStorage.getItem('token');
if (token) {
  api.defaults.headers['Authorization'] = token;
}

export default api;
