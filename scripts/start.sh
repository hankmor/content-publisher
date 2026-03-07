#!/bin/bash

# 启动脚本，用于启动微信发布工具

echo "启动微信发布工具..."

# 检查可执行文件是否存在
if [ -f "./wechat-publisher" ]; then
    EXECUTABLE="./wechat-publisher"
elif [ -f "./wechat-publisher-darwin-amd64" ]; then
    EXECUTABLE="./wechat-publisher-darwin-amd64"
elif [ -f "./wechat-publisher-linux-amd64" ]; then
    EXECUTABLE="./wechat-publisher-linux-amd64"
elif [ -f "./wechat-publisher-windows-amd64.exe" ]; then
    EXECUTABLE="./wechat-publisher-windows-amd64.exe"
elif [ -f "./output/wechat-publisher-darwin-amd64" ]; then
    EXECUTABLE="./output/wechat-publisher-darwin-amd64"
elif [ -f "./output/wechat-publisher-linux-amd64" ]; then
    EXECUTABLE="./output/wechat-publisher-linux-amd64"
elif [ -f "./output/wechat-publisher-windows-amd64.exe" ]; then
    EXECUTABLE="./output/wechat-publisher-windows-amd64.exe"
else
    echo "错误：未找到可执行文件！"
    echo "请先运行 build.sh 脚本构建应用。"
    exit 1
fi

# 检查配置文件是否存在
if [ ! -f "./configs/config.yaml" ]; then
    echo "警告：配置文件不存在！"
    echo "请根据 configs/config.yaml.example 创建 configs/config.yaml 文件并填写相关配置。"
    if [ -f "./configs/config.yaml.example" ]; then
        echo "正在复制配置文件模板..."
        cp ./configs/config.yaml.example ./configs/config.yaml
        echo "请编辑 configs/config.yaml 文件，填写微信公众号的 AppID、AppSecret 和 Token。"
    fi
fi

# 启动应用
echo "正在启动应用..."
$EXECUTABLE
