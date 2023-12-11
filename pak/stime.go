package pak

import (
	"github.com/atotto/clipboard"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"strconv"
	"time"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2023/12/11
 * @Desc:
 * @Project: pf_tools
 */

type STime struct{}

func (t *STime) GetCurrentTime() {
	timeUnix := strconv.FormatInt(time.Now().Unix(), 10)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	p := widgets.NewParagraph()
	p.Title = "时间戳"
	p.Text = "时间戳：" + timeUnix
	p.TextStyle.Fg = ui.ColorGreen
	p.BorderStyle.Fg = ui.ColorGreen
	p.SetRect(0, 0, 30, 3)
	ui.Render(p)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "c":
			_ = clipboard.WriteAll(timeUnix)
			return
		}
	}
}

func (t *STime) SetFormatTime(str string) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	p := widgets.NewParagraph()
	p.Title = "时间戳"
	formatTime, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		p.Text = "时间格式错误!!!"
		p.TextStyle.Fg = ui.ColorRed
		p.BorderStyle.Fg = ui.ColorRed
		p.SetRect(0, 0, 30, 3)
		ui.Render(p)
		uiEvents := ui.PollEvents()
		for {
			e := <-uiEvents
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		}
	} else {
		timeUnix := strconv.FormatInt(formatTime.Unix(), 10)
		p.Text = "时间戳：" + timeUnix
		p.TextStyle.Fg = ui.ColorGreen
		p.BorderStyle.Fg = ui.ColorGreen
		p.SetRect(0, 0, 30, 3)
		ui.Render(p)
		uiEvents := ui.PollEvents()
		for {
			e := <-uiEvents
			switch e.ID {
			case "q", "<C-c>":
				return
			case "c":
				_ = clipboard.WriteAll(timeUnix)
				return
			}
		}
	}
}
