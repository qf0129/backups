#!/bin/bash
set -e

# 编译Go程序
go build -o backups

# 打包到zip文件
zip backups-linux-amd-$(git describe --tags).zip backups conf.example.json

echo "编译并打包完成"