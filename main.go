package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vless-reality-proxy/internal/config"
	"vless-reality-proxy/internal/server"
	"vless-reality-proxy/pkg/logger"
)

var (
	configPath = flag.String("config", "config.json", "配置文件路径")
	version    = flag.Bool("version", false, "显示版本信息")
)

const VERSION = "1.0.0"

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("VLESS Reality Proxy v%s\n", VERSION)
		return
	}

	// 初始化日志
	logger.Init()
	log := logger.GetLogger()

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建服务器
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("创建服务器失败: %v", err)
	}

	// 启动服务器
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := srv.Start(ctx); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	log.Printf("VLESS Reality 代理服务器已启动，监听端口: %d", cfg.Port)

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("正在关闭服务器...")
	cancel()

	// 给服务器一些时间来优雅关闭
	time.Sleep(2 * time.Second)
	log.Println("服务器已关闭")
}
