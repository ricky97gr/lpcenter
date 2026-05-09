#!/bin/sh
echo "======================================"
echo "  lpcenter 服务启动"
echo "======================================"

echo ""
echo "[1] 测试 Nginx 配置文件..."
/opt/bitnami/nginx/sbin/nginx -t
if [ $? -ne 0 ]; then
    echo "Nginx配置错误，退出"
    exit 1
fi

echo ""
echo "[2] 启动 Nginx (端口 9090)..."
/opt/bitnami/nginx/sbin/nginx
sleep 1
ps aux | grep nginx

echo ""
echo "[3] 启动后端服务 (端口 9091, 下载端口 9092)..."
export HTTP_PORT=9091
echo "访问地址: http://localhost:9090"
echo "======================================"

cd /app
exec ./lpcenter_server