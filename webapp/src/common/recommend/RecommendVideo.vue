<template>
  <div>
    <VideoPlayer v-if="state.videosInfo.length > 0" :videosInfo="state.videosInfo"></VideoPlayer>
  </div>
</template>

<script>
import VideoPlayer from "@/common/videoList/Video-Player.vue";
import { defineComponent, onMounted, reactive } from "vue";
import { recommendVideo } from "@/api/recommend";

export default defineComponent({
  components: { VideoPlayer },
  setup() {
    const state = reactive({
      videosInfo: [], // 推荐视频的信息列表
      vid: 0 // 当前播放视频的ID
    });

    const fetchRecommentVideos = async () => {
      try {
        const response = await recommendVideo(localStorage.getItem('uid'));
        state.videosInfo = response.data.info
        console.log("state.videoList:",state.videosInfo)
      } catch (error) {
        console.error("获取推荐视频失败", error);
      }
    };

    onMounted(() => {
      fetchRecommentVideos();
    });

    return { state };
  }
});
</script>

<style>
/* 您的样式代码 */
</style>
