# Reality 域名选择指南

## 🎯 什么是 Reality 回落域名

Reality 回落域名（`dest` 字段）是 Reality 协议的核心组件，它的作用是：

1. **流量伪装**：让你的代理服务器看起来像在访问这个网站
2. **抗检测**：当有探测流量时，会被转发到真实网站
3. **TLS 特征模拟**：模拟真实网站的 TLS 握手特征

## 📋 域名选择标准

### ✅ 好的域名特征

1. **高可用性**
   - 99.9%+ 的在线时间
   - 全球 CDN 分布
   - 快速响应时间

2. **TLS 配置良好**
   - 支持 TLS 1.3
   - 现代密码套件
   - 有效的证书链

3. **流量特征普通**
   - 大众化网站
   - 正常的访问模式
   - 不敏感的内容

4. **地理位置合适**
   - 与服务器地理位置相近
   - 低延迟连接
   - 稳定的路由

### ❌ 避免的域名特征

1. **政治敏感**
   - 政府网站
   - 新闻媒体
   - 社交平台（在某些地区）

2. **技术特征明显**
   - VPN/代理服务商
   - 技术论坛
   - 开发者工具

3. **不稳定**
   - 经常宕机
   - 证书过期
   - 配置变更频繁

## 🏆 推荐域名列表

### 一线推荐（最佳选择）

```json
{
  "dest": "www.microsoft.com:443",
  "server_names": ["www.microsoft.com", "microsoft.com"]
}
```

```json
{
  "dest": "www.apple.com:443",
  "server_names": ["www.apple.com", "apple.com"]
}
```

```json
{
  "dest": "www.cloudflare.com:443",
  "server_names": ["www.cloudflare.com", "cloudflare.com"]
}
```

### 二线推荐（备选方案）

**电商平台**：
- `www.amazon.com:443`
- `www.ebay.com:443`
- `www.shopify.com:443`

**技术公司**：
- `github.com:443`
- `stackoverflow.com:443`
- `www.docker.com:443`

**云服务**：
- `aws.amazon.com:443`
- `cloud.google.com:443`
- `azure.microsoft.com:443`

**娱乐平台**：
- `www.youtube.com:443`
- `www.netflix.com:443`
- `www.spotify.com:443`

## 🛠️ 域名测试工具

### 使用内置测试工具

```bash
# 测试所有推荐域名
make test-domains

# 或者直接运行
go run tools/test-domains.go
```

### 手动测试延迟

```bash
# 测试连接延迟
for domain in www.microsoft.com www.apple.com www.cloudflare.com; do
  echo "测试 $domain:"
  time openssl s_client -connect $domain:443 -servername $domain < /dev/null
  echo "---"
done
```

### 测试 TLS 配置

```bash
# 检查 TLS 版本和密码套件
openssl s_client -connect www.microsoft.com:443 -servername www.microsoft.com -tls1_3
```

## 🌍 地区化建议

### 北美地区
- `www.microsoft.com:443`
- `www.apple.com:443`
- `www.amazon.com:443`

### 欧洲地区
- `www.microsoft.com:443`
- `www.apple.com:443`
- `www.cloudflare.com:443`

### 亚太地区
- `www.microsoft.com:443`
- `www.apple.com:443`
- `aws.amazon.com:443`

## ⚙️ 配置示例

### 基础配置

```json
{
  "reality_settings": {
    "show": false,
    "dest": "www.microsoft.com:443",
    "xver": 0,
    "server_names": [
      "www.microsoft.com",
      "microsoft.com"
    ],
    "private_key": "your-private-key-here",
    "short_ids": ["your-short-id-here"]
  }
}
```

### 多域名配置

```json
{
  "reality_settings": {
    "show": false,
    "dest": "www.microsoft.com:443",
    "xver": 0,
    "server_names": [
      "www.microsoft.com",
      "microsoft.com",
      "docs.microsoft.com",
      "azure.microsoft.com"
    ],
    "private_key": "your-private-key-here",
    "short_ids": ["short-id-1", "short-id-2"]
  }
}
```

## 🔄 域名轮换策略

### 定期更换

建议每 2-3 个月更换一次域名：

1. 测试新域名的性能
2. 更新服务器配置
3. 通知客户端更新
4. 监控连接稳定性

### 多域名部署

在不同服务器上使用不同域名：

```bash
# 服务器 A
"dest": "www.microsoft.com:443"

# 服务器 B  
"dest": "www.apple.com:443"

# 服务器 C
"dest": "www.cloudflare.com:443"
```

## 🚨 安全注意事项

1. **避免使用相同域名**
   - 不要在多个服务器上使用相同的回落域名
   - 定期更换域名避免特征识别

2. **监控域名状态**
   - 定期检查域名可用性
   - 关注证书更新
   - 监控访问延迟

3. **备用方案**
   - 准备 2-3 个备用域名
   - 测试备用域名的可用性
   - 制定快速切换方案

## 📊 性能监控

### 监控指标

- **连接延迟**：< 200ms
- **成功率**：> 99%
- **TLS 握手时间**：< 1s
- **证书有效期**：> 30天

### 监控脚本

```bash
#!/bin/bash
# reality-monitor.sh

DOMAIN="www.microsoft.com:443"
THRESHOLD=200  # ms

LATENCY=$(timeout 5 bash -c "time openssl s_client -connect $DOMAIN -servername ${DOMAIN%:*} < /dev/null" 2>&1 | grep real | awk '{print $2}')

if [[ $LATENCY > ${THRESHOLD}ms ]]; then
    echo "警告：域名 $DOMAIN 延迟过高: $LATENCY"
    # 发送告警
fi
```

## 🔧 故障排除

### 常见问题

1. **连接超时**
   ```bash
   # 检查网络连通性
   telnet www.microsoft.com 443
   ```

2. **TLS 握手失败**
   ```bash
   # 检查 TLS 配置
   openssl s_client -connect www.microsoft.com:443 -servername www.microsoft.com
   ```

3. **证书验证失败**
   ```bash
   # 检查证书链
   openssl s_client -connect www.microsoft.com:443 -servername www.microsoft.com -verify_return_error
   ```

### 解决方案

1. **更换域名**：选择延迟更低的域名
2. **调整配置**：优化 TLS 设置
3. **网络优化**：检查路由和 DNS 设置

记住：选择合适的回落域名是 Reality 协议成功的关键！
