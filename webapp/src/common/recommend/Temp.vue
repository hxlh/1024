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
      <div class="pause-overlay" v-show="!isPlaying"></div>
      <div class="buttons">
        <div class="prev-button" @click="switchToPrevVideo">⬆</div>
        <div class="next-button" @click="switchToNextVideo">⬇</div>
        <div class="like" >
          <img v-if="isLiked" @click="unLikeThis" src="../../assets/icons/like%20(1).svg"  alt="Liked" />
          <img v-else @click="likeThis" src="../../assets/icons/like.svg"  alt="Like" />
        </div>
      </div>
      <div class="menu" @mouseenter="showControls" @mouseleave="hideControls">
        <div class="play" @click="toggleVideo">{{ isPlaying ? '暂停' : '播放' }}</div>
        <div class="time">{{ currentTime }} / {{ duration }}</div>
        <div class="progress-bar" ref="progressBar" @mousedown="startDragging">
          <div class="progress" :style="{ width: progressWidth }"></div>
          <div class="dot" :style="{ left: progressWidth }"></div>
        </div>

        <div class="quick" @click="toggleSpeedMenu">{{ speed }}</div>
        <div class="quick-list" v-show="showSpeedMenu">
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
import {computed, onMounted, reactive, ref, toRefs, watch} from "vue";
import {likeVideo, unLikeVideo} from "@/api/like";
export default {

  props: {
    videosInfo: {
      type: Array,
      required: true,
    },

  },
  setup(props) {
    const { videosInfo } = toRefs(props); // 把 props 转换为响应式的引用
    console.log(videosInfo)
    const videoPlayer = ref(null);
    const currentVideoIndex = ref(0);
    const isPlaying = ref(false);
    const isLiked = ref(false);

    const state = reactive({
      videoList: [], // 推荐视频的路径列表
      isPlaying: false,
      progressWidth: '0%',
      currentTime: '0:00',
      duration: '0:00',
      speed: '正常',
      showSpeedMenu: false,
      draggingProgress: false,
      // videoVid: 41, //当前视频的vid

    });
    const currentVideoVid = computed(() => {
      // 确保有视频并且currentVideoIndex是有效的
      if (props.videosInfo.length > 0 && props.videosInfo[currentVideoIndex.value]) {
        return props.videosInfo[currentVideoIndex.value].vid;
      }
      return null; // 没有有效的当前视频
    });
    const updateLikeStatus = () => {
      if (isLiked.value) {
        isLiked.value = currentVideo.is_like;
      }
    };
    onMounted(() => {
      console.log(state.videoList)
      console.log(currentVideoIndex)
      // fetchRecommendedVideos()
      updateLikeStatus();

      const handleSpacebar = (event) => {
        if (event.code === 'Space') {
          if (videoPlayer.value.paused) {
            videoPlayer.value.play();

          } else {
            videoPlayer.value.pause();
          }
          event.preventDefault();
        }
      };
      const currentVideo = computed(() => {
        return props.videosInfo[currentVideoIndex.value] || {};
      });

      const currentSrc = computed(() => {
        return currentVideo.value.video || '';
      });

      const currentVid = computed(() => {
        return currentVideo.value.vid;
      });
      watch(currentVideo, (newVideo) => {
        isLiked.value = newVideo.is_like;
      }, { immediate: true });
      document.addEventListener('keydown', handleSpacebar);
      // 记得在卸载组件时移除事件监听器
      return () => {
        document.removeEventListener('keydown', handleSpacebar);
      };
    });
    state.videoList = props.videosInfo.map(item=> {
      item.video,item.vid
    })

    //setup函数返回的对象中的属性和方法将可用于模板中
    return {
      videoPlayer,
      currentVideoIndex,
      isLiked,
      isPlaying,
      currentVideoVid,
      currentSrc,
      currentVid,
      ...toRefs(state),
    };
  },


  methods: {
    toggleVideo() {
      const video = this.$refs.videoPlayer;
      if (this.isPlaying) {
        video.pause();
      } else {
        video.play();
      }
      this.isPlaying = !this.isPlaying;
    },
    switchToNextVideo() {
      this.currentVideoIndex = (this.currentVideoIndex + 1) % this.videoList.length;
      this.isPlaying = false;
      this.updateLikeStatus();

    },
    switchToPrevVideo(){
      this.currentVideoIndex = (this.currentVideoIndex -1 + this.videoList.length) % this.videoList.length;
      this.isPlaying = false;
      this.updateLikeStatus()

    },

    resetVideoPlayer() {
      const video = this.$refs.videoPlayer;
      video.currentTime = 0; // 重置视频时间为0
      video.load(); // 重新加载视频，这也会更新视频的持续时间等
    },
    handleMouseWheel(event) {
      const delta = Math.sign(event.deltaY);
      if (delta > 0) {
        this.switchToNextVideo();
      } else if (delta < 0) {
        // 切换到上一个视频
        this.currentVideoIndex = (this.currentVideoIndex - 1 + this.videoList.length) % this.videoList.length;
        this.isPlaying = false;
      }
    },

    showControls() {
      // 显示控制条
    },
    hideControls() {
      // 隐藏控制条
    },
    updateProgress() {
      if (!this.draggingProgress) {
        const video = this.$refs.videoPlayer;
        const currentTime = video.currentTime;
        const duration = video.duration;
        this.currentTime = this.formatTime(currentTime);
        this.duration = this.formatTime(duration);
        this.progressWidth = (currentTime / duration * 100) + '%';
      }
    },
    startDragging(event) {
      this.draggingProgress = true;
      document.addEventListener('mouseup', this.dragging);
      document.addEventListener('mousemove', this.stopDragging);
      this.dragging(event);
    },
    dragging(event) {
      if (this.draggingProgress) {
        const progressBar = this.$refs.progressBar;
        const video = this.$refs.videoPlayer;
        const rect = progressBar.getBoundingClientRect();
        const offsetX = event.clientX - rect.left;
        const newProgress = Math.max(0, Math.min(1, offsetX / rect.width));
        video.currentTime = newProgress * video.duration;
        this.progressWidth = (newProgress * 100) + '%';
      }
    },
    stopDragging() {
      this.draggingProgress = false;
      document.removeEventListener('mousemove', this.dragging);
      document.removeEventListener('mouseup', this.stopDragging);
    },
    formatTime(time) {
      const minutes = Math.floor(time / 60);
      const seconds = Math.floor(time % 60);
      return `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`;
    },
    toggleSpeedMenu() {
      this.showSpeedMenu = !this.showSpeedMenu;
    },
    changeSpeed(speed) {
      const video = this.$refs.videoPlayer;
      video.playbackRate = speed;
      this.speed = speed === 1 ? '正常' : `${speed}倍速`;
      this.showSpeedMenu = false;
    },
    async likeThis(){
      isLiked.value = true;
      const uid = parseInt(localStorage.getItem("uid"));
      try {
        await likeVideo(currentVid.value, uid);
      } catch (error) {
        console.error("点赞失败", error);
      }
    },
    async unLikeThis(){
      isLiked.value = false;
      const uid = parseInt(localStorage.getItem("uid"));
      try {
        await unLikeVideo(currentVid.value, uid);
      } catch (error) {
        console.error("取消点赞失败", error);
      }
    },
  },
};


</script>
<style>
@import "./video.css";
</style>

