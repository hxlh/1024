import {createRouter, createWebHistory} from 'vue-router'
import Upload from "../common/upload/Upload.vue";
import Home from "../views/Home.vue";
import Videos from "../common/videoList/Videos.vue";
import LoginModule from "@/common/user/LoginModule.vue";
import VideoSearchResult from "@/common/search/VideoSearchResult.vue";
import UserInfo from "@/common/user/UserInfo.vue";
import Recommend from "@/common/recommend/RecommendVideo.vue";
import VideoSearchDetail from "@/common/search/VideoSearchDetail.vue";
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [

    {
      path: '/videos',
      name: 'videos',
      component: Videos
    },
    {
      path: '/test',
      name: 'test',
      component: Videos
    },
    {
      path: '/upload',
      name: 'upload',
      component: Upload
    },
    {
      path:'/login',
      name:'login',
      component:LoginModule,
    },
    {
      path:'/search',
      name:'search',
      component:VideoSearchResult,
    },
    {
      path:'/userinfo',
      name:'userinfo',
      component:UserInfo,
    },
    {
    path:'/',
    name:'recommend',
      component: Recommend
    },
    {
      path: '/videodetail/:vid', // 路径参数
      name:'videodetail',
      component:VideoSearchDetail,
      props: true // 允许将路由参数作为props传递给组件

    }
  ]
})

export default router
