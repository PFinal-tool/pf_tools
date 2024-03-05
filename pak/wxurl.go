package pak

import (
	"encoding/json"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"io"
	"log"
	"net/http"
	"net/url"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2024/3/4
 * @Desc:
 * @Project: pf_tools
 */

type WxUrlInfo struct {
	Data   string
	ReCode int
}

func GetWxUrlInfo(urlString string) {
	api := "https://cgi.urlsec.qq.com/index.php?m=url&a=validUrl&url=" + urlString
	resp, err := http.Get(api)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	out, _ := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	table := widgets.NewTable()
	table.Title = "微信URL安全检测"
	// table.BorderStyle = ui.NewStyle(ui.ColorRed)
	table.Rows = [][]string{
		[]string{"网址", "检测结果      "},
	}
	//fmt.Println(string(out))
	urlResponse := &WxUrlInfo{}
	if err := json.Unmarshal(out, &urlResponse); err != nil {
		fmt.Println("解析json失败:", err)
		return
	}
	enStr := "网址未被微信屏蔽"
	if urlResponse.ReCode == 0 {
		enStr = "网址被微信屏蔽"
	}
	table.Rows = append(table.Rows, []string{urlString, enStr})
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
		}
	}
}

func CheckUrl(urlString string) bool {
	// 解析url
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return false
	}
	// 检测是否有协议
	if parsedUrl.Scheme == "" {
		return false
	}
	// 检测是否有主机名
	if parsedUrl.Host == "" {
		return false
	}
	return true
}
