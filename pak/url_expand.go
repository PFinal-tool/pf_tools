package pak

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: 短链 / 重定向追踪
 * @Project: pf_tools
 */

type ExpandTool struct {
	Quiet bool
	Max   int
}

func (e *ExpandTool) Run(rawURL string) {
	view := ResultView{Quiet: e.Quiet}
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		view.ShowText("Expand", "请输入 URL", "", false)
		return
	}
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		view.ShowText("Expand", "URL 无效", "", false)
		return
	}

	max := e.Max
	if max <= 0 {
		max = 10
	}
	if max > 20 {
		max = 20
	}

	client := &http.Client{
		Timeout: 12 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	current := rawURL
	rows := [][]string{{"#", "状态", "URL"}}
	finalURL := rawURL
	for i := 0; i < max; i++ {
		req, err := http.NewRequest(http.MethodGet, current, nil)
		if err != nil {
			view.ShowText("Expand", err.Error(), "", false)
			return
		}
		req.Header.Set("User-Agent", "pf_tools/1.3")
		resp, err := client.Do(req)
		if err != nil {
			view.ShowText("Expand", "请求失败: "+err.Error(), "", false)
			return
		}
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
		_ = resp.Body.Close()

		rows = append(rows, []string{
			fmt.Sprintf("%d", i+1),
			fmt.Sprintf("%d", resp.StatusCode),
			truncate(current, 70),
		})
		finalURL = current

		if resp.StatusCode < 300 || resp.StatusCode >= 400 {
			break
		}
		loc := resp.Header.Get("Location")
		if loc == "" {
			break
		}
		next, err := resp.Request.URL.Parse(loc)
		if err != nil {
			break
		}
		current = next.String()
	}

	if e.Quiet {
		view.Print(finalURL)
		return
	}
	view.ShowTable("短链展开", rows, finalURL)
}
