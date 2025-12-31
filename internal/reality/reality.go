package reality

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/curve25519"
	"vless-reality-proxy/internal/config"
	"vless-reality-proxy/pkg/logger"
)

// RealityHandler Reality 处理器
type RealityHandler struct {
	config *config.RealitySettings
	logger logger.Logger
}

// NewRealityHandler 创建 Reality 处理器
func NewRealityHandler(cfg *config.RealitySettings) *RealityHandler {
	return &RealityHandler{
		config: cfg,
		logger: logger.GetLogger(),
	}
}

// WrapConnection 包装连接以支持 Reality
func (r *RealityHandler) WrapConnection(conn net.Conn) (net.Conn, error) {
	// 创建 Reality TLS 配置
	tlsConfig, err := r.createTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("创建 TLS 配置失败: %v", err)
	}

	// 包装为 TLS 连接
	tlsConn := tls.Server(conn, tlsConfig)
	
	// 执行 TLS 握手
	if err := tlsConn.Handshake(); err != nil {
		return nil, fmt.Errorf("TLS 握手失败: %v", err)
	}

	return tlsConn, nil
}

// createTLSConfig 创建 TLS 配置
func (r *RealityHandler) createTLSConfig() (*tls.Config, error) {
	// 生成自签名证书
	cert, err := r.generateSelfSignedCert()
	if err != nil {
		return nil, fmt.Errorf("生成证书失败: %v", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   r.getServerName(),
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
		NextProtos: []string{"h2", "http/1.1"},
	}, nil
}

// generateSelfSignedCert 生成自签名证书
func (r *RealityHandler) generateSelfSignedCert() (tls.Certificate, error) {
	// 这里应该实现真正的证书生成逻辑
	// 为了简化，我们使用一个基本的实现
	
	// 在实际应用中，你需要：
	// 1. 生成密钥对
	// 2. 创建证书模板
	// 3. 自签名证书
	// 4. 返回 tls.Certificate
	
	// 临时返回空证书，实际使用时需要实现完整的证书生成
	return tls.Certificate{}, fmt.Errorf("证书生成功能需要完整实现")
}

// getServerName 获取服务器名称
func (r *RealityHandler) getServerName() string {
	if len(r.config.ServerNames) > 0 {
		return r.config.ServerNames[0]
	}
	return "example.com"
}

// ValidateShortID 验证 Short ID
func (r *RealityHandler) ValidateShortID(shortID string) bool {
	for _, validID := range r.config.ShortIds {
		if shortID == validID {
			return true
		}
	}
	return false
}

// GenerateKeyPair 生成 Curve25519 密钥对
func GenerateKeyPair() (privateKey, publicKey []byte, err error) {
	privateKey = make([]byte, 32)
	if _, err = rand.Read(privateKey); err != nil {
		return nil, nil, err
	}

	publicKey, err = curve25519.X25519(privateKey, curve25519.Basepoint)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}

// EncodePrivateKey 编码私钥为 base64
func EncodePrivateKey(privateKey []byte) string {
	// 这里应该使用适当的编码方式
	// 通常使用 base64 或 hex 编码
	return fmt.Sprintf("%x", privateKey)
}

// DecodePrivateKey 解码私钥
func DecodePrivateKey(encoded string) ([]byte, error) {
	// 对应的解码实现
	privateKey := make([]byte, 32)
	n, err := fmt.Sscanf(encoded, "%x", &privateKey)
	if err != nil || n != 1 {
		return nil, fmt.Errorf("解码私钥失败")
	}
	return privateKey, nil
}

// RealityConfig Reality 配置生成器
type RealityConfig struct {
	Dest         string
	ServerNames  []string
	PrivateKey   string
	ShortIds     []string
	MinClientVer string
	MaxClientVer string
	MaxTimeDiff  int
}

// NewRealityConfig 创建新的 Reality 配置
func NewRealityConfig(dest string, serverNames []string) (*RealityConfig, error) {
	// 生成密钥对
	privateKey, _, err := GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("生成密钥对失败: %v", err)
	}

	// 生成随机 Short ID
	shortIds := make([]string, 3)
	for i := range shortIds {
		shortID := make([]byte, 8)
		rand.Read(shortID)
		shortIds[i] = fmt.Sprintf("%x", shortID)
	}

	return &RealityConfig{
		Dest:         dest,
		ServerNames:  serverNames,
		PrivateKey:   EncodePrivateKey(privateKey),
		ShortIds:     shortIds,
		MinClientVer: "1.8.0",
		MaxClientVer: "",
		MaxTimeDiff:  60000, // 60 秒
	}, nil
}

// Validate 验证 Reality 配置
func (rc *RealityConfig) Validate() error {
	if rc.Dest == "" {
		return fmt.Errorf("目标地址不能为空")
	}

	if len(rc.ServerNames) == 0 {
		return fmt.Errorf("服务器名称列表不能为空")
	}

	if rc.PrivateKey == "" {
		return fmt.Errorf("私钥不能为空")
	}

	if len(rc.ShortIds) == 0 {
		return fmt.Errorf("Short ID 列表不能为空")
	}

	// 验证目标地址格式
	if !strings.Contains(rc.Dest, ":") {
		return fmt.Errorf("目标地址格式无效: %s", rc.Dest)
	}

	return nil
}

// IsRealityConnection 检查是否为 Reality 连接
func IsRealityConnection(conn net.Conn, timeout time.Duration) bool {
	// 设置读取超时
	conn.SetReadDeadline(time.Now().Add(timeout))
	defer conn.SetReadDeadline(time.Time{})

	// 读取 TLS Client Hello
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}

	// 简单检查是否为 TLS 握手
	if n < 6 {
		return false
	}

	// TLS 记录类型 (22 = Handshake)
	if buffer[0] != 22 {
		return false
	}

	// TLS 版本检查
	if buffer[1] != 3 || (buffer[2] != 1 && buffer[2] != 3) {
		return false
	}

	return true
}
