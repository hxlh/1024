<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div>
    <RouterView />
    <div class="switch-tab">
      <!-- 创建一排按钮 -->
      <button>
        <RouterLink to="/">首页</RouterLink>
      </button>
      <button>
        <RouterLink to="/home">我的</RouterLink>
      </button>
      <!-- <button class="button"></button>
          <button class="button"></button> -->
    </div>
  </div>
</template>

<style scoped>
.switch-tab {
  width: 400px;
  height: 50px;
  display: flex;
  /* 平分空间 */
  justify-content: space-between;
  background-color: black;
}

.switch-tab button {
  background-color: black;
  color: white;
  flex-grow: 1;
  /* 按钮自动占满可用空间 */
  height: 50px;
  padding: 5px 0;
  text-align: center;
  border: none;
  cursor: pointer;
}

.switch-tab a {
  background-color: black;
  color: white;
  flex-grow: 1;
  text-decoration: none;
}
</style>
