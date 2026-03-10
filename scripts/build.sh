#!/bin/bash

# 构建脚本，支持构建不同平台的可执行文件
# 用法：
# ./build.sh                # 构建所有平台
# ./build.sh darwin         # 只构建 macOS 版本
# ./build.sh linux          # 只构建 Linux 版本（包含 amd64 和 arm64）
# ./build.sh linux-amd64    # 只构建 Linux AMD64 版本
# ./build.sh linux-arm64    # 只构建 Linux ARM64 版本
# ./build.sh windows        # 只构建 Windows 版本

# 定义函数：构建指定平台
build_platform() {
  local os=$1
  local arch=$2
  local ext=$3
  local name=$4
  
  echo "构建 ${name} 版本..."
  export GOOS=${os}
  export GOARCH=${arch}
  go build -o output/wechat-publisher-${os}-${arch}${ext} ./cmd/publisher
}

# 检查参数
PLATFORM="all"
if [ $# -gt 0 ]; then
  PLATFORM=$1
fi

echo "开始构建微信发布工具..."

# 创建输出目录
mkdir -p output

# 构建指定平台
case ${PLATFORM} in
  "darwin")
    build_platform "darwin" "amd64" "" "macOS AMD64"
    build_platform "darwin" "arm64" "" "macOS ARM64"
    ;;
  "linux")
    build_platform "linux" "amd64" "" "Linux AMD64"
    build_platform "linux" "arm64" "" "Linux ARM64"
    ;;
  "linux-amd64")
    build_platform "linux" "amd64" "" "Linux AMD64"
    ;;
  "linux-arm64")
    build_platform "linux" "arm64" "" "Linux ARM64"
    ;;
  "windows")
    build_platform "windows" "amd64" ".exe" "Windows AMD64"
    ;;
  "all")
    # 构建 macOS 版本
    build_platform "darwin" "amd64" "" "macOS AMD64"
    build_platform "darwin" "arm64" "" "macOS ARM64"
    
    # 构建 Linux 版本
    build_platform "linux" "amd64" "" "Linux AMD64"
    build_platform "linux" "arm64" "" "Linux ARM64"
    
    # 构建 Windows 版本
    build_platform "windows" "amd64" ".exe" "Windows AMD64"
    ;;
  *)
    echo "错误：不支持的平台 ${PLATFORM}"
    echo "支持的平台：darwin, linux, linux-amd64, linux-arm64, windows, all"
    exit 1
    ;;
esac

# 复制配置文件和 demo.html 到输出目录
echo "复制配置文件和 demo.html 到输出目录..."
mkdir -p output/configs
cp -r configs/* output/configs/
cp demo.html output/
cp API_GUIDE.md output/

echo "构建完成！可执行文件和配置文件已输出到 output 目录。"
echo ""
echo "可用的可执行文件："
ls -la output/
