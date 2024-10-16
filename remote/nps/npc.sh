#!/bin/sh

# 创建服务文件
cat  >/lib/systemd/system/nps.service<<EOF
[Unit]
Description=NPC Service
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/cache/nps/npc -server=ip:端口 -vkey=密钥 -type=tcp
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target

EOF


# 重新加载Systemd，启用并启动服务
systemctl daemon-reload
systemctl enable nps.service
systemctl start nps.service

echo "nps服务已创建并启动。"
