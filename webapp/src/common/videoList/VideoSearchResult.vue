<template>
  <div class="video-gallery">
    <div v-for="video in videos" :key="video.id" class="video-card">
      <img :src="video.thumbnail" alt="视频缩略图" class="thumbnail"/>
      <div class="video-info">
        <div class="video-title">{{ video.title }}</div>
        <div class="video-details">
          <span class="play-time">{{ video.playTime }}</span>
          <span class="likes">{{ video.likes }} ❤️</span>
          <span class="upload-time">{{ video.uploadTime }}前</span>
        </div>
      </div>
    </div>
  </div>
</template>



<script>
export default {
  name: "VideoSearchResult",
  props: {
    videos: {
      type: Array,
      required: true,
      default:() => [
        {
          playTime: '1:00',
          likes: 100,
          uploadTime: '1天',
        }
      ]
    }
  }
}
</script>

<style scoped>
.video-gallery {
  position: relative;
  top: 80px;
  background-color: green;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.video-card {
  position: relative;
  overflow: hidden;
  border-radius: 10px;
  background-color: pink;
}

.thumbnail {
  width: 100%;
  height: auto;
  background-color: red;
}

.video-info {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  padding: 10px;
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
}

.video-title {
  font-weight: bold;
}

.video-details {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
}
</style>

