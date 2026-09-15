package pak

import (
	"fmt"
	"log"
	"strings"

	"github.com/atotto/clipboard"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: 统一结果展示（termui / quiet）
 * @Project: pf_tools
 */

// ResultView 统一结果输出：Quiet 模式下只打印到 stdout，便于管道。
type ResultView struct {
	Quiet bool
}

func (r *ResultView) Print(text string) {
	fmt.Println(text)
}

func (r *ResultView) ShowText(title, text, copyText string, success bool) {
	if r.Quiet {
		out := copyText
		if out == "" {
			out = text
		}
		if out != "" {
			fmt.Println(strings.TrimRight(out, "\n"))
		}
		return
	}
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	lines := splitStringByLength(text, 100)
	if len(lines) == 0 {
		lines = []string{""}
	}
	p := widgets.NewParagraph()
	p.Title = title
	p.Text = strings.Join(lines, "\n")
	p.WrapText = true
	if success {
		p.TextStyle.Fg = ui.ColorGreen
		p.BorderStyle.Fg = ui.ColorGreen
	} else {
		p.TextStyle.Fg = ui.ColorRed
		p.BorderStyle.Fg = ui.ColorRed
	}
	width := 105
	if w := len(text) + 8; w < width && w > 30 {
		width = w
	}
	if width < 30 {
		width = 30
	}
	p.SetRect(0, 0, width, len(lines)+2)
	ui.Render(p)
	waitCopyOrQuit(copyText, success)
}

func (r *ResultView) ShowTable(title string, rows [][]string, copyText string) {
	if r.Quiet {
		for i, row := range rows {
			if i == 0 {
				continue
			}
			fmt.Println(strings.Join(row, "\t"))
		}
		if copyText != "" && len(rows) <= 1 {
			fmt.Println(copyText)
		}
		return
	}
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table := widgets.NewTable()
	table.Title = title
	table.Rows = rows
	table.TextStyle = ui.NewStyle(ui.ColorGreen)
	table.TitleStyle = ui.NewStyle(ui.ColorGreen)
	height := len(rows) + 2
	if height < 5 {
		height = 5
	}
	table.SetRect(0, 0, 100, height)
	ui.Render(table)
	waitCopyOrQuit(copyText, true)
}

func waitCopyOrQuit(copyText string, allowCopy bool) {
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "c":
			if allowCopy && copyText != "" {
				_ = clipboard.WriteAll(strings.TrimRight(copyText, "\n"))
			}
			return
		}
	}
}
