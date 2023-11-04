<template>
  <div class="home-page">
    <div class="top-nav">
      <router-link to="/" class="nav-link front_page">首页</router-link>
      <router-link to="/videos" class="nav-link">推荐</router-link>
      <input  class="search-input" placeholder="搜索..." />
      <router-link to="/upload">
        <img src="../assets/icons/加号.svg" alt="上传" class="upload-product" />
      </router-link>
      <div class="user" @click="handleUserClick">
        <img src="../assets/icons/用户.svg" alt="用户图像" class="user-avatar" />
        <span class="user-text" >{{ username ? username : '登录' }}</span>
      <LoginModule :showModal="showModal"  @close="closeModal" @loggedIn="updateLoginStatus" />
      </div>
    </div>
    <router-view></router-view>
  </div>
</template>

<script  setup>
import { ref } from "vue";
import LoginModule from "@/common/user/LoginModule.vue";
import RegisterModule from "@/common/user/RegisterModule.vue";
import router from "@/router";

const showModal = ref(false); //展示登录框
const username = localStorage.getItem("username")
const isLogged= ref(false)
function showLogin() {
  console.log("showLogin")
  showModal.value = true; //为true 确保显示登录模态框

}

function updateLoginStatus(status){
  isLogged.value = status;
}
function handleUserClick(){
  console.log("handleUserClick")
  console.log(isLogged.value)
  if(isLogged.value){
    router.push('/userinfo');

  }else{
    showModal.value = true;
  }
}
function closeModal() {
  showModal.value = false;
}





</script>


<style scoped lang="scss">
.home-page {
  //display: flex;
  flex-direction: column;
  align-items: center;

  .top-nav {
    z-index: 1000;
    width: 100%;
    background-color: #181818;
    color: white;
    display: flex;
    justify-content: space-between; /* 使子元素在水平方向上均匀分布 */
    align-items: center; /* 使子元素在垂直方向上居中对齐 */
    padding: 10px;
    position: absolute;
    top: 0;
    font-size: 24px;
    text-decoration: none;

    .nav-link {
      color: white;
      text-decoration: none;
    }
    .front_page{
      margin-left: 20%;
    }

    .search-input {
      width: 20%;
      padding: 8px;
      margin-left: 10px;
      border: 1px solid white;
      background-color: #181818;
      color: white;
      border-radius: 5px;
      outline: none;
      font-size: 16px;
    }
    .search-input:focus {
      border-color: white !important; /* 设置激活时的边框颜色为白色 */
      box-shadow: 0 0 5px rgba(255, 255, 255, 0.5) !important; /* 设置激活时的阴影效果 */
    }
    .user {
      width: 120px;
      height: 50px;
      background-color: rgb(254, 44, 85);
      display: flex; /* 使用Flex布局 */
      align-items: center; /* 垂直居中 */
      border-radius: 15px;
      .user-avatar {
        width: 30px;
        height: 30px;
        object-fit: cover;
        border-radius: 50%;
        margin-left: 10px;
      }
      .user-text{
        font-size: 15px;
        margin-left: 10px;
      }
    }

    .user:hover{
      cursor: pointer; /* 当鼠标悬停在元素上时显示小手光标 */
    }


    .upload-product {
      width: 40px; /* 修改上传图像的宽度 */
      height: 40px; /* 修改上传图像的高度 */
      object-fit: cover;
      margin-right: 10px; /* 调整上传图像与输入框之间的间距 */
      background-color: rgb(254, 44, 85);
      /* 添加阴影 */
      box-shadow: 0 4px 8px rgba(254, 44, 85, 0.1); /* x偏移量, y偏移量, 模糊半径, 阴影颜色的透明度 */
    }

    .upload-product:hover {
      content: url("../assets/icons/加号1.svg");
      cursor: pointer; /* 当鼠标悬停在元素上时显示小手光标 */
      background-color: rgb(254, 44, 85);
      /* 添加阴影 */
      box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2); /* x偏移量, y偏移量, 模糊半径, 阴影颜色的透明度 */
    }



    /* 输入框激活时的样式 */
    .search-input:focus {
      border-color: #3498db;
      box-shadow: 0 0 5px rgba(52, 152, 219, 0.5);
    }
  }
}
</style>
