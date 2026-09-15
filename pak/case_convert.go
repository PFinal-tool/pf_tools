package pak

import (
	"regexp"
	"strings"
	"unicode"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: 命名风格转换
 * @Project: pf_tools
 */

type CaseTool struct {
	Quiet bool
	To    string // camel / snake / kebab / pascal / export
}

func (c *CaseTool) Run(input string) {
	view := ResultView{Quiet: c.Quiet}
	input = strings.TrimSpace(input)
	if input == "" {
		view.ShowText("Case", "请输入标识符", "", false)
		return
	}
	words := splitIdent(input)
	if len(words) == 0 {
		view.ShowText("Case", "无法解析标识符", "", false)
		return
	}

	to := strings.ToLower(c.To)
	if to == "" {
		// 未指定则输出多种格式
		rows := [][]string{
			{"风格", "结果"},
			{"camel", toCamel(words, false)},
			{"pascal", toCamel(words, true)},
			{"snake", toSnake(words)},
			{"kebab", toKebab(words)},
			{"export", strings.ToUpper(toSnake(words))},
		}
		view.ShowTable("命名转换", rows, toCamel(words, false))
		return
	}

	var out string
	switch to {
	case "camel":
		out = toCamel(words, false)
	case "pascal":
		out = toCamel(words, true)
	case "snake":
		out = toSnake(words)
	case "kebab", "dash":
		out = toKebab(words)
	case "export", "upper", "screaming":
		out = strings.ToUpper(toSnake(words))
	default:
		view.ShowText("Case", "不支持的目标: "+c.To+"（camel/pascal/snake/kebab/export）", "", false)
		return
	}
	view.ShowText("命名转换", out, out, true)
}

func splitIdent(s string) []string {
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	if strings.Contains(s, "_") {
		parts := strings.Split(s, "_")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				out = append(out, strings.ToLower(p))
			}
		}
		return out
	}
	// camel / Pascal: UserID -> user, id
	re := regexp.MustCompile(`([A-Z]+[a-z0-9]*|[a-z0-9]+)`)
	parts := re.FindAllString(s, -1)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.ToLower(p))
	}
	return out
}

func toCamel(words []string, pascal bool) string {
	var b strings.Builder
	for i, w := range words {
		if i == 0 && !pascal {
			b.WriteString(w)
			continue
		}
		b.WriteString(upperFirst(w))
	}
	return b.String()
}

func toSnake(words []string) string {
	return strings.Join(words, "_")
}

func toKebab(words []string) string {
	return strings.Join(words, "-")
}

func upperFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
