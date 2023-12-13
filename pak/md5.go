package pak

import (
	"crypto/md5"
	"fmt"
	"github.com/atotto/clipboard"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2023/11/14
 * @Desc:
 * @Project: pf_tools
 */

type Md5 struct{}
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
	//fmt.Println(enStr)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	table := widgets.NewTable()
	table.Title = "MD5加密"
	// table.BorderStyle = ui.NewStyle(ui.ColorRed)
	table.Rows = [][]string{
		[]string{"待加密", "加密后      "},
	}
	table.Rows = append(table.Rows, []string{str, enStr})
	table.TextStyle = ui.NewStyle(ui.ColorGreen)
	table.TitleStyle = ui.NewStyle(ui.ColorGreen)
	table.SetRect(0, 0, 80, 5)
	ui.Render(table)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q":
			return
		case "c":
			_ = clipboard.WriteAll(enStr)
			return
		}
	}
}
