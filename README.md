# VLESS Reality 代理服务器

一个支持 VLESS 协议和 Reality TLS 伪装的高性能代理服务器，使用 Go 语言开发。

## 功能特性

- ✅ **VLESS 协议支持**: 完整实现 VLESS 协议，支持 TCP 传输
- ✅ **Reality TLS 伪装**: 支持 Reality 协议，提供更好的抗检测能力
- ✅ **高性能**: 基于 Go 语言开发，支持高并发连接
- ✅ **灵活配置**: JSON 格式配置文件，支持多客户端管理
- ✅ **安全可靠**: 支持 UUID 验证和流量加密
- ✅ **日志系统**: 完整的日志记录和错误处理

## 项目结构

```
vless-reality-proxy/
├── main.go                 # 主程序入口
├── go.mod                  # Go 模块配置
├── config.json             # 示例配置文件
├── internal/               # 内部包
│   ├── config/            # 配置管理
│   │   └── config.go
│   ├── server/            # 服务器核心
│   │   └── server.go
│   ├── vless/             # VLESS 协议实现
│   │   └── protocol.go
│   └── reality/           # Reality TLS 实现
│       └── reality.go
└── pkg/                   # 公共包
    ├── logger/            # 日志系统
    │   └── logger.go
    └── utils/             # 工具函数
        └── utils.go
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置服务器

编辑 `config.json` 文件：

```json
{
  "port": 443,
  "log_level": "info",
  "inbounds": [
    {
      "protocol": "vless",
      "port": 443,
      "settings": {
        "clients": [
          {
            "id": "b831381d-6324-4d53-ad4f-8cda48b30811",
            "flow": "xtls-rprx-vision",
            "email": "user1@example.com"
          }
        ],
        "decryption": "none"
      },
      "stream_settings": {
        "network": "tcp",
        "security": "reality",
        "reality_settings": {
          "show": false,
          "dest": "www.microsoft.com:443",
          "server_names": ["www.microsoft.com"],
          "private_key": "your-private-key-here",
          "short_ids": ["6ba85179e30d4fc2"]
        }
      }
    }
  ]
}
```

### 3. 运行服务器

```bash
# 使用默认配置
go run main.go

# 指定配置文件
go run main.go -config=/path/to/config.json

# 查看版本信息
go run main.go -version
```

### 4. 编译可执行文件

```bash
# 编译当前平台
go build -o vless-reality-proxy main.go

# 交叉编译 Linux
GOOS=linux GOARCH=amd64 go build -o vless-reality-proxy-linux main.go

# 交叉编译 Windows
GOOS=windows GOARCH=amd64 go build -o vless-reality-proxy.exe main.go
```

## 配置说明

### 基本配置

- `port`: 服务器监听端口
- `log_level`: 日志级别 (debug, info, warn, error)

### VLESS 客户端配置

每个客户端需要配置：
- `id`: 客户端 UUID (可使用在线工具生成)
- `flow`: 流控模式 (推荐 `xtls-rprx-vision`)
- `email`: 客户端标识 (可选)

### Reality 配置

- `dest`: 回落目标地址 (如 `www.microsoft.com:443`)
- `server_names`: 服务器名称列表
- `private_key`: Reality 私钥
- `short_ids`: 短 ID 列表 (用于客户端识别)

## 生成配置工具

### 生成 UUID

```go
package main

import (
    "fmt"
    "github.com/google/uuid"
)

func main() {
    fmt.Println(uuid.New().String())
}
```

### 生成 Reality 密钥

```bash
# 使用项目内置工具
go run -c 'import "vless-reality-proxy/pkg/utils"; fmt.Println(utils.GenerateUUID())'
```

## 客户端配置示例

### V2Ray/Xray 客户端配置

```json
{
  "outbounds": [
    {
      "protocol": "vless",
      "settings": {
        "vnext": [
          {
            "address": "your-server-ip",
            "port": 443,
            "users": [
              {
                "id": "b831381d-6324-4d53-ad4f-8cda48b30811",
                "flow": "xtls-rprx-vision"
              }
            ]
          }
        ]
      },
      "streamSettings": {
        "network": "tcp",
        "security": "reality",
        "realitySettings": {
          "serverName": "www.microsoft.com",
          "fingerprint": "chrome",
          "shortId": "6ba85179e30d4fc2",
          "publicKey": "your-public-key-here"
        }
      }
    }
  ]
}
```

## 性能优化

### 系统调优

```bash
# 增加文件描述符限制
echo "* soft nofile 65535" >> /etc/security/limits.conf
echo "* hard nofile 65535" >> /etc/security/limits.conf

# 优化网络参数
echo "net.core.rmem_max = 134217728" >> /etc/sysctl.conf
echo "net.core.wmem_max = 134217728" >> /etc/sysctl.conf
echo "net.ipv4.tcp_rmem = 4096 65536 134217728" >> /etc/sysctl.conf
echo "net.ipv4.tcp_wmem = 4096 65536 134217728" >> /etc/sysctl.conf

sysctl -p
```

## 部署建议

### 使用 systemd 服务

创建服务文件 `/etc/systemd/system/vless-reality-proxy.service`:

```ini
[Unit]
Description=VLESS Reality Proxy Server
After=network.target

[Service]
Type=simple
User=nobody
Group=nobody
ExecStart=/usr/local/bin/vless-reality-proxy -config=/etc/vless-reality-proxy/config.json
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl enable vless-reality-proxy
sudo systemctl start vless-reality-proxy
sudo systemctl status vless-reality-proxy
```

## 安全注意事项

1. **定期更换密钥**: 建议定期更换 Reality 私钥和客户端 UUID
2. **防火墙配置**: 只开放必要的端口
3. **日志管理**: 定期清理日志文件，避免敏感信息泄露
4. **证书更新**: 如果使用自定义证书，需要定期更新

## 故障排除

### 常见问题

1. **连接失败**
   - 检查服务器端口是否正确开放
   - 验证客户端 UUID 是否匹配
   - 确认 Reality 配置是否正确

2. **性能问题**
   - 检查系统资源使用情况
   - 优化网络参数
   - 考虑使用更高性能的服务器

3. **日志错误**
   - 查看详细错误信息
   - 检查配置文件格式
   - 验证网络连通性

## 开发

### 构建要求

- Go 1.21+
- Linux/macOS/Windows

### 依赖包

- `github.com/google/uuid`: UUID 生成
- `github.com/gorilla/websocket`: WebSocket 支持
- `golang.org/x/crypto`: 加密算法
- `golang.org/x/net`: 网络工具

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 免责声明

本项目仅供学习和研究使用，请遵守当地法律法规。
