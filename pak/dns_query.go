package pak

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: DNS 查询
 * @Project: pf_tools
 */

type DNSTool struct {
	Quiet bool
	Type  string
}

func (d *DNSTool) Run(host string) {
	view := ResultView{Quiet: d.Quiet}
	host = strings.TrimSpace(host)
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	if i := strings.Index(host, "/"); i >= 0 {
		host = host[:i]
	}
	if host == "" {
		view.ShowText("DNS", "请输入域名", "", false)
		return
	}

	typ := strings.ToUpper(strings.TrimSpace(d.Type))
	if typ == "" {
		typ = "A"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resolver := net.DefaultResolver

	rows := [][]string{{"类型", "值"}}
	var copyText string

	switch typ {
	case "A":
		ips, err := resolver.LookupIP(ctx, "ip4", host)
		if err != nil {
			view.ShowText("DNS", "查询失败: "+err.Error(), "", false)
			return
		}
		for _, ip := range ips {
			rows = append(rows, []string{"A", ip.String()})
			if copyText == "" {
				copyText = ip.String()
			}
		}
	case "AAAA":
		ips, err := resolver.LookupIP(ctx, "ip6", host)
		if err != nil {
			view.ShowText("DNS", "查询失败: "+err.Error(), "", false)
			return
		}
		for _, ip := range ips {
			rows = append(rows, []string{"AAAA", ip.String()})
			if copyText == "" {
				copyText = ip.String()
			}
		}
	case "CNAME":
		cname, err := resolver.LookupCNAME(ctx, host)
		if err != nil {
			view.ShowText("DNS", "查询失败: "+err.Error(), "", false)
			return
		}
		rows = append(rows, []string{"CNAME", strings.TrimSuffix(cname, ".")})
		copyText = strings.TrimSuffix(cname, ".")
	case "MX":
		mxs, err := resolver.LookupMX(ctx, host)
		if err != nil {
			view.ShowText("DNS", "查询失败: "+err.Error(), "", false)
			return
		}
		for _, mx := range mxs {
			val := fmt.Sprintf("%d %s", mx.Pref, strings.TrimSuffix(mx.Host, "."))
			rows = append(rows, []string{"MX", val})
			if copyText == "" {
				copyText = val
			}
		}
	case "TXT":
		txts, err := resolver.LookupTXT(ctx, host)
		if err != nil {
			view.ShowText("DNS", "查询失败: "+err.Error(), "", false)
			return
		}
		for _, t := range txts {
			rows = append(rows, []string{"TXT", truncate(t, 80)})
			if copyText == "" {
				copyText = t
			}
		}
	case "NS":
		nss, err := resolver.LookupNS(ctx, host)
		if err != nil {
			view.ShowText("DNS", "查询失败: "+err.Error(), "", false)
			return
		}
		for _, ns := range nss {
			val := strings.TrimSuffix(ns.Host, ".")
			rows = append(rows, []string{"NS", val})
			if copyText == "" {
				copyText = val
			}
		}
	case "ANY", "ALL":
		d.queryAll(view, resolver, ctx, host)
		return
	default:
		view.ShowText("DNS", "不支持的类型: "+typ+"（A/AAAA/CNAME/MX/TXT/NS/ALL）", "", false)
		return
	}

	if len(rows) == 1 {
		view.ShowText("DNS", "无记录", "", false)
		return
	}
	view.ShowTable("DNS "+host, rows, copyText)
}

func (d *DNSTool) queryAll(view ResultView, resolver *net.Resolver, ctx context.Context, host string) {
	rows := [][]string{{"类型", "值"}}
	var first string
	add := func(typ, val string) {
		rows = append(rows, []string{typ, truncate(val, 80)})
		if first == "" {
			first = val
		}
	}

	if ips, err := resolver.LookupIP(ctx, "ip4", host); err == nil {
		for _, ip := range ips {
			add("A", ip.String())
		}
	}
	if ips, err := resolver.LookupIP(ctx, "ip6", host); err == nil {
		for _, ip := range ips {
			add("AAAA", ip.String())
		}
	}
	if cname, err := resolver.LookupCNAME(ctx, host); err == nil && cname != "" {
		add("CNAME", strings.TrimSuffix(cname, "."))
	}
	if mxs, err := resolver.LookupMX(ctx, host); err == nil {
		for _, mx := range mxs {
			add("MX", fmt.Sprintf("%d %s", mx.Pref, strings.TrimSuffix(mx.Host, ".")))
		}
	}
	if txts, err := resolver.LookupTXT(ctx, host); err == nil {
		for _, t := range txts {
			add("TXT", t)
		}
	}
	if nss, err := resolver.LookupNS(ctx, host); err == nil {
		for _, ns := range nss {
			add("NS", strings.TrimSuffix(ns.Host, "."))
		}
	}
	if len(rows) == 1 {
		view.ShowText("DNS", "无记录", "", false)
		return
	}
	view.ShowTable("DNS "+host, rows, first)
}
