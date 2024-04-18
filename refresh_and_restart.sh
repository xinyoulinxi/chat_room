#!/bin/bash

# 查找web_server进程的PID
pid=$(ps -ef | grep 'web_server' | grep -v 'grep' | awk '{print $2}')

# 判断PID是否存在
if [ -n "$pid" ]; then
    echo "找到web_server进程，PID为：$pid"
    echo "正在杀掉进程..."
    kill -9 $pid
    echo "进程已被杀掉。"
else
    echo "未找到web_server进程。"
fi
echo "更新git仓库..."
git pull -r
echo "更新git仓库完成。"
echo "开始编译..."
./build_and_run_backend.sh
echo "编译完成。已在后台运行中。"
