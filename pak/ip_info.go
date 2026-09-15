package pak

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: IP 查询
 * @Project: pf_tools
 */

type IPTool struct {
	Quiet bool
	Local bool
}

func (t *IPTool) Run(ip string) {
	view := ResultView{Quiet: t.Quiet}

	if t.Local {
		rows := [][]string{{"网卡", "IP"}}
		ifaces, err := net.Interfaces()
		if err != nil {
			view.ShowText("IP", err.Error(), "", false)
			return
		}
		for _, iface := range ifaces {
			addrs, _ := iface.Addrs()
			for _, a := range addrs {
				rows = append(rows, []string{iface.Name, a.String()})
			}
		}
		view.ShowTable("本地IP", rows, "")
		return
	}

	ip = strings.TrimSpace(ip)
	if ip == "" {
		pub, err := fetchPublicIP()
		if err != nil {
			view.ShowText("IP", "获取公网 IP 失败: "+err.Error(), "", false)
			return
		}
		view.ShowText("公网IP", pub, pub, true)
		return
	}

	parsed := net.ParseIP(ip)
	if parsed == nil {
		view.ShowText("IP", "IP 地址无效", "", false)
		return
	}
	info := fmt.Sprintf("地址: %s\n类型: %s\n回环: %v\n私有: %v",
		parsed.String(),
		ipFamily(parsed),
		parsed.IsLoopback(),
		parsed.IsPrivate(),
	)
	view.ShowText("IP信息", info, parsed.String(), true)
}

func ipFamily(ip net.IP) string {
	if ip.To4() != nil {
		return "IPv4"
	}
	return "IPv6"
}

func fetchPublicIP() (string, error) {
	client := &http.Client{Timeout: 6 * time.Second}
	endpoints := []string{
		"https://api.ipify.org?format=json",
		"https://ifconfig.me/ip",
	}
	for _, ep := range endpoints {
		req, _ := http.NewRequest(http.MethodGet, ep, nil)
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		_ = resp.Body.Close()
		if resp.StatusCode != 200 {
			continue
		}
		text := strings.TrimSpace(string(b))
		if strings.HasPrefix(text, "{") {
			var m map[string]string
			if json.Unmarshal(b, &m) == nil {
				if v := m["ip"]; v != "" {
					return v, nil
				}
			}
			continue
		}
		if net.ParseIP(text) != nil {
			return text, nil
		}
	}
	return "", fmt.Errorf("所有查询接口均失败")
}
