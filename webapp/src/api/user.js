

const baseURL = 'http://47.236.105.41:8080';
import axios from 'axios'

axios.defaults.baseURL = baseURL;


export function login(username, pwd){
  const data = {
    username,
    pwd
  }

  return axios.post('/account/login',data).then(res=>{
    console.log('token',res.data.data.token)
    localStorage.setItem('token',res.data.data.token)
    localStorage.setItem('uid',res.data.data.uid)
    return res.data

  }).catch(err=>{
    console.log(err)
    //抛出错误内容
    throw new Error(err)
  })
}
export function register(username,pwd,nickname) {
  const data = {
    username,
    pwd,
    nickname
  }
  return axios.post('/account/register',data).then(res=>{
    return res.data
  }).catch(err=>{
    console.log(err)
    throw new Error(err)
  })
}

