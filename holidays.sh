#!/bin/bash

set -e

DIR=/root/holidays
SERVICE_FILE="/etc/systemd/system/holidays.service"
NGINX_FILE="/etc/nginx/modules-enabled/holidays.conf"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | tr '[:upper:]' '[:lower:]')

GREEN='\e[32m'
RED='\e[31m'
RESET='\e[0m'

function get_arch() {
    local arch=""
    case $ARCH in
    x86_64)
        arch=amd64
        ;;
    i386)
        arch=386
        ;;
    *)
        arch=$ARCH
        ;;
    esac
    echo "$arch"
}

function add_systemd() {
    cat <<EOL | sudo tee "$SERVICE_FILE" >/dev/null
[Unit]
Description=holidays api Service
After=network.target

[Service]
Environment=HOLIDAY_PORT=$1
ExecStart=$DIR/holidays
WorkingDirectory=$DIR
Type=simple
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOL

    systemctl daemon-reload
    systemctl start holidays.service
    systemctl enable holidays.service
}

function add_nginx() {
    if [ -x "$(command -v nginx)" ]; then
        local nginxFile="/etc/nginx/modules-enabled"
        if [ -d "$nginxFile" ]; then
            cat <<EOL | sudo tee "$NGINX_FILE" >/dev/null
stream {
    upstream holidays_api {
        server 127.0.0.1:$2;
    }

    server {
        listen $1;
        proxy_pass holidays_api;
    }
}
EOL
            systemctl reload nginx
        else
            echo -e "${RED}nginx 配置文件夹[$nginxFile]不存在！请手动配置nginx反向代理${RESET}"
            exit 1
        fi
    else
        echo -e "${RED}nginx service not exist!${RESET}"
        exit 1
    fi
}

function install() {
    mkdir -p "$DIR"
    cd "$DIR" || exit 1

    arch=$(get_arch)
    UNAME="${OS}_${arch}"
    URL=$(curl -s https://api.github.com/repos/xiaoxuan6/chinese-holidays-api/releases/latest | grep "browser_download_url" | grep "tar.gz" | cut -d '"' -f 4 | grep "$UNAME")
    if [ -z "$URL" ]; then
        echo -e "${RED}Unsupported platform: $(uname -s) $(uname -m)${RESET}"
        exit 1
    fi

    curl -L -O "$URL"
    FILENAME=$(echo "$URL" | cut -d '/' -f 9)
    if [ ! -f "$FILENAME" ]; then
        echo "url: $URL"
        echo -e "${RED}filename $FILENAME dose not exist${RESET}"
        exit 1
    fi

    tar xf "$FILENAME"
    rm "$FILENAME" "LICENSE" "README.md"
    mv chinese-holidays-api holidays

    if [ -z "$1" ]; then
        # shellcheck disable=SC2162
        read -p "请输入 holidays 服务有效端口号：" port
    else
        port=$1
    fi

    if [[ "$port" =~ ^[0-9]+$ ]]; then
        add_systemd "$port"
    else
        echo "无效的端口号: $port, 端口必须是数字"
        exit 1
    fi

    # shellcheck disable=SC2162
    read -p "是否设置nginx反向代理?(y/n)" disable_port
    if [[ "$disable_port" =~ ^[Yy]$ ]]; then
        # shellcheck disable=SC2155
        local new_port=$((port + 1))
        add_nginx "$new_port" "$port"
        echo -e "${GREEN}反向代理设置成功，代理地址：127.0.0.1:$new_port${RESET}"
    fi

    echo "holidays service install done."
}

function remove() {
    systemctl stop holidays.service
    systemctl disable holidays.service

    if [ -f "$SERVICE_FILE" ]; then
        rm "$SERVICE_FILE"
        echo "服务文件 $SERVICE_FILE 已删除"
    else
        echo "服务文件 $SERVICE_FILE 不存在！"
    fi

    systemctl daemon-reload

    if [ -d "$DIR" ]; then
        rm -rf "$DIR"
        echo "文件夹 $DIR 已删除"
    else
        echo "文件夹 $DIR 不存在！"
    fi

    if [ -f "$NGINX_FILE" ]; then
        rm "$NGINX_FILE"
        systemctl reload nginx
    fi

    echo -e "${GREEN}holidays service remove successful!${RESET}"
}

case $1 in
install) install "$2" ;;
remove) remove ;;
*)
    echo "Not found $1 option"
    echo "Usage: $0 {install|remove}"
    echo ""
    exit 1
    ;;
esac
