package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"sort"
	"time"
)

// DomainResult 域名测试结果
type DomainResult struct {
	Domain    string
	Latency   time.Duration
	TLSWorks  bool
	Error     error
	TLSVersion uint16
}

// 推荐的域名列表
var recommendedDomains = []string{
	"www.microsoft.com:443",
	"www.apple.com:443",
	"www.cloudflare.com:443",
	"www.amazon.com:443",
	"github.com:443",
	"stackoverflow.com:443",
	"www.docker.com:443",
	"aws.amazon.com:443",
	"cloud.google.com:443",
	"azure.microsoft.com:443",
	"www.google.com:443",
	"www.youtube.com:443",
	"www.facebook.com:443",
	"www.twitter.com:443",
	"www.instagram.com:443",
	"www.linkedin.com:443",
	"www.netflix.com:443",
	"www.spotify.com:443",
}

func main() {
	fmt.Println("🔍 Reality 域名测试工具")
	fmt.Println("正在测试推荐域名的连接性能...")
	fmt.Println()

	results := make([]DomainResult, 0, len(recommendedDomains))

	for _, domain := range recommendedDomains {
		result := testDomain(domain)
		results = append(results, result)
		
		status := "❌"
		if result.TLSWorks {
			status = "✅"
		}
		
		fmt.Printf("%s %s - 延迟: %v\n", status, result.Domain, result.Latency)
		if result.Error != nil {
			fmt.Printf("   错误: %v\n", result.Error)
		}
	}

	fmt.Println("\n📊 推荐排序（按延迟排序）:")
	
	// 过滤出可用的域名并排序
	var workingDomains []DomainResult
	for _, result := range results {
		if result.TLSWorks {
			workingDomains = append(workingDomains, result)
		}
	}

	sort.Slice(workingDomains, func(i, j int) bool {
		return workingDomains[i].Latency < workingDomains[j].Latency
	})

	fmt.Println("\n🏆 最佳选择（前5名）:")
	for i, result := range workingDomains {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s (延迟: %v, TLS版本: %s)\n", 
			i+1, result.Domain, result.Latency, getTLSVersionString(result.TLSVersion))
	}

	if len(workingDomains) > 0 {
		best := workingDomains[0]
		fmt.Printf("\n🎯 推荐配置:\n")
		fmt.Printf(`"dest": "%s",\n`, best.Domain)
		fmt.Printf(`"server_names": ["%s"],\n`, getHostFromDomain(best.Domain))
	}

	fmt.Println("\n💡 选择建议:")
	fmt.Println("1. 优先选择延迟最低的域名")
	fmt.Println("2. 建议选择大厂域名（微软、苹果、亚马逊等）")
	fmt.Println("3. 避免选择可能被封锁的域名")
	fmt.Println("4. 定期测试和更换域名")
}

func testDomain(domain string) DomainResult {
	start := time.Now()
	
	conn, err := net.DialTimeout("tcp", domain, 5*time.Second)
	if err != nil {
		return DomainResult{
			Domain:   domain,
			Latency:  time.Since(start),
			TLSWorks: false,
			Error:    err,
		}
	}
	defer conn.Close()

	// 测试 TLS 连接
	tlsConn := tls.Client(conn, &tls.Config{
		ServerName:         getHostFromDomain(domain),
		InsecureSkipVerify: false,
	})

	err = tlsConn.Handshake()
	latency := time.Since(start)
	
	if err != nil {
		return DomainResult{
			Domain:   domain,
			Latency:  latency,
			TLSWorks: false,
			Error:    err,
		}
	}

	tlsVersion := tlsConn.ConnectionState().Version
	tlsConn.Close()

	return DomainResult{
		Domain:     domain,
		Latency:    latency,
		TLSWorks:   true,
		Error:      nil,
		TLSVersion: tlsVersion,
	}
}

func getHostFromDomain(domain string) string {
	host, _, err := net.SplitHostPort(domain)
	if err != nil {
		return domain
	}
	return host
}

func getTLSVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown"
	}
}
