<template>
  <div class="profile-card">
    <div class="profile-header">
      <img :src="user.avatar" alt="user avatar" class="avatar">  <!--  这里放用户头像  -->
      <div class="user-info">
        <h3 class="username">{{ user.name }}</h3>
        <p class="user-details">关注 {{ user.following }} | 粉丝 {{ user.followers }}</p>
      </div>
      <button @click="logout" class="logout-button">退出登录</button>

      <!--   其他信息和按钮等   -->
    </div>
    <!--   其他个人信息   -->
  </div>
</template>

<script setup>
import { ref } from 'vue';
import router from "@/router";
import { defineEmits } from 'vue';

const emit = defineEmits(['logged-out']);

const user = ref({
  name: localStorage.getItem('username'),
  avatar: 'path_to_avatar_image',
  followers: 216, //模拟粉丝数量
  following: 29, // 模拟关注数量
});
async function logout() {
  // 清除本地存储中的token
  localStorage.removeItem('token');
  localStorage.removeItem('username'); // 清除用户名
  localStorage.removeItem('isLogged'); // 清除登录状态标志
  emit('logged-out');

  // 等待路由完成跳转
  await router.push('/');
  // 路由跳转完成后刷新页面
  location.reload();
}
</script>

<style scoped>
.profile-card {
  position: relative;
  top: 80px;
  background-color: #f4f7fa;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);
  transition: all 0.3s ease-in-out;
}

.profile-card:hover {
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
}

.profile-header {
  display: flex;
  align-items: center;
  border-bottom: 1px solid #e1e4e8;
  padding-bottom: 20px;
  margin-bottom: 20px;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  border: 2px solid #3498db;
  margin-right: 20px;
  transition: transform 0.3s ease-in-out;
}

.avatar:hover {
  transform: scale(1.05);
}

.username {
  margin: 0;
  color: #333;
  font-weight: bold;
  font-size: 1.2rem;
}

.user-details {
  margin: 4px 0 0;
  font-size: 0.9rem;
  color: #555;
}

.logout-button {
  margin-left: auto;
  padding: 8px 16px;
  font-size: 0.9rem;
  color: #fff;
  background-color: #3498db;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  outline: none;
  transition: background-color 0.3s ease-in-out;
}

.logout-button:hover {
  background-color: #2980b9;
}
</style>
