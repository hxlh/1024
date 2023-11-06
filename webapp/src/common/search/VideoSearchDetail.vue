<template>
  <div class="videoDetail" >
    <VideoDetail v-if="videoInfo" :videoInfo="videoInfo"></VideoDetail>
    <div v-else>没有视频信息可供显示。</div>
  </div>
</template>

<script>
import VideoDetail from "@/common/videoList/Video-Detail.vue";
import { ref, onMounted, watch } from "vue";
import { useRoute } from "vue-router";

export default {
  name: "VideoSearchDetail",
  components: {VideoDetail},

  setup() {
    const videoInfo = ref(null); // 初始化为 null
    const route = useRoute();
    // 根据视频ID获取视频信息
    const fetchVideoInfo = (vid) => {
      const allVideosInfo = JSON.parse(localStorage.getItem("videoSearchInfo"));
      if (allVideosInfo) {
        const matchingVideo = allVideosInfo.find(video => video.vid.toString() === vid);
        videoInfo.value = matchingVideo || null;
      }
    };

    onMounted(() => {
      // 首次挂载时获取视频信息
      if (route.params.vid) {
        fetchVideoInfo(route.params.vid);
      }
    });

    watch(() => route.params.vid, (newVid) => {
      // 当路由参数变化时获取新的视频信息
      fetchVideoInfo(newVid);
    });

    // 返回响应式状态
    return {
      videoInfo
    };
  },
};
</script>

<style scoped>

</style>
