# WebApp

该前端项目为`七牛云1024校园马拉松`比赛作品，使用Vite构建的前端Web。

## 项目截图

![image-20231107224606495](/Users/chenhetian/Library/Application Support/typora-user-images/image-20231107224606495.png)

## 快速开始

克隆项目仓库并安装依赖。

```bash
git clone [项目仓库URL]
cd webapp
npm install
```

## 开发

启动开发服务器。

```
bashCopy code
npm run dev
```

访问 [http://localhost:5173](http://localhost:5173查看应用。

## 构建

构建生产环境下的应用。

```
bashCopy code
npm run build
```

## 目录结构


```bash
  src/
	├── api/                     # 后端API交互函数
  │   ├── ...                  # 其他API文件
  ├── assets/                  # 静态资源如图标、图片和样式表
  │   ├── icons/               # 图标资源
  │   ├── images/              # 图片资源
  │   ├── ...                  # 其他资源
  ├── common/                  # 通用组件
  │   ├── search/              # 搜索组件
  │   ├── upload/              # 上传组件
  │   ├── user/                # 用户组件
  │   ├── videoList/           # 视频列表组件
  │   ├── ...                  # 其他通用组件
  ├── router/                  # 路由配置
  │   ├── ...                  # 其他路由配置
  ├── utils/                   # 实用工具函数
  │   ├── debounce.js          # 防抖函数
  │   ├── ...                  # 其他工具函数
  ├── views/                   # 视图组件
  │   ├── ...                  # 其他视图组件
  ├── App.vue                  # 主应用组件
  ├── main.js                  # 入口文件
  ├── index.html               # HTML入口文件
  ├── .gitignore               # Git忽略文件配置
  ├── package.json             # 项目依赖和脚本配置
  ├── vite.config.js           # Vite配置文件
  
```

## 主要依赖

- Vue 3
- Vite
- Axios
- Vue Router
- Qiniu-js
- 其他...



## 贡献

如果您想为项目贡献代码，请遵循以下步骤：

1. Fork 仓库。
2. 创建新分支。
3. 提交您的更改。
4. 推送您的分支并开启一个Pull Request。

## 许可证



## 致谢

感谢所有贡献者的努力和支持！
