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
      <div class="video-info">
        <div class="meta-info">
          <h3>@{{ currentNickName }}</h3>
          <p>title:{{currentSubTitle}}</p>
          <p>上传时间：{{ formatDate(currentUploadTime) }}</p>
        </div>
      </div>

      <div class="pause-overlay" v-show="!isPlaying"></div>
        <div class="buttons">
          <div class="prev-button" @click="switchToPrevVideo">⬆</div>
          <div class="next-button" @click="switchToNextVideo">⬇</div>
          <div class="like" >
            <img v-if="currentLiked" @click="unLikeThis"  src="../../assets/icons/like%20(1).svg"  alt="Liked" />
            <img v-else @click="likeThis" src="../../assets/icons/like.svg"  alt="Like" />
            <div class="likesNumber">{{likesNumber}}</div>
          </div>
        </div>
      <div class="menu"
          >
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
    const isLiked = ref(false); // 初始化点赞状态
    const likesNumber = ref(0); // 初始化点赞数为 0

    const state = reactive({
      videoList: [], // 推荐视频的路径列表
      isPlaying: false,
      progressWidth: '0%',
      currentTime: '0:00',
      duration: '0:00',
      speed: '正常',
      showSpeedMenu: false,
      draggingProgress: false,
      // showControlsBar: false, // 用于控制控制条显示的状态

      // videoVid: 41, //当前视频的vid

    });
    // const currentVideoVid = computed(() => {
    //   // 确保有视频并且currentVideoIndex是有效的
    //   if (props.videosInfo.length > 0 && props.videosInfo[currentVideoIndex.value]) {
    //     return props.videosInfo[currentVideoIndex.value].vid;
    //   }
    //   return null; // 没有有效的当前视频
    // });
    const updateLikeStatus = () => {
      if (isLiked.value) {
        console.log("this.currentVideo.is_like:",this.currentVideo.is_like)
        isLiked.value = this.currentVideo.is_like;
      }
    };
    onMounted(() => {
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
        watch(() => props.videosInfo[currentVideoIndex.value].likes, (newLikes) => {
          likesNumber.value = newLikes; // 更新likesNumber的值
        });
      document.addEventListener('keydown', handleSpacebar);
      // 记得在卸载组件时移除事件监听器
      return () => {
        document.removeEventListener('keydown', handleSpacebar);
      };
    });
    // 计算属性，用于获取当前视频的点赞状态
    const currentLiked = computed(() => {
      const currentVideo = props.videosInfo[currentVideoIndex.value];
      return currentVideo ? currentVideo.is_like : false;
    });

    state.videoList = props.videosInfo.map(item=> {
      item.video,item.vid
    })

    //setup函数返回的对象中的属性和方法将可用于模板中
    return {
      videoPlayer,
      currentVideoIndex,
      // currentVideo,
      updateLikeStatus,
      isPlaying,
      currentLiked,
      likesNumber,
      ...toRefs(state),
    };
  },

  computed: {
    // 当前播放视频的信息
    currentVideo() {
      return this.videosInfo[this.currentVideoIndex];
    },
    // 当前视频的src属性
    currentSrc() {
      return this.currentVideo.video;
    },
    // 当前视频的vid，用于点赞功能
    currentVid() {
      return this.currentVideo.vid;
    },
    currentLiked(){
      return this.currentVideo.is_like
    },
    currentNickName(){
      return this.currentVideo.uploader_username
    },
    currentUploadTime(){
      return this.currentVideo.upload_time
    },
    currentLikes(){
      return this.currentVideo.likes
    },
    currentSubTitle(){
      return this.currentVideo.subtitled
    }
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
    //格式化日期
    formatDate(timestamp) {
      const date = new Date(timestamp);
      return date.toLocaleDateString("zh-CN");
    },
    switchToPrevVideo() {
      this.currentVideoIndex = (this.currentVideoIndex - 1 + this.videoList.length) % this.videoList.length;
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

    // showControls() {
    //   this.state.showControlsBar = true; // 显示控制条
    //   console.log(this.state.showControlsBar)
    //   // 显示控制条
    // },
    // hideControls() {
    //   this.state.showControlsBar = false; // 显示控制条
    //   console.log(this.state.showControlsBar)
    //   // 隐藏控制条
    // },
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
      this.likesNumber += 1;
      if (!this.currentLiked) {
        try {
          const response = await likeVideo(this.currentVid, parseInt(localStorage.getItem('uid')));
          console.log("点赞成功", response);
          // 更新当前视频对象的点赞状态
          this.videosInfo[this.currentVideoIndex].is_like = true;
          // 由于currentLiked是计算属性，这将自动更新它
        } catch (error) {
          console.error("点赞失败", error);
        }
      }
    },
    async unLikeThis() {
      this.likesNumber -= 1;
      if (this.currentLiked) {
        try {
          const response = await unLikeVideo(this.currentVid, parseInt(localStorage.getItem('uid')));
          console.log("取消点赞成功", response);
          // 更新当前视频对象的点赞状态
          this.videosInfo[this.currentVideoIndex].is_like = false;
          // 由于currentLiked是计算属性，这将自动更新它
        } catch (error) {
          console.error("取消点赞失败", error);
        }
      }
    },

  }
};


</script>
<style>
@import "./video.css";


.video-info {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center; /* 居中对齐子项 */
  margin-bottom: 20px; /* 和视频容器之间的距离 */
  position: absolute;
  bottom: 20px;
  border-radius: 10%;
}
.meta-info {
  background: rgba(0, 0, 0, 0.3); /* 透明黑背景 */
  color: #fff;
  width: calc(100% - 20px); /* 减去边距的宽度 */
  padding: 10px;
  box-sizing: border-box; /* 确保内边距不会影响元素的总宽度 */
  border-radius: 4px; /* 给元数据信息添加圆角 */
}

.meta-info h3 {
  font-size: 1.5em; /* 较大的标题字号 */
  margin: 0 0 10px 0; /* 只在标题下方添加间距 */
}

.meta-info p {
  font-size: 1em; /* 段落的字号与默认相同 */
  margin: 5px 0; /* 在段落之间添加垂直间距 */
}



/* 视频容器的样式 */
.video-container {
  /* 您可以添加适用于视频容器的样式，例如宽度、高度等 */
}
</style>
