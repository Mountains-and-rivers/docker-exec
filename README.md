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
mkdir -p /usr/local/expireSession/bin
cp -rf expireSession /usr/local/expireSession/bin
chmod +x /usr/local/expireSession/bin/expireSession

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

systemctl daemon-reload
systemctl start expireSession
systemctl stop expireSession
systemctl restart expireSession
systemctl enable expireSession
systemctl status expireSession
journalctl -xe
```
