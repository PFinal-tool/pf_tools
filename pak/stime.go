package pak

import (
	"strconv"
	"strings"
	"time"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2023/12/11
 * @Desc:
 * @Project: pf_tools
 */

type STime struct {
	Quiet  bool
	Millis bool
	Zone   string
	Layout string
}

func (t *STime) loc() *time.Location {
	if t.Zone == "" {
		return time.Local
	}
	loc, err := time.LoadLocation(t.Zone)
	if err != nil {
		return time.Local
	}
	return loc
}

func (t *STime) layout() string {
	if t.Layout != "" {
		return t.Layout
	}
	return "2006-01-02 15:04:05"
}

func (t *STime) GetCurrentTime() {
	now := time.Now().In(t.loc())
	var val string
	var title string
	if t.Millis {
		val = strconv.FormatInt(now.UnixMilli(), 10)
		title = "毫秒时间戳"
	} else {
		val = strconv.FormatInt(now.Unix(), 10)
		title = "时间戳"
	}
	view := ResultView{Quiet: t.Quiet}
	view.ShowText(title, title+"："+val, val, true)
}

func (t *STime) Convert(input string) {
	input = strings.TrimSpace(input)
	view := ResultView{Quiet: t.Quiet}

	// 纯数字：时间戳 → 时间
	if isAllDigits(input) {
		n, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			view.ShowText("时间戳", "时间戳解析失败", "", false)
			return
		}
		var tm time.Time
		useMilli := t.Millis
		// 自动识别：13 位按毫秒，10 位按秒
		if len(input) >= 13 {
			useMilli = true
		} else if len(input) <= 10 {
			useMilli = false
		}
		if useMilli {
			tm = time.UnixMilli(n).In(t.loc())
		} else {
			tm = time.Unix(n, 0).In(t.loc())
		}
		formatted := tm.Format(t.layout())
		view.ShowText("时间", "时间："+formatted, formatted, true)
		return
	}

	// 时间字符串 → 时间戳
	formatTime, err := time.ParseInLocation(t.layout(), input, t.loc())
	if err != nil {
		view.ShowText("时间戳", "时间格式错误!!! 期望: "+t.layout(), "", false)
		return
	}
	var val string
	var title string
	if t.Millis {
		val = strconv.FormatInt(formatTime.UnixMilli(), 10)
		title = "毫秒时间戳"
	} else {
		val = strconv.FormatInt(formatTime.Unix(), 10)
		title = "时间戳"
	}
	view.ShowText(title, title+"："+val, val, true)
}

// 兼容旧调用
func (t *STime) SetFormatTime(str string) {
	t.Convert(str)
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
