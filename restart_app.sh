#!/bin/bash

# 查找 tg-file-share 进程并终止
pkill -f tg-file-share

# 等待进程完全退出
sleep 1

# 启动 tg-file-share
nohup ./tg-file-share > /dev/null 2>&1 &

echo "tg-file-share 已重启"