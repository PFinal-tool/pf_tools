package pak

import (
	"encoding/json"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"io"
	"log"
	"net/http"
	"strings"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2023/11/15
 * @Desc:
 * @Project: pf_tools
 */

type Speak struct{}
type SpeakInfo struct {
	Name  string
	Trans []string
}

func (s *Speak) GetSpeakInfo(str string) {
	info := s.getSpeakInfo(str)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	l := widgets.NewList()
	l.Title = "查询信息"
	l.Rows = []string{
		"[0] 查询的词:" + info[0].Name,
	}
	for _, v := range info[0].Trans {
		l.Rows = append(l.Rows, "    "+v)
	}
	l.TextStyle = ui.NewStyle(ui.ColorGreen)
	l.TitleStyle = ui.NewStyle(ui.ColorGreen)
	l.WrapText = false
	l.SetRect(0, 0, 40, 8)
	ui.Render(l)
	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}
		ui.Render(l)
	}
}

func (s *Speak) getSpeakInfo(str string) (info []SpeakInfo) {
	api := "https://lab.magiconch.com/api/nbnhhsh/guess"
	data := make(map[string]string)
	//var data map[string]string
	data["text"] = str
	params, _ := json.Marshal(data)
	resp, err := http.Post(api, "application/json", strings.NewReader(string(params)))
	if err != nil {
		return nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	out, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(out, &info)
	if err != nil {
		return nil
	}
	return
}
