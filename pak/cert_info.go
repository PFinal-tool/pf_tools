package pak

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: TLS 证书信息查询
 * @Project: pf_tools
 */

type CertTool struct {
	Quiet bool
	Days  bool
}

func (c *CertTool) Run(hostPort string) {
	view := ResultView{Quiet: c.Quiet}
	hostPort = strings.TrimSpace(hostPort)
	if hostPort == "" {
		view.ShowText("证书", "请输入域名，如 example.com", "", false)
		return
	}
	if !strings.Contains(hostPort, ":") {
		hostPort += ":443"
	}

	host, _, err := net.SplitHostPort(hostPort)
	if err != nil {
		view.ShowText("证书", "地址无效: "+err.Error(), "", false)
		return
	}

	dialer := &net.Dialer{Timeout: 8 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", hostPort, &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: true, // 仅读取证书信息，不因过期/自签中断
	})
	if err != nil {
		view.ShowText("证书", "连接失败: "+err.Error(), "", false)
		return
	}
	defer conn.Close()

	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		view.ShowText("证书", "未获取到证书", "", false)
		return
	}
	cert := state.PeerCertificates[0]
	now := time.Now()
	remain := cert.NotAfter.Sub(now)
	remainDays := int(remain.Hours() / 24)

	if c.Days {
		msg := fmt.Sprintf("%d", remainDays)
		if remain < 0 {
			msg = fmt.Sprintf("%d", remainDays) // 负数表示已过期天数
		}
		text := fmt.Sprintf("剩余天数: %d\n过期时间: %s", remainDays, cert.NotAfter.Local().Format("2006-01-02 15:04:05"))
		view.ShowText("证书有效期", text, msg, remain >= 0)
		return
	}

	dns := strings.Join(cert.DNSNames, ", ")
	if dns == "" {
		dns = "-"
	}
	issuer := cert.Issuer.CommonName
	if issuer == "" && len(cert.Issuer.Organization) > 0 {
		issuer = cert.Issuer.Organization[0]
	}
	subject := cert.Subject.CommonName
	if subject == "" {
		subject = host
	}

	status := "有效"
	if now.Before(cert.NotBefore) {
		status = "尚未生效"
	} else if now.After(cert.NotAfter) {
		status = "已过期"
	}

	rows := [][]string{
		{"字段", "值"},
		{"地址", hostPort},
		{"主题", subject},
		{"颁发者", issuer},
		{"生效", cert.NotBefore.Local().Format("2006-01-02 15:04:05")},
		{"过期", cert.NotAfter.Local().Format("2006-01-02 15:04:05")},
		{"剩余天数", fmt.Sprintf("%d", remainDays)},
		{"状态", status},
		{"SAN", truncate(dns, 80)},
		{"算法", cert.PublicKeyAlgorithm.String()},
	}
	view.ShowTable("TLS证书", rows, fmt.Sprintf("%d", remainDays))
}
