<template>
  <video controls class="custom-video" @ended="restartVideo">
    <source v-bind:src="url" type="video/mp4">
    Your browser does not support the video tag.
  </video>
  <div class="custom-controls">
    <button class="play-pause-button" @click="clickPlay">&#9658;</button>
    <div class="video-progress">
      <input type="range" min="0" max="100" value="0" step="1">
    </div>
    <button class="mute-button">&#128266;</button>
    <button class="fullscreen-button">&#9723;</button>
  </div>
</template>

<script setup>
import './video.css'

import {ref} from 'vue'
const url = ref('./tmp/video1.mp4')
const progress = ref(0);

function updateProgress() {
  const video = this.$refs.video;
  progress.value = (video.currentTime / video.duration) * 100;
}
function btnUpClicked() {
  url.value = "./tmp/video1.mp4"
  const video = document.querySelector('.myVideo')
  video.load()
}
function btnDownClicked() {
  url.value = "./tmp/video2.mp4"
  const video = document.querySelector('.myVideo')
  video.load()
}

function restartVideo(){
  const video = document.querySelector(".myVideo")
  video.currentTime = 0 //视频播放完毕后，重新播放视频
  video.play()
}
function clickPlay(){
  const video = document.querySelector(".myVideo")
  video.play()
}
function clickPause(){
  const video = document.querySelector(".myVideo")
  video.pause()
}



</script>

<style scoped>

</style>
