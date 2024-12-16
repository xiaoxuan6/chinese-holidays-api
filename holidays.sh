#!/bin/bash

set -e

DIR=/root/holidays
SERVICE_FILE="/etc/systemd/system/holidays.service"

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

    echo -e "${GREEN}holidays service install done.${RESET}"
}

function uninstall() {
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

    echo -e "${GREEN}holidays service uninstall successful!${RESET}"
}

case $1 in
install) install "$2" ;;
uninstall) uninstall ;;
reinstall)
    remove
    install
    ;;
*)
    echo "Not found $1 option"
    echo "Usage: $0 {install|uninstall|reinstall}"
    echo ""
    exit 1
    ;;
esac
