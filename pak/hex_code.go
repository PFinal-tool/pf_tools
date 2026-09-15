package pak

import (
	"encoding/hex"
	"strings"
	"unicode/utf8"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: 十六进制编解码
 * @Project: pf_tools
 */

type HexTool struct {
	Quiet  bool
	Decode bool
}

func (h *HexTool) Run(input string) {
	view := ResultView{Quiet: h.Quiet}
	input = strings.TrimSpace(input)
	if input == "" {
		view.ShowText("Hex", "请输入字符串或十六进制", "", false)
		return
	}
	if h.Decode {
		clean := strings.ReplaceAll(input, " ", "")
		clean = strings.TrimPrefix(strings.ToLower(clean), "0x")
		b, err := hex.DecodeString(clean)
		if err != nil {
			view.ShowText("Hex", "解码失败: "+err.Error(), "", false)
			return
		}
		if !utf8.Valid(b) {
			out := hex.EncodeToString(b)
			view.ShowText("Hex", "非 UTF-8 文本，原始字节长度 "+itoa(len(b)), out, true)
			return
		}
		text := string(b)
		view.ShowText("Hex解码", text, text, true)
		return
	}
	out := hex.EncodeToString([]byte(input))
	view.ShowText("Hex编码", out, out, true)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
