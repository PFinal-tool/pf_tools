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
 * @Date: 2023/11/8
 * @Desc:
 * @Project: pf_tools
 */

type Weather struct{}

type IPInfo struct {
	Status string `json:"status"`
	Data   []IP   `json:"data`
}

type IP struct {
	Location string `json:"location"`
}

type WeatherRes struct {
	City    string        `json:"city"`
	Weather []WeatherInfo `json:"weather"`
}
type WeatherInfo struct {
	Date    string `json:"date"`
	Weather string `json:"weather"`
	Temp    string `json:"temp"`
	W       string `json:"w"`
	Wind    string `json:"wind"`
}

func (w *Weather) GetWeather(city string) {
	if city == "" {
		city = w._getLocalCity()
		// fmt.Println(city)
	}
	// fmt.Println(city)
	res := getWeatherInfo(city)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	table := widgets.NewTable()
	table.Title = res.City + "天气"
	table.BorderStyle = ui.NewStyle(ui.ColorRed)
	table.Rows = [][]string{
		[]string{"日期", "天气", "风向", "体感温度"},
	}
	for _, v := range res.Weather {
		table.Rows = append(table.Rows, []string{v.Date, v.Weather, v.Wind, v.Temp})
	}
	table.TextStyle = ui.NewStyle(ui.ColorGreen)
	table.TitleStyle = ui.NewStyle(ui.ColorGreen)
	table.SetRect(0, 0, 60, 10)
	ui.Render(table)
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "c":

		}
	}
}

func (w *Weather) _getLocalCity() string {
	externalIp, _ := w._getClientIp()
	result := getIpInfo(externalIp)
	return extractCity(result.Data[0].Location)
}

func (w *Weather) _getClientIp() (string, error) {
	resp, err := http.Get("https://ipv4.netarm.com/")
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	content, _ := io.ReadAll(resp.Body)
	return string(content), nil
}

func getIpInfo(ip string) *IPInfo {
	url := "https://opendata.baidu.com/api.php?query=" + ip + "&co=&resource_id=6006&oe=utf8"
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	out, _ := io.ReadAll(resp.Body)
	var result IPInfo
	if err := json.Unmarshal(out, &result); err != nil {
		return nil
	}
	return &result
}

func extractCity(input string) string {
	cityKeyword := "市"
	provinceKeyword := "省"
	cityIndex := strings.Index(input, cityKeyword)

	if cityIndex != -1 {
		provinceIndex := strings.LastIndex(input[:cityIndex], provinceKeyword)
		if provinceIndex != -1 {
			return input[provinceIndex+len(provinceKeyword) : cityIndex+len(cityKeyword)]
		}
		return input[:cityIndex+len(cityKeyword)]
	}

	return "未找到市"
}

func getWeatherInfo(city string) (weatherResponse *WeatherRes) {
	api := "https://api.asilu.com/weather/?city=" + city
	resp, err := http.Get(api)
	if err != nil {
		return nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	out, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(out))
	if err := json.Unmarshal(out, &weatherResponse); err != nil {
		return weatherResponse
	}
	return weatherResponse
}
