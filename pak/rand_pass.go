package pak

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: 随机串 / 密码生成
 * @Project: pf_tools
 */

const (
	charsetAlphaNum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetDigit    = "0123456789"
	charsetHex      = "0123456789abcdef"
	charsetBase58   = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	charsetSpecial  = "!@#$%^&*()-_=+[]{}?,."
	charsetPass     = charsetAlphaNum + charsetSpecial
)

type RandTool struct {
	Quiet bool
	Len   int
	Type  string // alphanum / digit / hex / base58
	Count int
}

type PassTool struct {
	Quiet     bool
	Len       int
	Special   bool
	Count     int
}

func (r *RandTool) Run() {
	view := ResultView{Quiet: r.Quiet}
	n := r.Len
	if n <= 0 {
		n = 32
	}
	if n > 4096 {
		view.ShowText("随机串", "长度过大", "", false)
		return
	}
	count := r.Count
	if count <= 0 {
		count = 1
	}
	if count > 100 {
		view.ShowText("随机串", "一次最多生成 100 个", "", false)
		return
	}

	charset, err := charsetByType(r.Type)
	if err != nil {
		view.ShowText("随机串", err.Error(), "", false)
		return
	}

	list := make([]string, 0, count)
	for i := 0; i < count; i++ {
		s, err := randomFromCharset(n, charset)
		if err != nil {
			view.ShowText("随机串", err.Error(), "", false)
			return
		}
		list = append(list, s)
	}
	text := strings.Join(list, "\n")
	view.ShowText("随机串", text, text, true)
}

func (p *PassTool) Run() {
	view := ResultView{Quiet: p.Quiet}
	n := p.Len
	if n <= 0 {
		n = 16
	}
	if n < 4 {
		view.ShowText("密码", "长度至少为 4", "", false)
		return
	}
	if n > 256 {
		view.ShowText("密码", "长度过大", "", false)
		return
	}
	count := p.Count
	if count <= 0 {
		count = 1
	}
	if count > 100 {
		view.ShowText("密码", "一次最多生成 100 个", "", false)
		return
	}

	list := make([]string, 0, count)
	for i := 0; i < count; i++ {
		s, err := randomPassword(n, p.Special)
		if err != nil {
			view.ShowText("密码", err.Error(), "", false)
			return
		}
		list = append(list, s)
	}
	text := strings.Join(list, "\n")
	view.ShowText("密码", text, text, true)
}

func charsetByType(t string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "", "alphanum", "an":
		return charsetAlphaNum, nil
	case "digit", "num", "number":
		return charsetDigit, nil
	case "hex":
		return charsetHex, nil
	case "base58":
		return charsetBase58, nil
	default:
		return "", fmt.Errorf("不支持的类型: %s（可用 alphanum/digit/hex/base58）", t)
	}
}

func randomFromCharset(n int, charset string) (string, error) {
	b := make([]byte, n)
	max := big.NewInt(int64(len(charset)))
	for i := 0; i < n; i++ {
		v, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = charset[v.Int64()]
	}
	return string(b), nil
}

func randomPassword(n int, special bool) (string, error) {
	// 保证至少包含大小写和数字；开启 special 时再保证一个特殊字符
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := charsetDigit
	specials := charsetSpecial

	required := []string{lower, upper, digits}
	if special {
		required = append(required, specials)
	}
	if n < len(required) {
		return "", fmt.Errorf("长度不足以覆盖字符类别")
	}

	out := make([]byte, 0, n)
	for _, set := range required {
		ch, err := randomFromCharset(1, set)
		if err != nil {
			return "", err
		}
		out = append(out, ch[0])
	}

	pool := charsetAlphaNum
	if special {
		pool = charsetPass
	}
	rest, err := randomFromCharset(n-len(required), pool)
	if err != nil {
		return "", err
	}
	out = append(out, rest...)

	// Fisher–Yates 洗牌
	for i := len(out) - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		j := int(jBig.Int64())
		out[i], out[j] = out[j], out[i]
	}
	return string(out), nil
}
