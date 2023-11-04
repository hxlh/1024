# webApp

## 开发流程

### 安装依赖

安装项目所需的依赖：

```bash
npm install
```

### 通过vite运行开发服务器

启动开发服务器进行开发：

```bash
npm run dev
```

### 构建生产

构建用于生产的应用：

```bash
npm run build
```

## 目录结构

```plaintext
src/
├── api/                   # 后端API交互函数
│   ├── upload.js          # 文件上传API
│   └── user.js            # 用户相关API
├── assets/                # 静态资源
│   ├── icons/             # 图标资源
│   ├── images/            # 图片资源
│   │   ├── arrow-down.png # 下箭头图标
│   │   └── arrow-up.png   # 上箭头图标
│   ├── base.css           # 基础样式
│   └── main.css           # 主要样式
├── common/                # 通用组件
│   ├── search/            # 搜索组件
│   │   └── VideoSearchResult.vue //搜索视频组件
│   ├── upload/            # 上传组件 
│   │   └── Upload.vue
│   ├── user/              # 用户组件
│   │   ├── LoginModule.vue # 用户登录模态框
│   │   ├── RegisterModule.vue # 用户注册模态框
│   │   └── UserInfo.vue	# 用户个人信息跳转
│   └── videoList/         # 视频列表组件
│       ├── Menu.vue # 视频播放进度条 - 未优化
│       ├── Videos.vue	# 视频组件
│       ├── VideoShow.vue
│       └── video.css	# 个别视频样式（用于测试） -可删除
├── router/               # 路由配置
│   ├── home.js           # 首页路由
│   └── index.js          # 路由索引
├── views/                # 视图组件
│   ├── AboutView.vue     # 关于页面
│   ├── HomeView.vue      # 主页视图
│   ├── Layout.vue        # 布局组件
│   └── VideoView.vue     # 视频视图
├── App.vue               # 主应用组件
├── App.copy.vue          # 应用副本
├── main.js               # 入口文件
├── .gitignore            # Git忽略文件
└── index.html            # HTML入口文件
```

## 组件说明

- `VideoSearchResult.vue` - 展示视频搜索结果。
- `Upload.vue` - 处理文件上传逻辑。
- `LoginModule.vue` - 用户登录模态框组件。
- `RegisterModule.vue` - 用户注册模态框组件。
- `UserInfo.vue` - 显示用户信息。
- `Menu.vue` - 视频进度条组件。
- `Videos.vue` - 视频列表显示。

## 路由配置

- `index.js` - 项目的路由入口文件。
