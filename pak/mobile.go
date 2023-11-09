package pak

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/zheng-ji/gophone"
	"log"
	"regexp"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2023/11/9
 * @Desc:
 * @Project: pf_tools
 */

type Mobile struct{}

type MobileInfo struct {
	PhoneNum string `json:"phone_num"`
	Province string `json:"province"`
	AreaZone string `json:"area_zone"`
	City     string `json:"city"`
	ZipCode  string `json:"zip_code"`
	CardType string `json:"card_type"`
}

func (m *Mobile) GetInfo(numb string) {
	pr, _ := gophone.Find(numb)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	l := widgets.NewList()
	l.Title = "号码详细信息"
	l.Rows = []string{
		"[0] 查询的号码:" + pr.PhoneNum,
		"[1] 号码运营商:" + pr.CardType,
		"[2] 号码所在省份:" + pr.Province,
		"[3] 号码所在城市:" + pr.City,
		"[4] 所在城市邮编:" + pr.ZipCode,
		"[5] 所在地区编码:" + pr.AreaZone,
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

func CheckMobile(phone string) bool {
	// 匹配规则
	// ^1第一位为一
	// [345789]{1} 后接一位345789 的数字
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1[345789]{1}\\d{9}$"
	// 正则调用规则
	reg := regexp.MustCompile(regRuler)
	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)
}
