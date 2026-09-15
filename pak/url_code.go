package pak

import (
	"fmt"
	"net/url"
	"strings"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: URL 编解码 / 拆解
 * @Project: pf_tools
 */

type URLTool struct {
	Quiet  bool
	Decode bool
	Parse  bool
}

func (u *URLTool) Run(input string) {
	view := ResultView{Quiet: u.Quiet}
	if input == "" {
		view.ShowText("URL", "请输入 URL 或待编码字符串", "", false)
		return
	}

	if u.Parse {
		u.showParse(input)
		return
	}

	if u.Decode {
		out, err := url.QueryUnescape(input)
		if err != nil {
			view.ShowText("URL", "解码失败: "+err.Error(), "", false)
			return
		}
		view.ShowText("URL解码", out, out, true)
		return
	}

	out := url.QueryEscape(input)
	view.ShowText("URL编码", out, out, true)
}

func (u *URLTool) showParse(input string) {
	view := ResultView{Quiet: u.Quiet}
	parsed, err := url.Parse(input)
	if err != nil {
		view.ShowText("URL", "解析失败: "+err.Error(), "", false)
		return
	}
	rows := [][]string{
		{"字段", "值"},
		{"scheme", parsed.Scheme},
		{"host", parsed.Host},
		{"path", parsed.Path},
		{"query", parsed.RawQuery},
		{"fragment", parsed.Fragment},
	}
	if parsed.User != nil {
		rows = append(rows, []string{"user", parsed.User.String()})
	}
	q := parsed.Query()
	for k, vs := range q {
		rows = append(rows, []string{"query."+k, strings.Join(vs, ", ")})
	}

	if u.Quiet {
		var b strings.Builder
		for i, row := range rows {
			if i == 0 {
				continue
			}
			fmt.Fprintf(&b, "%s\t%s\n", row[0], row[1])
		}
		view.Print(strings.TrimRight(b.String(), "\n"))
		return
	}
	view.ShowTable("URL解析", rows, parsed.String())
}
