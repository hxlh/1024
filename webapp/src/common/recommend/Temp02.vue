<template>
  <div class="video-container" @wheel="handleMouseWheel">
    <div class="videoWrapper">
      <video
        class="myVideo"
        ref="videoPlayer"
        :src="currentSrc"
        type="video/mp4"
        @ended="switchToNextVideo"
        @timeupdate="updateProgress"
        @click="toggleVideo">
      </video>
      <div class="pause-overlay" v-show="!state.isPlaying"></div>
      <div class="buttons">
        <div class="prev-button" @click="switchToPrevVideo">⬆</div>
        <div class="next-button" @click="switchToNextVideo">⬇</div>
        <div class="like">
          <img v-if="state.isLiked" @click="unLikeThis" src="../../assets/icons/like%20(1).svg"  alt="Liked" />
          <img v-else @click="likeThis" src="../../assets/icons/like.svg"  alt="Like" />
        </div>
      </div>
      <div class="menu" @mouseenter="showControls" @mouseleave="hideControls">
        <div class="play" @click="toggleVideo">{{ state.isPlaying ? '暂停' : '播放' }}</div>
        <div class="time">{{ state.currentTime }} / {{ state.duration }}</div>
        <div class="progress-bar" ref="progressBar" @mousedown="startDragging">
          <div class="progress" :style="{ width: state.progressWidth }"></div>
          <div class="dot" :style="{ left: state.progressWidth }"></div>
        </div>

        <div class="quick" @click="toggleSpeedMenu">{{ state.speed }}</div>
        <div class="quick-list" v-show="state.showSpeedMenu">
          <ul>
            <li @click="changeSpeed(1)">正常</li>
            <li @click="changeSpeed(1.25)">1.25倍速</li>
            <li @click="changeSpeed(1.5)">1.5倍速</li>
            <li @click="changeSpeed(2)">2倍速</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed, onMounted, onUnmounted, reactive, ref, watch } from "vue";
import { likeVideo, unLikeVideo } from "@/api/like";

export default {
  props: {
    videosInfo: {
      type: Array,
      required: true,
    },
  },
  setup(props) {
    const videoPlayer = ref(null);
    const progressBar = ref(null);
    const currentVideoIndex = ref(0);

    const state = reactive({
      videoList: props.videosInfo, // 推荐视频的路径列表
      isPlaying: false,
      isLiked: false,
      progressWidth: '0%',
      currentTime: '0:00',
      duration: '0:00',
      speed: '正常',
      showSpeedMenu: false,
      draggingProgress: false,
    });

    const currentVideo = computed(() => props.videosInfo[currentVideoIndex.value] || {});
    const currentSrc = computed(() => currentVideo.value.video || '');
    const currentVid = computed(() => currentVideo.value.vid || '');

    const handleSpacebar = (event) => {
      if (event.code === 'Space') {
        toggleVideo();
        event.preventDefault();
      }
    };

    onMounted(() => {
      document.addEventListener('keydown', handleSpacebar);
      // 如果有必要的话，这里可以进行一些初始化操作
    });

    onUnmounted(() => {
      document.removeEventListener('keydown', handleSpacebar);
    });

    watch(currentVideo, (newVideo) => {
      state.isLiked = newVideo.is_like;
    }, { immediate: true });

    // 所有的方法都应该在 setup() 函数内定义
    const toggleVideo = () => {
      if (videoPlayer.value.paused) {
        videoPlayer.value.play();
      } else {
        videoPlayer.value.pause();
      }
      state.isPlaying = !state.isPlaying;
    };
    const updateLikeStatus = () => {
      state.isLiked = currentVideo.value.is_like;
    };
    const switchToNextVideo=() => {
      currentVideoIndex.value = (currentVideoIndex.value + 1) % state.videoList.length;
      this.isPlaying = false;
      updateLikeStatus();
    };
    const switchToPrevVideo=() =>{
      currentVideoIndex.value = (currentVideoIndex.value - 1 + state.videoList.length) % state.videoList.length;
      this.isPlaying = false;
      this.updateLikeStatus()

    };
    const resetVideoPlayer=() => {
      const video = this.$refs.videoPlayer;
      video.currentTime = 0; // 重置视频时间为0
      video.load(); // 重新加载视频，这也会更新视频的持续时间等
    };
    const handleMouseWheel = (event) => {
      const delta = Math.sign(event.deltaY);
      if (delta > 0) {
        switchToNextVideo();
      } else if (delta < 0) {
        switchToPrevVideo();
      }
    };

    const showControls=() =>  {
      // 显示控制条
    };
    const hideControls=() =>  {
      // 隐藏控制条
    };
    const updateProgress=() =>  {
      if (!this.draggingProgress) {
        const video = this.$refs.videoPlayer;
        const currentTime = video.currentTime;
        const duration = video.duration;
        this.currentTime = this.formatTime(currentTime);
        this.duration = this.formatTime(duration);
        this.progressWidth = (currentTime / duration * 100) + '%';
      }
    };
    const startDragging=(event) =>  {
      this.draggingProgress = true;
      document.addEventListener('mouseup', this.dragging);
      document.addEventListener('mousemove', this.stopDragging);
      this.dragging(event);
    };
    const dragging=(event) =>  {
      if (this.draggingProgress) {
        const progressBar = this.$refs.progressBar;
        const video = this.$refs.videoPlayer;
        const rect = progressBar.getBoundingClientRect();
        const offsetX = event.clientX - rect.left;
        const newProgress = Math.max(0, Math.min(1, offsetX / rect.width));
        video.currentTime = newProgress * video.duration;
        this.progressWidth = (newProgress * 100) + '%';
      }
    };
    const stopDragging=() => {
      this.draggingProgress = false;
      document.removeEventListener('mousemove', this.dragging);
      document.removeEventListener('mouseup', this.stopDragging);
    };
    const formatTime=(time) => {
      const minutes = Math.floor(time / 60);
      const seconds = Math.floor(time % 60);
      return `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`;
    };
    const toggleSpeedMenu=() => {
      this.showSpeedMenu = !this.showSpeedMenu;
    };
    const changeSpeed=(speed) => {
      const video = this.$refs.videoPlayer;
      video.playbackRate = speed;
      this.speed = speed === 1 ? '正常' : `${speed}倍速`;
      this.showSpeedMenu = false;
    };
    // ... 其他方法同理

    return {
      videoPlayer,
      progressBar,
      currentVideoIndex,
      currentSrc,
      currentVid,
      toggleVideo,
      switchToNextVideo,
      resetVideoPlayer,
      handleMouseWheel,
      showControls,
      hideControls,
      updateProgress,
      startDragging,
      dragging,
      stopDragging,
      formatTime,
      toggleSpeedMenu,
      changeSpeed,
      updateLikeStatus,


      // ... 其他方法和数据
      state,
    };
  },
};
</script>

<style scoped>
@import "./video.css";

</style>
