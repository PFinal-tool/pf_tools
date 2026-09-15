package pak

import (
	"strings"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: Cookie / Header 解析
 * @Project: pf_tools
 */

type CookieTool struct {
	Quiet bool
}

type HeaderTool struct {
	Quiet bool
}

func (c *CookieTool) Run(input string) {
	view := ResultView{Quiet: c.Quiet}
	input = strings.TrimSpace(input)
	if input == "" {
		view.ShowText("Cookie", "请输入 Cookie 字符串", "", false)
		return
	}
	// 支持 "Cookie: a=1; b=2" 或裸串
	if i := strings.Index(strings.ToLower(input), "cookie:"); i >= 0 {
		input = strings.TrimSpace(input[i+7:])
	}
	rows := [][]string{{"Key", "Value"}}
	parts := strings.Split(input, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		k := strings.TrimSpace(kv[0])
		v := ""
		if len(kv) == 2 {
			v = strings.TrimSpace(kv[1])
		}
		rows = append(rows, []string{k, truncate(v, 80)})
	}
	if len(rows) == 1 {
		view.ShowText("Cookie", "未解析到键值对", "", false)
		return
	}
	view.ShowTable("Cookie", rows, rows[1][1])
}

func (h *HeaderTool) Run(input string) {
	view := ResultView{Quiet: h.Quiet}
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.TrimSpace(input)
	if input == "" {
		view.ShowText("Header", "请输入 HTTP Header 文本", "", false)
		return
	}

	rows := [][]string{{"Header", "Value"}}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if isRequestLine(line) {
			rows = append(rows, []string{"(request)", line})
			continue
		}
		idx := strings.Index(line, ":")
		if idx <= 0 {
			rows = append(rows, []string{"(raw)", truncate(line, 80)})
			continue
		}
		k := strings.TrimSpace(line[:idx])
		v := strings.TrimSpace(line[idx+1:])
		rows = append(rows, []string{k, truncate(v, 80)})
	}
	if len(rows) == 1 {
		view.ShowText("Header", "未解析到 Header", "", false)
		return
	}
	view.ShowTable("HTTP Header", rows, "")
}

func isRequestLine(line string) bool {
	methods := []string{"GET ", "POST ", "PUT ", "DELETE ", "PATCH ", "HEAD ", "OPTIONS "}
	upper := strings.ToUpper(line)
	for _, m := range methods {
		if strings.HasPrefix(upper, m) {
			return true
		}
	}
	return strings.HasPrefix(upper, "HTTP/")
}
