#!/bin/bash

# 通过命令行参数获取端口号
PORT=$1

if [ -z "$PORT" ]; then
  echo "请提供要查询的端口号"
  echo "示例: ./script.sh 8080"
  exit 1
fi

CHECK_PORT=$(lsof -t -i:$PORT)

if [ -z "$CHECK_PORT" ]; then
  echo "端口 $PORT 未被占用"
else
  echo "端口 $PORT 被进程 $CHECK_PORT 占用"
  kill -9 $CHECK_PORT
  echo "进程 $CHECK_PORT 已被杀死"
fi
