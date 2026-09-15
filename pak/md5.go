package pak

import (
	"crypto/md5"
	"fmt"
	"os"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2023/11/14
 * @Desc:
 * @Project: pf_tools
 */

type Md5 struct {
	Quiet bool
}

type Options struct {
	Salt     string
	Position int `json:"position"` // 0.为左边 1 为右边
}

func (m *Md5) Encipher(str string, options Options) {
	if str == "" {
		str = "123456"
	}
	if options.Position == 0 {
		str = options.Salt + str
	} else {
		str = str + options.Salt
	}
	data := []byte(str)
	has := md5.Sum(data)
	enStr := fmt.Sprintf("%x", has)

	view := ResultView{Quiet: m.Quiet}
	if m.Quiet {
		view.Print(enStr)
		return
	}
	rows := [][]string{
		{"待加密", "加密后"},
		{str, enStr},
	}
	view.ShowTable("MD5加密", rows, enStr)
}

func (m *Md5) EncipherFile(path string) {
	data, err := os.ReadFile(path)
	view := ResultView{Quiet: m.Quiet}
	if err != nil {
		view.ShowText("MD5加密", "读取文件失败: "+err.Error(), "", false)
		return
	}
	has := md5.Sum(data)
	enStr := fmt.Sprintf("%x", has)
	if m.Quiet {
		view.Print(enStr)
		return
	}
	rows := [][]string{
		{"文件", "MD5"},
		{path, enStr},
	}
	view.ShowTable("MD5加密", rows, enStr)
}
