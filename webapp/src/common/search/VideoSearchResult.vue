<template>
  <div class="video-gallery">
    <div v-if="ResultVideo.length === 0" class="no-data" style="font-size: 30px">
      抱歉，没有此类视频数据哦～
    </div>
    <div
      v-else
      v-for="video in videos"
      :key="video.vid"
      class="video-card-link"
    >
      <router-link :to="{ name: 'videodetail', params: { vid: video.vid } }">
      <div class="video-card" @click="checkDetail(video.vid)">
        <img :src="video.thumbnail" alt="视频缩略图" class="thumbnail" />
        <div class="video-info">
          <div class="video-title" v-html="video.highlight_subtitled"></div>
          <div class="video-details">
            <span class="likes">{{ video.likes}} ❤️</span>
            <span class="likes">标题：：{{ video.subtitled }}</span>
            <span class="upload-time">发布于：{{ formatUploadTime(video.upload_time) }}前</span>
          </div>
        </div>
      </div>
      </router-link>
    </div>
  </div>
</template>
<script>
import {useRouter} from "vue-router";
import {ref, watchEffect} from "vue";

export default {
  name: "VideoSearchResult",
  props: {
    videos: {
      type: Array,
      required: true
    },
  },

  setup(props) {
    const router = useRouter();
    const ResultVideo = ref([])
    watchEffect(() => {
      console.log("videos prop has changed:", props.videos);
      ResultVideo.value =props.videos
    });
    const checkDetail = (vid) => {
      console.log("checkDetail方法被调用了");
      router.push({
        name: 'videodetail',
        params: { vid:vid }
      });
    };
    return {
      ResultVideo,
      checkDetail
    }
  },
  methods: {
    formatUploadTime(timestamp) {
      const date = new Date(timestamp);
      const now = new Date();
      const diff = now - date;
      // 这只是一个简单的示例，你需要根据需要实现逻辑
      const hours = Math.floor(diff / 3600000);
      return `${hours} 小时`;
    },

  },
};
</script>

<style scoped>
.video-gallery {
  display: grid;
  /* 设置为最多显示两列，当只有一个视频时会自动居中 */
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  /* 确保grid子项在整个可用空间中居中对齐 */
  justify-items: center;
  justify-content: center;
  gap: 20px;
  padding: 20px;
  position: absolute;
  width: 60%; /* 可以尝试设置为100%，确保占据全部父元素宽度 */
  left: 50%; /* 将左边距设置为50% */
  transform: translateX(-50%); /* 通过transform偏移来确保完全居中 */
  top: 120px;
  background-color: #f7f7f7;
}

.video-card {
  position: relative;
  overflow: hidden;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0,0,0,0.1); /* 添加阴影效果 */
}

.thumbnail {
  width: 100%;
  height: auto;
  display: block; /* 避免底部空白 */
}

.video-info {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  padding: 10px;
  background-color: rgba(0, 0, 0, 0.5); /* 提高透明度 */
  color: #ffffff;
  text-shadow: 0 1px 3px rgba(0,0,0,0.1); /* 文本阴影 */
}

.video-title {
  font-size: 16px; /* 增大字体 */
  margin-bottom: 5px;
}

.video-details {
  display: flex;
  justify-content: space-between;
  font-size: 14px; /* 调整字体大小 */
}
</style>


