package pak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: JSON 格式化 / 压缩 / 校验 / 路径取值
 * @Project: pf_tools
 */

type JSONTool struct {
	Quiet  bool
	Compact bool
	Validate bool
	Path   string
}

func (j *JSONTool) Run(input string, fromFile bool) {
	view := ResultView{Quiet: j.Quiet}
	raw := input
	if fromFile {
		data, err := os.ReadFile(input)
		if err != nil {
			view.ShowText("JSON", "读取文件失败: "+err.Error(), "", false)
			return
		}
		raw = string(data)
	}

	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		view.ShowText("JSON", "无效 JSON: "+err.Error(), "", false)
		return
	}

	if j.Validate {
		msg := "JSON 有效"
		if j.Quiet {
			view.Print("ok")
			return
		}
		view.ShowText("JSON", msg, "ok", true)
		return
	}

	if j.Path != "" {
		got, err := jsonPath(v, j.Path)
		if err != nil {
			view.ShowText("JSON", err.Error(), "", false)
			return
		}
		out, err := json.MarshalIndent(got, "", "  ")
		if err != nil {
			view.ShowText("JSON", err.Error(), "", false)
			return
		}
		// 纯标量更友好
		if s, ok := got.(string); ok {
			view.ShowText("JSON", s, s, true)
			return
		}
		text := string(out)
		view.ShowText("JSON", text, text, true)
		return
	}

	var out []byte
	var err error
	if j.Compact {
		out, err = json.Marshal(v)
	} else {
		out, err = json.MarshalIndent(v, "", "  ")
	}
	if err != nil {
		view.ShowText("JSON", err.Error(), "", false)
		return
	}
	text := string(out)
	view.ShowText("JSON", text, text, true)
}

func jsonPath(v any, path string) (any, error) {
	path = strings.TrimPrefix(path, ".")
	if path == "" {
		return v, nil
	}
	parts := strings.Split(path, ".")
	cur := v
	for _, p := range parts {
		switch node := cur.(type) {
		case map[string]any:
			next, ok := node[p]
			if !ok {
				return nil, fmt.Errorf("路径不存在: %s", p)
			}
			cur = next
		case []any:
			idx, err := atoi(p)
			if err != nil || idx < 0 || idx >= len(node) {
				return nil, fmt.Errorf("数组下标无效: %s", p)
			}
			cur = node[idx]
		default:
			return nil, fmt.Errorf("无法继续解析路径: %s", p)
		}
	}
	return cur, nil
}

func atoi(s string) (int, error) {
	n := 0
	if s == "" {
		return 0, fmt.Errorf("empty")
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not number")
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

// EnsureValidJSON 供测试或其它包复用
func EnsureValidJSON(b []byte) error {
	var v any
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return dec.Decode(&v)
}
