package pak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: 轻量 HTTP 探活
 * @Project: pf_tools
 */

type HTTPTool struct {
	Quiet   bool
	Method  string
	Body    string
	Header  []string
	JSONFmt bool
	ShowTime bool
}

func (h *HTTPTool) Run(rawURL string) {
	view := ResultView{Quiet: h.Quiet}
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		view.ShowText("HTTP", "请输入 URL", "", false)
		return
	}
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}

	method := strings.ToUpper(strings.TrimSpace(h.Method))
	if method == "" {
		method = http.MethodGet
	}

	var bodyReader io.Reader
	if h.Body != "" {
		bodyReader = strings.NewReader(h.Body)
	}
	req, err := http.NewRequest(method, rawURL, bodyReader)
	if err != nil {
		view.ShowText("HTTP", err.Error(), "", false)
		return
	}
	req.Header.Set("User-Agent", "pf_tools/1.4")
	hasCT := false
	for _, hv := range h.Header {
		k, v, ok := strings.Cut(hv, ":")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		req.Header.Set(k, v)
		if strings.EqualFold(k, "Content-Type") {
			hasCT = true
		}
	}
	if h.Body != "" && !hasCT {
		if json.Valid([]byte(h.Body)) {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	client := &http.Client{Timeout: 15 * time.Second}
	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		view.ShowText("HTTP", "请求失败: "+err.Error(), "", false)
		return
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20)) // 2MB
	if err != nil {
		view.ShowText("HTTP", "读取响应失败: "+err.Error(), "", false)
		return
	}

	bodyText := string(raw)
	ct := resp.Header.Get("Content-Type")
	if h.JSONFmt || strings.Contains(ct, "json") {
		var pretty bytes.Buffer
		if json.Indent(&pretty, raw, "", "  ") == nil {
			bodyText = pretty.String()
		}
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s %s\n", method, rawURL))
	b.WriteString(fmt.Sprintf("状态: %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode)))
	if h.ShowTime {
		b.WriteString(fmt.Sprintf("耗时: %s\n", elapsed.Round(time.Millisecond)))
	}
	b.WriteString("响应头:\n")
	for k, vs := range resp.Header {
		b.WriteString(fmt.Sprintf("  %s: %s\n", k, strings.Join(vs, ", ")))
	}
	b.WriteString("响应体:\n")
	b.WriteString(truncate(bodyText, 3000))

	text := strings.TrimRight(b.String(), "\n")
	copyText := bodyText
	if h.Quiet {
		if h.ShowTime {
			view.Print(fmt.Sprintf("%d\t%s\t%s", resp.StatusCode, elapsed.Round(time.Millisecond), truncate(oneLine(bodyText), 200)))
			return
		}
		view.Print(fmt.Sprintf("%d\n%s", resp.StatusCode, bodyText))
		return
	}
	ok := resp.StatusCode >= 200 && resp.StatusCode < 400
	view.ShowText("HTTP", text, copyText, ok)
}

func oneLine(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	return strings.TrimSpace(s)
}
