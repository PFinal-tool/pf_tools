package pak

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: Unicode / \\uXXXX 处理
 * @Project: pf_tools
 */

type UnicodeTool struct {
	Quiet  bool
	Decode bool
}

func (u *UnicodeTool) Run(input string) {
	view := ResultView{Quiet: u.Quiet}
	input = strings.TrimSpace(input)
	if input == "" {
		view.ShowText("Unicode", "请输入文本或码点", "", false)
		return
	}

	if u.Decode || looksLikeUnicodeEscapes(input) {
		out, err := decodeUnicodeInput(input)
		if err != nil {
			view.ShowText("Unicode", err.Error(), "", false)
			return
		}
		view.ShowText("Unicode解码", out, out, true)
		return
	}

	var b strings.Builder
	var escapes []string
	for _, r := range input {
		escapes = append(escapes, fmt.Sprintf("U+%04X", r))
		if r <= 0xFFFF {
			b.WriteString(fmt.Sprintf("\\u%04x", r))
		} else {
			b.WriteString(fmt.Sprintf("\\U%08x", r))
		}
	}
	text := fmt.Sprintf("字符数: %d\n码点: %s\n转义: %s",
		utf8.RuneCountInString(input),
		strings.Join(escapes, " "),
		b.String(),
	)
	view.ShowText("Unicode", text, b.String(), true)
}

func looksLikeUnicodeEscapes(s string) bool {
	return strings.Contains(s, `\u`) || strings.Contains(s, `\U`) || strings.Contains(strings.ToUpper(s), "U+")
}

func decodeUnicodeInput(s string) (string, error) {
	// U+4F60 U+597D
	if strings.Contains(strings.ToUpper(s), "U+") {
		parts := strings.Fields(s)
		var runes []rune
		for _, p := range parts {
			p = strings.TrimSpace(p)
			p = strings.TrimPrefix(strings.ToUpper(p), "U+")
			v, err := strconv.ParseInt(p, 16, 32)
			if err != nil {
				return "", fmt.Errorf("码点无效: %s", p)
			}
			runes = append(runes, rune(v))
		}
		return string(runes), nil
	}

	// 处理 \uXXXX / \UXXXXXXXX
	var out strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '\\' && i+1 < len(s) {
			if s[i+1] == 'u' && i+6 <= len(s) {
				v, err := strconv.ParseInt(s[i+2:i+6], 16, 32)
				if err != nil {
					return "", err
				}
				out.WriteRune(rune(v))
				i += 6
				continue
			}
			if s[i+1] == 'U' && i+10 <= len(s) {
				v, err := strconv.ParseInt(s[i+2:i+10], 16, 32)
				if err != nil {
					return "", err
				}
				out.WriteRune(rune(v))
				i += 10
				continue
			}
		}
		out.WriteByte(s[i])
		i++
	}
	return out.String(), nil
}
