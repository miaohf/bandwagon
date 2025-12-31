package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

// GeographyAnalysis 地理位置分析
type GeographyAnalysis struct {
	Domain      string
	IP          string
	Latency     time.Duration
	Hops        int
	Region      string
	Suspicious  bool
	Reason      string
}

// 不同地区的推荐域名
var regionalDomains = map[string][]string{
	"北美": {
		"www.microsoft.com:443",
		"www.apple.com:443",
		"www.amazon.com:443",
		"github.com:443",
		"www.netflix.com:443",
	},
	"欧洲": {
		"www.microsoft.com:443",
		"www.apple.com:443",
		"www.spotify.com:443",
		"www.booking.com:443",
		"www.bbc.com:443",
	},
	"亚太": {
		"www.microsoft.com:443",
		"www.apple.com:443",
		"aws.amazon.com:443",
		"www.sony.com:443",
		"www.samsung.com:443",
	},
	"中国": {
		"www.microsoft.com:443",  // 微软在中国有节点
		"www.apple.com:443",      // 苹果在中国有CDN
		"github.com:443",         // GitHub虽然慢但常用
		"stackoverflow.com:443",  // 开发者常访问
		"www.docker.com:443",     // 技术相关
	},
}

func main() {
	fmt.Println("🌍 Reality 域名地理位置分析工具")
	fmt.Println("分析不同域名的地理位置特征和访问合理性")
	fmt.Println()

	// 检测当前服务器位置（简化版）
	serverRegion := detectServerRegion()
	fmt.Printf("检测到服务器可能位于: %s\n", serverRegion)
	fmt.Println()

	// 分析推荐域名
	fmt.Printf("📊 %s地区推荐域名分析:\n", serverRegion)
	recommendedDomains := regionalDomains[serverRegion]
	
	var analyses []GeographyAnalysis
	for _, domain := range recommendedDomains {
		analysis := analyzeDomain(domain, serverRegion)
		analyses = append(analyses, analysis)
		
		status := "✅"
		if analysis.Suspicious {
			status = "⚠️"
		}
		
		fmt.Printf("%s %s\n", status, domain)
		fmt.Printf("   IP: %s | 延迟: %v | 跳数: %d\n", 
			analysis.IP, analysis.Latency, analysis.Hops)
		if analysis.Suspicious {
			fmt.Printf("   ⚠️  %s\n", analysis.Reason)
		}
		fmt.Println()
	}

	// 给出建议
	fmt.Println("💡 地理位置选择建议:")
	printGeographyAdvice(serverRegion)

	// 展示流量特征分析
	fmt.Println("\n🔍 流量特征分析:")
	analyzeTrafficPatterns(serverRegion)
}

func detectServerRegion() string {
	// 简化的地区检测（实际应用中可以使用 IP 地理位置 API）
	// 这里基于一些启发式规则
	
	// 检查时区
	cmd := exec.Command("timedatectl", "show", "-p", "Timezone", "--value")
	output, err := cmd.Output()
	if err == nil {
		timezone := strings.TrimSpace(string(output))
		switch {
		case strings.Contains(timezone, "Asia/Shanghai") || strings.Contains(timezone, "Asia/Beijing"):
			return "中国"
		case strings.Contains(timezone, "Asia/"):
			return "亚太"
		case strings.Contains(timezone, "Europe/"):
			return "欧洲"
		case strings.Contains(timezone, "America/"):
			return "北美"
		}
	}

	// 检查公网 IP（简化版）
	cmd = exec.Command("curl", "-s", "ifconfig.me")
	output, err = cmd.Output()
	if err == nil {
		ip := strings.TrimSpace(string(output))
		// 这里可以查询 IP 地理位置数据库
		fmt.Printf("检测到公网 IP: %s\n", ip)
	}

	return "未知地区"
}

func analyzeDomain(domain, serverRegion string) GeographyAnalysis {
	host := strings.Split(domain, ":")[0]
	
	// 解析 IP
	ips, err := net.LookupIP(host)
	var ip string
	if err != nil || len(ips) == 0 {
		ip = "解析失败"
	} else {
		ip = ips[0].String()
	}

	// 测试延迟
	latency := testLatency(domain)
	
	// 简单的跳数估计（基于延迟）
	hops := estimateHops(latency)
	
	// 分析是否可疑
	suspicious, reason := analyzeSuspiciousness(domain, serverRegion, latency)

	return GeographyAnalysis{
		Domain:     domain,
		IP:         ip,
		Latency:    latency,
		Hops:       hops,
		Region:     guessRegionFromDomain(domain),
		Suspicious: suspicious,
		Reason:     reason,
	}
}

func testLatency(domain string) time.Duration {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", domain, 5*time.Second)
	if err != nil {
		return time.Duration(-1)
	}
	defer conn.Close()
	return time.Since(start)
}

func estimateHops(latency time.Duration) int {
	if latency < 0 {
		return -1
	}
	// 粗略估计：每 10ms 约 1-2 跳
	return int(latency.Milliseconds() / 15)
}

func guessRegionFromDomain(domain string) string {
	host := strings.Split(domain, ":")[0]
	switch {
	case strings.Contains(host, "microsoft.com"):
		return "全球CDN"
	case strings.Contains(host, "apple.com"):
		return "全球CDN"
	case strings.Contains(host, "amazon.com"):
		return "主要北美"
	case strings.Contains(host, "github.com"):
		return "主要北美"
	case strings.Contains(host, "baidu.com"):
		return "中国"
	default:
		return "未知"
	}
}

func analyzeSuspiciousness(domain, serverRegion string, latency time.Duration) (bool, string) {
	// 延迟过高
	if latency > 500*time.Millisecond {
		return true, "延迟过高，可能引起注意"
	}
	
	// 地理位置不匹配的情况
	host := strings.Split(domain, ":")[0]
	if serverRegion == "中国" {
		if strings.Contains(host, "gov") {
			return true, "政府域名在中国可能敏感"
		}
		if latency > 300*time.Millisecond {
			return true, "海外域名延迟过高"
		}
	}
	
	return false, ""
}

func printGeographyAdvice(region string) {
	switch region {
	case "中国":
		fmt.Println("• 优先选择有中国 CDN 节点的国际大厂")
		fmt.Println("• 避免被墙的域名和政治敏感网站")
		fmt.Println("• 延迟控制在 200ms 以内")
		fmt.Println("• 推荐：微软、苹果等有本土化的服务")
		
	case "北美":
		fmt.Println("• 可以选择大部分美国本土网站")
		fmt.Println("• 延迟通常很低，选择面广")
		fmt.Println("• 避免明显的技术/代理相关域名")
		fmt.Println("• 推荐：AWS、GitHub、Netflix 等")
		
	case "欧洲":
		fmt.Println("• 选择欧洲本土或全球 CDN 网站")
		fmt.Println("• 注意 GDPR 合规的网站特征更自然")
		fmt.Println("• 避免美国政府相关域名")
		fmt.Println("• 推荐：Spotify、Booking 等欧洲公司")
		
	case "亚太":
		fmt.Println("• 选择亚太地区常用的国际网站")
		fmt.Println("• 考虑当地的网络基础设施")
		fmt.Println("• 日韩用户常访问的网站是好选择")
		fmt.Println("• 推荐：Sony、Samsung 等亚洲公司")
		
	default:
		fmt.Println("• 建议先确定服务器的具体地理位置")
		fmt.Println("• 选择全球 CDN 覆盖好的大厂域名")
		fmt.Println("• 测试延迟选择最优的域名")
	}
}

func analyzeTrafficPatterns(region string) {
	fmt.Println("正常用户访问模式分析:")
	fmt.Println()
	
	patterns := map[string][]string{
		"中国用户": {
			"• 经常访问：微软、苹果、GitHub（尽管慢）",
			"• 较少访问：Netflix、Facebook、Twitter",
			"• 时间模式：主要在北京时间工作时间",
			"• 特征：对海外网站延迟容忍度较高",
		},
		"美国用户": {
			"• 经常访问：本土网站延迟极低",
			"• 较多访问：社交媒体、流媒体服务",
			"• 时间模式：分布在各个时区",
			"• 特征：期望低延迟响应",
		},
		"欧洲用户": {
			"• 经常访问：注重隐私的服务",
			"• 较多访问：本地化服务和全球服务",
			"• 时间模式：集中在欧洲工作时间",
			"• 特征：对数据保护敏感",
		},
	}
	
	var userType string
	switch region {
	case "中国":
		userType = "中国用户"
	case "北美":
		userType = "美国用户"
	case "欧洲":
		userType = "欧洲用户"
	default:
		userType = "中国用户" // 默认显示
	}
	
	if patterns[userType] != nil {
		fmt.Printf("%s 的正常访问模式:\n", userType)
		for _, pattern := range patterns[userType] {
			fmt.Println(pattern)
		}
	}
	
	fmt.Println("\n⚠️  异常模式警告:")
	fmt.Println("• 中国服务器频繁访问南美网站 → 异常")
	fmt.Println("• 美国服务器只访问中国网站 → 可疑")
	fmt.Println("• 延迟模式与地理位置不符 → 需要注意")
	fmt.Println("• 访问时间与当地时区不符 → 可能被识别")
}
