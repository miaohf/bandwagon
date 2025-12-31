package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"vless-reality-proxy/internal/config"
	"vless-reality-proxy/internal/reality"
	"vless-reality-proxy/pkg/utils"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]
	switch command {
	case "uuid":
		generateUUID()
	case "short-id":
		generateShortID()
	case "reality-keys":
		generateRealityKeys()
	case "full-config":
		generateFullConfig()
	default:
		fmt.Printf("未知命令: %s\n", command)
		showHelp()
	}
}

func showHelp() {
	fmt.Println("配置生成工具")
	fmt.Println("使用方法:")
	fmt.Println("  go run tools/generate-config.go <command>")
	fmt.Println("")
	fmt.Println("可用命令:")
	fmt.Println("  uuid         - 生成新的 UUID")
	fmt.Println("  short-id     - 生成 Reality Short ID")
	fmt.Println("  reality-keys - 生成 Reality 密钥对")
	fmt.Println("  full-config  - 生成完整配置文件")
}

func generateUUID() {
	fmt.Printf("新 UUID: %s\n", uuid.New().String())
}

func generateShortID() {
	fmt.Printf("Short ID: %s\n", utils.GenerateShortID())
}

func generateRealityKeys() {
	privateKey, publicKey, err := reality.GenerateKeyPair()
	if err != nil {
		log.Fatalf("生成密钥对失败: %v", err)
	}

	fmt.Printf("私钥 (服务器用): %s\n", reality.EncodePrivateKey(privateKey))
	fmt.Printf("公钥 (客户端用): %s\n", reality.EncodePrivateKey(publicKey))
}

func generateFullConfig() {
	// 生成客户端配置
	clients := []config.VLESSClient{
		{
			ID:    uuid.New().String(),
			Flow:  "xtls-rprx-vision",
			Email: "user1@example.com",
		},
		{
			ID:    uuid.New().String(),
			Flow:  "xtls-rprx-vision",
			Email: "user2@example.com",
		},
	}

	// 生成 Reality 配置
	realityConfig, err := reality.NewRealityConfig(
		"www.microsoft.com:443",
		[]string{"www.microsoft.com", "microsoft.com"},
	)
	if err != nil {
		log.Fatalf("生成 Reality 配置失败: %v", err)
	}

	// 创建完整配置
	cfg := &config.Config{
		Port:     443,
		LogLevel: "info",
		Inbounds: []config.InboundConfig{
			{
				Protocol: "vless",
				Port:     443,
				Settings: config.VLESSSettings{
					Clients:    clients,
					Decryption: "none",
					Fallbacks: []config.Fallback{
						{
							Dest: "80",
							Xver: 1,
						},
					},
				},
				StreamSettings: &config.StreamSettings{
					Network:  "tcp",
					Security: "reality",
					RealitySettings: &config.RealitySettings{
						Show:         false,
						Dest:         realityConfig.Dest,
						XVer:         0,
						ServerNames:  realityConfig.ServerNames,
						PrivateKey:   realityConfig.PrivateKey,
						MinClientVer: realityConfig.MinClientVer,
						MaxClientVer: realityConfig.MaxClientVer,
						MaxTimeDiff:  realityConfig.MaxTimeDiff,
						ShortIds:     realityConfig.ShortIds,
					},
				},
			},
		},
	}

	// 输出配置
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Fatalf("序列化配置失败: %v", err)
	}

	fmt.Println("生成的配置文件:")
	fmt.Println(string(data))

	// 保存到文件
	if err := os.WriteFile("config-generated.json", data, 0644); err != nil {
		log.Fatalf("保存配置文件失败: %v", err)
	}

	fmt.Println("\n配置已保存到 config-generated.json")
	fmt.Println("\n客户端信息:")
	for i, client := range clients {
		fmt.Printf("客户端 %d UUID: %s\n", i+1, client.ID)
	}
}
