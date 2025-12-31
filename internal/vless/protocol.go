package vless

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/google/uuid"
)

const (
	// VLESS 版本
	Version = 0

	// 地址类型
	AddressTypeIPv4   = 1
	AddressTypeDomain = 2
	AddressTypeIPv6   = 3

	// 命令类型
	CommandTCP = 1
	CommandUDP = 2
	CommandMux = 3
)

// VLESSHeader VLESS 协议头
type VLESSHeader struct {
	Version    byte
	UUID       [16]byte
	AddOnLen   byte
	Command    byte
	Port       uint16
	AddressType byte
	Address     []byte
}

// VLESSConnection VLESS 连接
type VLESSConnection struct {
	conn   net.Conn
	header *VLESSHeader
}

// NewVLESSConnection 创建新的 VLESS 连接
func NewVLESSConnection(conn net.Conn) *VLESSConnection {
	return &VLESSConnection{
		conn: conn,
	}
}

// ReadHeader 读取 VLESS 头部
func (vc *VLESSConnection) ReadHeader() (*VLESSHeader, error) {
	// 读取版本
	versionBuf := make([]byte, 1)
	if _, err := io.ReadFull(vc.conn, versionBuf); err != nil {
		return nil, fmt.Errorf("读取版本失败: %v", err)
	}

	if versionBuf[0] != Version {
		return nil, fmt.Errorf("不支持的 VLESS 版本: %d", versionBuf[0])
	}

	header := &VLESSHeader{
		Version: versionBuf[0],
	}

	// 读取 UUID
	if _, err := io.ReadFull(vc.conn, header.UUID[:]); err != nil {
		return nil, fmt.Errorf("读取 UUID 失败: %v", err)
	}

	// 读取附加数据长度
	addOnLenBuf := make([]byte, 1)
	if _, err := io.ReadFull(vc.conn, addOnLenBuf); err != nil {
		return nil, fmt.Errorf("读取附加数据长度失败: %v", err)
	}
	header.AddOnLen = addOnLenBuf[0]

	// 跳过附加数据
	if header.AddOnLen > 0 {
		addOnData := make([]byte, header.AddOnLen)
		if _, err := io.ReadFull(vc.conn, addOnData); err != nil {
			return nil, fmt.Errorf("读取附加数据失败: %v", err)
		}
	}

	// 读取命令
	commandBuf := make([]byte, 1)
	if _, err := io.ReadFull(vc.conn, commandBuf); err != nil {
		return nil, fmt.Errorf("读取命令失败: %v", err)
	}
	header.Command = commandBuf[0]

	// 读取端口
	portBuf := make([]byte, 2)
	if _, err := io.ReadFull(vc.conn, portBuf); err != nil {
		return nil, fmt.Errorf("读取端口失败: %v", err)
	}
	header.Port = binary.BigEndian.Uint16(portBuf)

	// 读取地址类型
	addressTypeBuf := make([]byte, 1)
	if _, err := io.ReadFull(vc.conn, addressTypeBuf); err != nil {
		return nil, fmt.Errorf("读取地址类型失败: %v", err)
	}
	header.AddressType = addressTypeBuf[0]

	// 读取地址
	switch header.AddressType {
	case AddressTypeIPv4:
		header.Address = make([]byte, 4)
		if _, err := io.ReadFull(vc.conn, header.Address); err != nil {
			return nil, fmt.Errorf("读取 IPv4 地址失败: %v", err)
		}
	case AddressTypeIPv6:
		header.Address = make([]byte, 16)
		if _, err := io.ReadFull(vc.conn, header.Address); err != nil {
			return nil, fmt.Errorf("读取 IPv6 地址失败: %v", err)
		}
	case AddressTypeDomain:
		domainLenBuf := make([]byte, 1)
		if _, err := io.ReadFull(vc.conn, domainLenBuf); err != nil {
			return nil, fmt.Errorf("读取域名长度失败: %v", err)
		}
		domainLen := domainLenBuf[0]
		header.Address = make([]byte, domainLen)
		if _, err := io.ReadFull(vc.conn, header.Address); err != nil {
			return nil, fmt.Errorf("读取域名失败: %v", err)
		}
	default:
		return nil, fmt.Errorf("不支持的地址类型: %d", header.AddressType)
	}

	vc.header = header
	return header, nil
}

// WriteResponse 写入 VLESS 响应
func (vc *VLESSConnection) WriteResponse() error {
	// VLESS 响应格式很简单，只需要写入版本和附加数据长度
	response := []byte{Version, 0} // 版本 + 附加数据长度 (0)
	_, err := vc.conn.Write(response)
	return err
}

// GetDestination 获取目标地址
func (header *VLESSHeader) GetDestination() string {
	var address string
	
	switch header.AddressType {
	case AddressTypeIPv4:
		address = net.IP(header.Address).String()
	case AddressTypeIPv6:
		address = net.IP(header.Address).String()
	case AddressTypeDomain:
		address = string(header.Address)
	}

	return net.JoinHostPort(address, strconv.Itoa(int(header.Port)))
}

// GetUUID 获取 UUID 字符串
func (header *VLESSHeader) GetUUID() string {
	u, _ := uuid.FromBytes(header.UUID[:])
	return u.String()
}

// ValidateUUID 验证 UUID
func ValidateUUID(headerUUID [16]byte, validUUIDs []string) bool {
	headerUUIDStr, _ := uuid.FromBytes(headerUUID[:])
	headerUUIDString := headerUUIDStr.String()

	for _, validUUID := range validUUIDs {
		if headerUUIDString == validUUID {
			return true
		}
	}
	return false
}

// GenerateUUID 生成新的 UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// Relay 数据中继
func Relay(dst, src net.Conn) error {
	defer dst.Close()
	defer src.Close()

	errCh := make(chan error, 2)

	// 双向数据传输
	go func() {
		_, err := io.Copy(dst, src)
		errCh <- err
	}()

	go func() {
		_, err := io.Copy(src, dst)
		errCh <- err
	}()

	// 等待任意一个方向的传输完成
	return <-errCh
}
