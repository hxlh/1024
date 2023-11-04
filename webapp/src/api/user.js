import api from './api'
export function login(username, pwd){
  const data = {
    username,
    pwd
  }

  return api.post('/account/login',data).then(res=>{
    if(res.data.status === "ok"){
      console.log('登录成功')
      console.log('token',res.data.data.token)
      localStorage.setItem('token',res.data.data.token)
      localStorage.setItem('uid',res.data.data.uid)
      return res.data
    }

  }).catch(err=>{
    alert("账号或密码错误，请重试！")
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
  return api.post('/account/register',data).then(res=>{
    return res.data
  }).catch(err=>{
    console.log(err)
    throw new Error(err)
  })
}

