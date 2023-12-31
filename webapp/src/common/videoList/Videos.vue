<template>
  <div class="video-container" @wheel="handleMouseWheel">
    <div class="videoWrapper">
      <video class="myVideo" ref="videoPlayer" :src="videoSource" type="video/mp4" @ended="switchToNextVideo" @timeupdate="updateProgress" @click="toggleVideo"></video>
      <div class="pause-overlay" v-show="!isPlaying"></div>
      <div class="buttons">
        <div class="buttons">
          <div class="prev-button" @click="switchToPrevVideo">⬆</div>
          <div class="next-button" @click="switchToNextVideo">⬇</div>
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
import {computed, onMounted, reactive, ref, toRefs} from "vue";
import {recommendVideo} from "@/api/recommend";

export default {
    setup() {
      const videoPlayer = ref(null);
      const state = reactive({
        videoList: [], // 推荐视频的路径列表
        currentVideoIndex: 0,
        isPlaying: false,
        progressWidth: '0%',
        currentTime: '0:00',
        duration: '0:00',
        speed: '正常',
        showSpeedMenu: false,
        draggingProgress: false
      });

      const fetchRecommendedVideos = async () => {
        try {
          const response = await recommendVideo(); // 假设recommendVideo是一个异步函数
          state.videoList = response.info.map(item => item.video) || [];
          console.log(state.videoList)
        } catch (error) {
          console.error("获取推荐视频失败", error);
        }
      };
      onMounted(() => {
        fetchRecommendedVideos()
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
      const videoSource = computed(() => {
        console.log("return state.videoList[state.currentVideoIndex];\n:",state.videoList)
        return state.videoList[state.currentVideoIndex];
      });

      //setup函数返回的对象中的属性和方法将可用于模板中
      return {
        videoPlayer,
        videoSource,
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

    },
    switchToPrevVideo(){
      this.currentVideoIndex = (this.currentVideoIndex -1 + this.videoList.length) % this.videoList.length;
      this.isPlaying = false;

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
  },
};


</script>

<style scoped>
.video-container {
  height: 100vh;
  width: 100vw;
  position: relative;
  overflow: hidden;
}

.videoWrapper {
  width: 100%;
  height: 100%;
  position: relative;
}

.myVideo {
  position: absolute;
  top: 70px;
  width: 100%;
  height: 95%;
}

.menu {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
  display: flex;
  justify-content: space-between;
  padding: 8px;
  font-size: 14px;
  z-index: 2;
}

.play,
.time,
.progress-bar,
.quick {
  cursor: pointer;
}

.progress-bar {
  width: 50%;
  height: 8px;
  background-color: #333;
  position: relative;
  top: 50%;
  transform: translateY(-50%);
  border-radius: 4px;
  display: flex;
  align-items: center;
}

.progress {
  height: 100%;
  background-color: white;
  border-radius: 4px;
}
.dot {
  position: absolute;
  height: 10px;
  width: 10px;
  background-color: #ff5722;
  border-radius: 50%;
  top: 50%;
  transform: translateY(-50%);
  margin-left: -3px; /* 将小圆点向左偏移5像素 */
}
.quick-list {
  position: absolute;
  bottom: 100%;
  left: 0;
  display: none;
}

.quick-list ul {
  list-style: none;
  margin: 0;
  padding: 0;
}

.quick-list li {
  cursor: pointer;
  padding: 8px;
  transition: background-color 0.2s;
}

.quick-list li:hover {
  background-color: rgba(255, 255, 255, 0.1);
}
.next-button, .prev-button {
  position: absolute;
  width: 50px;
  height: 50px;
  background-color: gray;
  color: white;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  text-align: center;
  line-height: 50px;
  font-size: 20px;
  right: 40px;
}

.next-button {
  top: 50%;
}

.prev-button {
  top: 40%;
}

.prev-button:hover, .next-button:hover {
  background-color: #45a049;
}

.pause-overlay {
  position: absolute;
  top: 40%;
  left: 50%;
  width: 10%;
  height: 10%;
  background: url('../../assets/icons/播放.svg') center center no-repeat;
  background-size: contain;
  z-index: 1;
}

</style>
