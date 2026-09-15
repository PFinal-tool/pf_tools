package pak

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: CIDR / IP 归属判断
 * @Project: pf_tools
 */

type CIDRTool struct {
	Quiet bool
	List  bool
	Count int
}

func (c *CIDRTool) Run(args []string) {
	view := ResultView{Quiet: c.Quiet}
	if len(args) == 0 {
		view.ShowText("CIDR", "用法:\n  pf_cidr 10.0.0.8 10.0.0.0/24\n  pf_cidr 10.0.0.0/24 --list -n 5", "", false)
		return
	}

	// 单参数且含 / ：展示网段信息 / 列表
	if len(args) == 1 && strings.Contains(args[0], "/") {
		_, network, err := net.ParseCIDR(args[0])
		if err != nil {
			view.ShowText("CIDR", "CIDR 无效: "+err.Error(), "", false)
			return
		}
		if c.List {
			c.listIPs(view, network)
			return
		}
		ones, bits := network.Mask.Size()
		total := uint64(1) << uint(bits-ones)
		if bits-ones >= 64 {
			total = 0
		}
		info := fmt.Sprintf("网段: %s\n网络: %s\n掩码: %s\n前缀: /%d\n大约主机数: %s",
			args[0],
			network.IP.String(),
			net.IP(network.Mask).String(),
			ones,
			hostCountText(ones, bits, total),
		)
		view.ShowText("CIDR", info, network.String(), true)
		return
	}

	// IP + CIDR 判断
	if len(args) >= 2 {
		ip := net.ParseIP(args[0])
		if ip == nil {
			view.ShowText("CIDR", "IP 无效", "", false)
			return
		}
		_, network, err := net.ParseCIDR(args[1])
		if err != nil {
			view.ShowText("CIDR", "CIDR 无效: "+err.Error(), "", false)
			return
		}
		ok := network.Contains(ip)
		msg := fmt.Sprintf("%s 不在 %s", ip.String(), network.String())
		if ok {
			msg = fmt.Sprintf("%s 属于 %s", ip.String(), network.String())
		}
		if c.Quiet {
			if ok {
				view.Print("true")
			} else {
				view.Print("false")
			}
			return
		}
		view.ShowText("CIDR", msg, strconv.FormatBool(ok), ok)
		return
	}

	view.ShowText("CIDR", "参数不足", "", false)
}

func (c *CIDRTool) listIPs(view ResultView, network *net.IPNet) {
	n := c.Count
	if n <= 0 {
		n = 5
	}
	if n > 256 {
		n = 256
	}
	ips := make([]string, 0, n)
	for ip := network.IP.Mask(network.Mask); network.Contains(ip); incIP(ip) {
		// 跳过网络地址本身（常见习惯，仍输出前 n 个可用展示）
		cp := append(net.IP(nil), ip...)
		ips = append(ips, cp.String())
		if len(ips) >= n {
			break
		}
	}
	text := strings.Join(ips, "\n")
	view.ShowText("CIDR列表", text, text, true)
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func hostCountText(ones, bits int, total uint64) string {
	if bits-ones >= 32 {
		return fmt.Sprintf("2^%d", bits-ones)
	}
	if total == 0 {
		return "-"
	}
	// 粗略可用主机（减网络/广播，仅 IPv4 且前缀<31）
	if bits == 32 && ones < 31 && total >= 2 {
		return fmt.Sprintf("%d (约可用 %d)", total, total-2)
	}
	return fmt.Sprintf("%d", total)
}
