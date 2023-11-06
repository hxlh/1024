<template>
  <div class="home-page">
    <div class="top-nav">
      <router-link to="/" class="nav-link recommend">推荐</router-link>

      <div class="nav-link front_page" @mouseenter="showDropdown = true" @mouseleave="showDropdown = false">
        分类
        <div class="dropdown" v-if="showDropdown">
          <button class="dropdown-item" @click="navigateToCategory('游戏，娱乐')">游戏</button>
          <button class="dropdown-item" @click="navigateToCategory('音乐，乐器，唱歌')">音乐</button>
          <button class="dropdown-item" @click="navigateToCategory('体育，运动，舞蹈，跳舞')">运动</button>
          <button class="dropdown-item" @click="navigateToCategory('健康，运动')">健康</button>
          <button class="dropdown-item" @click="navigateToCategory('科技，芯片，华为')">科技</button>
          <button class="dropdown-item" @click="navigateToCategory('二次元，动漫')">二次元</button>

          <!-- ... 其他分类按钮 ... -->
        </div>
      </div>      <div class="search-area">
        <input
          class="search-input"
          placeholder="搜索..."
          v-model="searchKey"
          @keyup.enter="search"
        />
        <button class="search-btn" @click="search">搜索</button>
      </div>
      <router-link to="/upload">
        <img src="../assets/icons/加号.svg" alt="上传" class="upload-product" />
      </router-link>
      <div class="user" @click="handleUserClick">
        <img src="../assets/icons/用户.svg" alt="用户图像" class="user-avatar" />
        <span class="user-text" >{{ username ? username : '登录' }}</span>
      <LoginModule :showModal="showModal"  @close="closeModal" @loggedIn="updateLoginStatus" />
      </div>
    </div>
    <router-view :videos="videoInfo" @logged-out="handleLogout"></router-view>
  </div>
</template>

<script  setup>
import LoginModule from "@/common/user/LoginModule.vue";
import {onMounted, ref} from "vue";
import {searchVideo} from "@/api/search";
import { useRouter } from 'vue-router'
const router = useRouter()
const showModal = ref(false); //展示登录框
const username = ref(localStorage.getItem("username"))
const isLogged= ref(false)
const searchKey = ref(""); //搜索关键字
const offset = ref(0)
const videoInfo = ref([]); // 添加一个用于存放视频数据的响应式引用
const showDropdown = ref(false); // 控制下拉菜单的显示和隐藏

// 搜索方法

async function search() {
  console.log('key:', searchKey.value);
  try {
    const response = await searchVideo({ key: searchKey.value,offset:offset.value }); // 假设 searchVideo 已正确接收对象参数
    console.log("输出的值为:",response.data.data.info)
    videoInfo.value = response.data.data.info;
    // 在存储到localStorage之前，先转换为JSON字符串
    localStorage.setItem("videoSearchInfo", JSON.stringify(videoInfo.value));

// 从localStorage读取时，需要解析JSON字符串
    const storedData = JSON.parse(localStorage.getItem("videoSearchInfo"));
    console.log("Layout下的videoInfo.value：", storedData);
    await router.push({
      name: 'search',
      params: {videos: videoInfo.value}
    });

  } catch (error) {
    alert("请先登录！")
    console.error('Error fetching videos:', error);
  }
}

function handleLogout() {
  // 更新父组件中的用户名显示
  username.value = null;
}

function updateLoginStatus(newUsername){
  isLogged.value = true;
  username.value = newUsername; // 更新用户名
}
function handleUserClick(){
  if(isLogged.value){
    router.push('/userinfo');

  }else{
    showModal.value = true;
  }
}
async function navigateToCategory(category) {
    console.log('key:', category);
    try {
      const response = await searchVideo({ key: category,offset:offset.value }); // 假设 searchVideo 已正确接收对象参数
      console.log("输出的值为:",response.data.data.info)
      videoInfo.value = response.data.data.info;
      // 在存储到localStorage之前，先转换为JSON字符串
      localStorage.setItem("videoSearchInfo", JSON.stringify(videoInfo.value));

// 从localStorage读取时，需要解析JSON字符串
      const storedData = JSON.parse(localStorage.getItem("videoSearchInfo"));
      console.log("Layout下的videoInfo.value：", storedData);
      await router.push({
        name: 'search',
        params: {videos: videoInfo.value}
      });

    } catch (error) {
      alert("请先登录！")
      console.error('Error fetching videos:', error);
    }
}
function closeModal() {
  showModal.value = false;
}

function setInitialLoginStatus() {
  // 从localStorage中获取用户名
  const storedUsername = localStorage.getItem("username");
  // 如果用户名存在，则用户已登录
  if (storedUsername) {
    username.value = storedUsername;
    isLogged.value = true;
  }
}

// 组件初始化时调用
onMounted(() => {
  setInitialLoginStatus();
  // ...其他onMounted逻辑
});


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
    .recommend{
      margin-left: 10%;
    }
    .front_page {
      position: relative; // 保留定位上下文
      cursor: pointer;
      display: inline-block; // 确保元素是行内块元素以便可以居中

    }

    .dropdown {
      display: none; // 默认不显示
      position: absolute;
      min-width: 160px; // 调整为期望的宽度
      box-shadow: 0px 8px 16px 0px rgba(0,0,0,0.2);
      z-index: 1;
      right: 0;
      left: -130%; // 将下拉框的左边界设置为触发元素的中心点
      background-color: rgba(24, 24, 24, 0.9); // 调整透明度
      border-radius: 4px; // 添加圆角
      padding: 8px 0; // 添加上下内边距给下拉菜单
    }

    .dropdown-item {
      cursor: pointer;
      width: 100%;
      color: white;
      padding: 12px 16px; // 调整内边距以符合设计
      text-decoration: none;
      border-radius: 5%;
      display: block;
      font-size: 16px;
      background-color: transparent; // 初始背景色透明

      &:hover {
        background-color: #575757; // hover时的背景色
      }
    }

    // 鼠标悬停时显示下拉菜单
    .front_page:hover .dropdown {
      display: block;
    }

    /* 新增一个简单的淡入动画 */
    @keyframes fadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }
    .search-area {
      display: flex; // 使用Flex布局来排列搜索框和按钮
      align-items: center; // 垂直居中对齐搜索框和按钮
      border: 1px solid white; // 给整个搜索区域加上边框
      border-radius: 5px; // 边框圆角效果
      overflow: hidden; // 防止子元素溢出边框
      background-color: #181818; // 搜索区域背景色

    }

    .search-input {
      width: 10%;
      padding: 8px;
      border: none;
      color: white;
      border-radius: 5px;
      outline: none;
      font-size: 16px;
      flex-grow: 1; // 输入框会填充剩余空间
      margin: 0; // 移除外边距
      background-color: transparent; // 输入框背景透明


    }

    .search-btn {
      border: none; // 去除按钮自身的边框
      background-color: #3498db; // 按钮背景颜色
      color: white; // 按钮文字颜色
      padding: 8px 16px; // 按钮内边距
      margin: 0; // 移除外边距
      border-radius: 0; // 由于是在.search-area内，所以移除单独的圆角
      cursor: pointer; // 鼠标悬停时显示指针
      outline: none; // 去除聚焦时的边框
    }
    .search-btn:hover {
      background-color: #2980b9; // 按钮背景颜色变深
    }

    // 输入框聚焦时的样式

    .user {
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
        margin:10px;
        font-size: 15px;
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
  }
}
</style>
