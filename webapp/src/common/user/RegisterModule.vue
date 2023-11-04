<template>
  <div v-if="showRegisterModal" class="modal" @click.stop="RegisterClick">
    <div class="modal-header">
      <span class="welcome-text">欢迎注册</span>
      <button class="close-btn" @click="closeModal">×</button>
    </div>

    <div class="modal-body">
      <input type="text" placeholder="请输入别名" v-model="registerState.c">
      <input type="text" placeholder="请输入账号" v-model="registerState.username">
      <input type="password" placeholder="请输入密码" v-model="registerState.password">
      <div class="red-line"></div>

      <button class="login-btn" @click="submitRegister">注册</button>
    </div>
    <div class="actions-row" >
      <button class="register-btn" >已有账号？<a href="#" @click="switchToLogin">去登录</a></button>
    </div>
  </div>
</template>

<script setup>
import {ref, defineProps, defineEmits, onMounted, reactive} from "vue";
import { register } from "@/api/user";


const registerState = reactive({
  nickname: "",
  username: "",
  password: ""
});
const props = defineProps({
  showRegisterModal: {
    type: Boolean,
    default: false
  }});

const emit = defineEmits("closeRegister")
function RegisterClick(event){
  if (event.target.classList.contains("modal")) {
    emit("closeRegister");
  }
}
function switchToLogin() {
  emit("switchToLogin");
}

function submitRegister(){
  console.log("submitRegister方法被调用了");
  register(registerState.username, registerState.password,registerState.nickname)
    .then(response => {
      console.log(response);
      alert("注册成功");
      emit("closeRegister");
    })
    .catch(error => {
      console.log(error);
      alert("注册失败");
    });
}
function closeModal() {
  console.log("Close button clicked");
  emit("closeRegister");

}
// 监听 ESC 键按下事件，关闭注册框
onMounted(() => {
  window.addEventListener("keydown", event => {
    if (event.key === "Escape") {
      closeModal();
    }
  });
});
</script>

<style scoped lang="scss">
.modal {
  width: 400px;
  background-color: #fff;
  padding: 20px;
  box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.5);
  border-radius: 8px;
  position: fixed;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  z-index: 1000;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.welcome-text {
  color: #333;
  font-size: 20px;
  position: relative;
  top: 15px;
  left: 40%;
}

.modal-body input {
  width: 100%;
  padding: 10px;
  margin: 10px 0;
  border: 1px solid #ccc;
  border-radius: 4px;
  cursor: text;
}

.red-line {
  height: 2px;
  background-color: #f00;
  margin: 10px 0;
}

.login-btn {
  width: 100%;
  padding: 10px;
  background-color: #f66;
  color: #fff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin: 10px 0;
}

.actions-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.register-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  text-decoration: none;
  color: #f66;
}

.register-link {
  margin: 0;
  text-align: right;
  flex-grow: 1;
  color: gray;
  font-size: 12px;
}

.register-link a {
  font-size: 15px;
  color: #f66;
}

.modal, .modal * {
  cursor: default;
}
</style>
