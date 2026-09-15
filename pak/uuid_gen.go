package pak

import (
	"crypto/rand"
	"fmt"
	"strings"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: UUID 生成
 * @Project: pf_tools
 */

type UUIDTool struct {
	Quiet    bool
	Count    int
	NoHyphen bool
}

func (u *UUIDTool) Run() {
	view := ResultView{Quiet: u.Quiet}
	n := u.Count
	if n <= 0 {
		n = 1
	}
	if n > 100 {
		view.ShowText("UUID", "一次最多生成 100 个", "", false)
		return
	}

	list := make([]string, 0, n)
	for i := 0; i < n; i++ {
		id, err := newUUIDv4()
		if err != nil {
			view.ShowText("UUID", "生成失败: "+err.Error(), "", false)
			return
		}
		if u.NoHyphen {
			id = strings.ReplaceAll(id, "-", "")
		}
		list = append(list, id)
	}

	text := strings.Join(list, "\n")
	view.ShowText("UUID", text, text, true)
}

func newUUIDv4() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant RFC4122
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
