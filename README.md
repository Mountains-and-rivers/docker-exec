## docker容器中执行命令

### 版本信息

| software | version |
| -------- | ------- |
| golang   | 1.16.3  |
| docker   | 19.03.9 |

### 编译

```
cd docker-exec
go build -o expireSession
```

### system管理golang进程

```
 
cat > /usr/lib/systemd/system/expireSession.service << EOF
[Unit]
Description=expireSession
After=syslog.target network.target
[Service]
User=root
Type=simple
Environment=AWS_SHARED_CREDENTIALS_FILE=/home/username/.aws/credentials
# 环境变量
Environment="SGFOOT_ENV=pro"
Environment="SGFOOT_PATH=/data/conf"
Restart=on-failure
RestartSec=5s
WorkingDirectory=/usr/local/expireSession/bin
ExecStart=nohup /usr/local/expireSession/bin/expireSession >/dev/null 2>&1 &

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload # 更新配置
systemctl start expireSession # 启动
systemctl stop expireSession # 停止
systemctl restart expireSession # 重启
systemctl enable expireSession # 加入开机启动
systemctl status expireSession #查看状态
journalctl -xe #查看日志
```
