package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/pfinal/pf_tools/cmd"
)

var Logo = `
欢迎使用pf_tools小工具， 请按照下面的指示操作
`

func main() {
	// 仅在无子命令时展示欢迎语，避免污染 -q / 管道输出
	if len(os.Args) == 1 {
		colorPrint := color.New(color.Bold)
		colorPrint.Add(color.FgGreen)
		_, _ = colorPrint.Println(Logo)
	}
	cmd.Execute()
}
