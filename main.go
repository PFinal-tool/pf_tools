package main

import (
	"github.com/fatih/color"
	"github.com/pfinal/pf_tools/cmd"
)

var Logo = `
欢迎使用pf_tools小工具， 请按照下面的指示操作
`

func main() {
	// 实例化一个新的color对象，设置前景色为红色，背景色为绿色，文字斜体
	colorPrint := color.New(color.Bold)
	colorPrint.Add(color.FgGreen)
	_, _ = colorPrint.Println(Logo)
	cmd.Execute()
}
