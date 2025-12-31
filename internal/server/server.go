package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"vless-reality-proxy/internal/config"
	"vless-reality-proxy/internal/reality"
	"vless-reality-proxy/internal/vless"
	"vless-reality-proxy/pkg/logger"
)

// Server 代理服务器
type Server struct {
	config   *config.Config
	logger   logger.Logger
	listener net.Listener
	clients  map[string]*config.VLESSClient
	reality  *reality.RealityHandler
	mu       sync.RWMutex
}

// NewServer 创建新的服务器
func NewServer(cfg *config.Config) (*Server, error) {
	server := &Server{
		config:  cfg,
		logger:  logger.GetLogger(),
		clients: make(map[string]*config.VLESSClient),
	}

	// 解析客户端配置
	if err := server.parseClients(); err != nil {
		return nil, fmt.Errorf("解析客户端配置失败: %v", err)
	}

	// 初始化 Reality 处理器
	if err := server.initReality(); err != nil {
		return nil, fmt.Errorf("初始化 Reality 失败: %v", err)
	}

	return server, nil
}

// parseClients 解析客户端配置
func (s *Server) parseClients() error {
	for _, inbound := range s.config.Inbounds {
		if inbound.Protocol == "vless" {
			// 将 settings 转换为 VLESSSettings
			settingsBytes, err := json.Marshal(inbound.Settings)
			if err != nil {
				return fmt.Errorf("序列化 VLESS 设置失败: %v", err)
			}

			var vlessSettings config.VLESSSettings
			if err := json.Unmarshal(settingsBytes, &vlessSettings); err != nil {
				return fmt.Errorf("反序列化 VLESS 设置失败: %v", err)
			}

			// 添加客户端到映射
			for _, client := range vlessSettings.Clients {
				s.clients[client.ID] = &client
			}
		}
	}

	if len(s.clients) == 0 {
		return fmt.Errorf("没有找到有效的 VLESS 客户端配置")
	}

	s.logger.Infof("加载了 %d 个 VLESS 客户端配置", len(s.clients))
	return nil
}

// initReality 初始化 Reality
func (s *Server) initReality() error {
	for _, inbound := range s.config.Inbounds {
		if inbound.StreamSettings != nil && inbound.StreamSettings.RealitySettings != nil {
			s.reality = reality.NewRealityHandler(inbound.StreamSettings.RealitySettings)
			s.logger.Info("Reality 处理器已初始化")
			return nil
		}
	}

	s.logger.Info("未配置 Reality，使用普通 TLS")
	return nil
}

// Start 启动服务器
func (s *Server) Start(ctx context.Context) error {
	// 创建监听器
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return fmt.Errorf("创建监听器失败: %v", err)
	}
	s.listener = listener

	s.logger.Infof("服务器正在监听端口 %d", s.config.Port)

	// 接受连接循环
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("收到关闭信号，停止接受新连接")
			return s.listener.Close()
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				s.logger.Errorf("接受连接失败: %v", err)
				continue
			}

			// 异步处理连接
			go s.handleConnection(conn)
		}
	}
}

// handleConnection 处理连接
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// 设置连接超时
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	s.logger.Debugf("新连接来自: %s", conn.RemoteAddr())

	// 如果配置了 Reality，先处理 Reality
	if s.reality != nil {
		wrappedConn, err := s.reality.WrapConnection(conn)
		if err != nil {
			s.logger.Errorf("Reality 处理失败: %v", err)
			return
		}
		conn = wrappedConn
	}

	// 处理 VLESS 协议
	if err := s.handleVLESS(conn); err != nil {
		s.logger.Errorf("处理 VLESS 连接失败: %v", err)
	}
}

// handleVLESS 处理 VLESS 连接
func (s *Server) handleVLESS(conn net.Conn) error {
	// 创建 VLESS 连接
	vlessConn := vless.NewVLESSConnection(conn)

	// 读取 VLESS 头部
	header, err := vlessConn.ReadHeader()
	if err != nil {
		return fmt.Errorf("读取 VLESS 头部失败: %v", err)
	}

	// 验证客户端 UUID
	if !s.validateClient(header.GetUUID()) {
		return fmt.Errorf("无效的客户端 UUID: %s", header.GetUUID())
	}

	s.logger.Debugf("客户端 %s 请求连接到 %s", header.GetUUID(), header.GetDestination())

	// 连接到目标服务器
	targetConn, err := net.DialTimeout("tcp", header.GetDestination(), 10*time.Second)
	if err != nil {
		return fmt.Errorf("连接目标服务器失败: %v", err)
	}
	defer targetConn.Close()

	// 发送 VLESS 响应
	if err := vlessConn.WriteResponse(); err != nil {
		return fmt.Errorf("发送 VLESS 响应失败: %v", err)
	}

	// 开始数据中继
	s.logger.Debugf("开始为客户端 %s 中继数据", header.GetUUID())
	return vless.Relay(targetConn, conn)
}

// validateClient 验证客户端
func (s *Server) validateClient(uuid string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	_, exists := s.clients[uuid]
	return exists
}

// GetStats 获取服务器统计信息
func (s *Server) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"clients_count": len(s.clients),
		"port":         s.config.Port,
		"reality_enabled": s.reality != nil,
	}
}

// AddClient 添加客户端
func (s *Server) AddClient(client *config.VLESSClient) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.clients[client.ID] = client
	s.logger.Infof("添加新客户端: %s", client.ID)
}

// RemoveClient 移除客户端
func (s *Server) RemoveClient(uuid string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.clients[uuid]; exists {
		delete(s.clients, uuid)
		s.logger.Infof("移除客户端: %s", uuid)
		return true
	}
	return false
}

// GetClients 获取客户端列表
func (s *Server) GetClients() []*config.VLESSClient {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	clients := make([]*config.VLESSClient, 0, len(s.clients))
	for _, client := range s.clients {
		clients = append(clients, client)
	}
	return clients
}

// Stop 停止服务器
func (s *Server) Stop() error {
	if s.listener != nil {
		s.logger.Info("正在关闭服务器...")
		return s.listener.Close()
	}
	return nil
}
