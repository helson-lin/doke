#!/bin/bash

# 定义应用程序的名称
APP_NAME="doke"

# 定义应用程序的版本号
VERSION=$1  # 你可以根据需要修改版本号

# 定义要构建的平台
PLATFORMS=("darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64" "windows/amd64")

# 清除之前的构建并创建新的 release 目录
rm -rf release
mkdir release

export GIN_MODE=release

# 为每个平台构建
for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    OUTPUT_NAME=$APP_NAME_$VERSION_$GOOS_$GOARCH  # 在文件名中加入版本号
    if [ $GOOS = "windows" ]; then
        OUTPUT_NAME+='.exe'
    fi

    echo "Building for $GOOS $GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X main.Version=$VERSION" -o release/$OUTPUT_NAME
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    # 将构建的二进制文件打包为 tar.gz
    echo "Compressing $OUTPUT_NAME..."
    tar -czvf release/$OUTPUT_NAME.tar.gz -C release $OUTPUT_NAME
    if [ $? -ne 0 ]; then
        echo 'An error occurred while compressing! Aborting the script execution...'
        exit 1
    fi

    # 删除原始的二进制文件（可选）
    rm release/$OUTPUT_NAME
done
