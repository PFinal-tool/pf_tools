package pak

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/atotto/clipboard"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: 剪贴板流水线
 * @Project: pf_tools
 */

type ClipTool struct {
	Quiet bool
	Write bool
}

func (c *ClipTool) Run(action string) {
	view := ResultView{Quiet: c.Quiet}
	raw, err := clipboard.ReadAll()
	if err != nil {
		view.ShowText("Clip", "读取剪贴板失败: "+err.Error(), "", false)
		return
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		view.ShowText("Clip", "剪贴板为空", "", false)
		return
	}

	action = strings.ToLower(strings.TrimSpace(action))
	var out string
	switch action {
	case "", "show", "print":
		view.ShowText("剪贴板", truncate(raw, 2000), raw, true)
		return
	case "md5":
		out = fmt.Sprintf("%x", md5.Sum([]byte(raw)))
	case "b64", "base64":
		out = base64.StdEncoding.EncodeToString([]byte(raw))
	case "b64d", "base64d":
		b, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			view.ShowText("Clip", "Base64 解码失败", "", false)
			return
		}
		out = string(b)
	case "url":
		out = url.QueryEscape(raw)
	case "urld":
		v, err := url.QueryUnescape(raw)
		if err != nil {
			view.ShowText("Clip", "URL 解码失败", "", false)
			return
		}
		out = v
	case "jwt":
		parts := strings.Split(raw, ".")
		if len(parts) < 2 {
			view.ShowText("Clip", "剪贴板内容不是有效 JWT", "", false)
			return
		}
		payload, err := decodeJWTPart(parts[1])
		if err != nil {
			view.ShowText("Clip", "JWT payload 解码失败", "", false)
			return
		}
		out = prettyJSON(payload)
	case "json":
		var v any
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			view.ShowText("Clip", "JSON 无效: "+err.Error(), "", false)
			return
		}
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			view.ShowText("Clip", err.Error(), "", false)
			return
		}
		out = string(b)
	default:
		view.ShowText("Clip", "未知动作: "+action+"\n可用: show/md5/b64/b64d/url/urld/jwt/json", "", false)
		return
	}

	if c.Write {
		_ = clipboard.WriteAll(out)
	}
	view.ShowText("Clip:"+action, out, out, true)
}
