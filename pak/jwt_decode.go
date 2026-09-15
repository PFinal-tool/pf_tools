package pak

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: JWT 解码（不校验签名）
 * @Project: pf_tools
 */

type JWTTool struct {
	Quiet bool
	Exp   bool
}

func (j *JWTTool) Run(token string) {
	view := ResultView{Quiet: j.Quiet}
	token = strings.TrimSpace(token)
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		view.ShowText("JWT", "无效 JWT：分段不足", "", false)
		return
	}

	headerJSON, err := decodeJWTPart(parts[0])
	if err != nil {
		view.ShowText("JWT", "header 解码失败: "+err.Error(), "", false)
		return
	}
	payloadJSON, err := decodeJWTPart(parts[1])
	if err != nil {
		view.ShowText("JWT", "payload 解码失败: "+err.Error(), "", false)
		return
	}

	headerPretty := prettyJSON(headerJSON)
	payloadPretty := prettyJSON(payloadJSON)

	if j.Exp {
		msg := expMessage(payloadJSON)
		view.ShowText("JWT过期检查", msg, msg, !strings.HasPrefix(msg, "已过期"))
		return
	}

	text := "Header:\n" + headerPretty + "\n\nPayload:\n" + payloadPretty
	view.ShowText("JWT解码", text, payloadPretty, true)
}

func decodeJWTPart(part string) ([]byte, error) {
	// JWT 使用 Base64URL，无 padding
	b, err := base64.RawURLEncoding.DecodeString(part)
	if err == nil {
		return b, nil
	}
	b, err = base64.URLEncoding.DecodeString(part)
	if err == nil {
		return b, nil
	}
	return base64.RawStdEncoding.DecodeString(part)
}

func prettyJSON(raw []byte) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		return string(raw)
	}
	return buf.String()
}

func expMessage(payload []byte) string {
	var m map[string]any
	if err := json.Unmarshal(payload, &m); err != nil {
		return "payload 不是 JSON 对象"
	}
	expRaw, ok := m["exp"]
	if !ok {
		return "payload 中无 exp 字段"
	}
	var exp int64
	switch v := expRaw.(type) {
	case float64:
		exp = int64(v)
	case json.Number:
		n, _ := v.Int64()
		exp = n
	default:
		return fmt.Sprintf("exp 字段类型不支持: %T", expRaw)
	}
	expireAt := time.Unix(exp, 0)
	now := time.Now()
	remain := expireAt.Sub(now)
	if remain <= 0 {
		return fmt.Sprintf("已过期\n过期时间: %s\n已过期: %s",
			expireAt.Format("2006-01-02 15:04:05"), (-remain).String())
	}
	return fmt.Sprintf("未过期\n过期时间: %s\n剩余: %s",
		expireAt.Format("2006-01-02 15:04:05"), remain.Round(time.Second).String())
}
