#!/bin/bash

echo "启动后端服务器..."
cd server
export GOPROXY=https://proxy.golang.org,direct
go run main.go
