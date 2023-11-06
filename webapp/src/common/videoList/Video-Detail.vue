<template>
  <div class="video-container" @wheel="handleMouseWheel">
    <div class="backButton">
      <div class="back-button" @click="goBack">返回</div>
    </div>
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
      <div class="video-info">
        <div class="meta-info">
          <h3>@{{ currentNickName }}</h3>
          <p>title:{{currentSubTitle}}</p>
          <p>点赞数：{{ currentLikes }}❤️</p>
          <p>上传时间：{{ formatDate(currentUploadTime) }}</p>
        </div>
      </div>
      <div class="pause-overlay" v-show="!isPlaying"></div>
      <div class="buttons">
        <div class="like" >
          <img v-if="currentLiked" @click="unLikeThis" src="../../assets/icons/like%20(1).svg"  alt="Liked" />
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
import {useRouter} from "vue-router";
export default {

  props: {
    videoInfo: {  // 单个视频的信息
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const videoPlayer = ref(null);
    const currentVideoIndex = ref(0);
    const isPlaying = ref(false);

    const state = reactive({
      isPlaying: false,
      progressWidth: '0%',
      currentTime: '0:00',
      duration: '0:00',
      speed: '正常',
      showSpeedMenu: false,
      draggingProgress: false,
      // videoVid: 41, //当前视频的vid

    });
    const router = useRouter();
    function goBack() {
      router.back();
    }
    const currentSrc = computed(() => props.videoInfo.video);
    const currentVid = computed(() => props.videoInfo.vid);
    const currentNickName = computed(() => props.videoInfo.uploader_username);
    const currentLikes = computed(() => props.videoInfo.likes);
    const currentUploadTime = computed(() => props.videoInfo.upload_time);
    const currentSubTitle = computed(()=>{props.videoInfo.subtitled})
    onMounted(() => {
      console.log(currentVideoIndex)
      // fetchRecommendedVideos()
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

      document.addEventListener('keydown', handleSpacebar);
      // 记得在卸载组件时移除事件监听器
      return () => {
        document.removeEventListener('keydown', handleSpacebar);
      };
    });




    //setup函数返回的对象中的属性和方法将可用于模板中
    return {
      videoPlayer,
      currentVideoIndex,
      currentVid,
      currentSrc,
      isPlaying,
      goBack,
      currentNickName,
      currentLikes,
      currentUploadTime,
      currentSubTitle,
      ...toRefs(state),
    };

  },
  computed: {
    // 当前视频的src属性
    currentSrc() {
      console.log(this.videoInfo.video)
      return this.videoInfo.video;
    },
    // 当前视频的vid，用于点赞功能
    currentVid() {
      return this.videoInfo.vid;
    },
    currentLiked(){
      return this.videoInfo.is_like;
    },

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
      if (!this.videoList || this.videoList.length === 0) {
        this.resetVideoPlayer();
        this.isPlaying = true;
      } else {
        this.currentVideoIndex = (this.currentVideoIndex + 1) % this.videoList.length;
        this.isPlaying = false;
      }
    },

    switchToPrevVideo(){
      this.currentVideoIndex = (this.currentVideoIndex -1 + this.videoList.length) % this.videoList.length;
      this.isPlaying = false;
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
    //格式化日期
    formatDate(timestamp) {
      const date = new Date(timestamp);
      return date.toLocaleDateString("zh-CN");
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

    async likeThis() {
      console.log("likethis被调用了");
      if (!this.currentLiked) {
        try {
          const response = await likeVideo(this.currentVid, parseInt(localStorage.getItem('uid')));
          console.log("点赞成功", response);
          // 更新当前视频对象的点赞状态
          this.videoInfo.is_like = true;
          // 由于currentLiked是计算属性，这将自动更新它
        } catch (error) {
          console.error("点赞失败", error);
        }
      }
    },
    async unLikeThis() {
      if (this.currentLiked) {
        try {
          const response = await unLikeVideo(this.currentVid, parseInt(localStorage.getItem('uid')));
          console.log("取消点赞成功", response);
          // 更新当前视频对象的点赞状态
          this.videoInfo.is_like = false;
          // 由于currentLiked是计算属性，这将自动更新它
        } catch (error) {
          console.error("取消点赞失败", error);
        }
      }
    },
  },
};


</script>
<style scoped lang="scss">
@import "./video.css";
.buttons {
  width: 70px;
  height: 70px;
  /* 如果.buttons也需要居中内容，可以同样应用flex布局 */
  display: flex;
  justify-content: center;
  align-items: center;
  .like {
    width: 60px;
    height: 60px;
    display: flex;
    justify-content: center;
    align-items: center;
    img {
      max-width: 100%; /* 确保图片不会超出容器大小 */
      max-height: 100%; /* 确保图片高度不会超出容器大小 */
    }
  }



}

.back-button {
  position: absolute;
  top: 80px;
  left: 30px;
  cursor: pointer;
  /* Style your back button as needed */
  padding: 10px;
  background-color: #f0f0f0;
  border: 1px solid #dcdcdc;
  border-radius: 5px;
  margin: 10px;
  display: inline-block;
  z-index: 99;
}


</style>

