package pak

import (
	"bytes"
	"embed"
	"encoding/binary"
	"errors"
	"fmt"
	. "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"regexp"
	"strconv"
	"strings"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2023/11/9
 * @Desc:
 * @Project: pf_tools
 */
const (
	IntLen           = 4
	CharLen          = 1
	PhoneIndexLength = 9
	CHUNK            = 100
	PhoneDat         = "phone.dat"
)

//go:embed phone.dat
var fsContent embed.FS

type PhoneRecord struct {
	PhoneNum string
	Province string
	City     string
	ZipCode  string
	AreaZone string
	CardType string
}

var content []byte

func init() {
	var err error
	content, err = fsContent.ReadFile(PhoneDat)
	if err != nil {
		panic(err)
	}
}

func Display() {
	fmt.Println(getVersion())
	fmt.Println(getTotalRecord())
	fmt.Println(getFirstRecordOffset())
}

func (pr PhoneRecord) String() string {
	_str := fmt.Sprintf("PhoneNum: %s\nAreaZone: %s\nCardType: %s\nCity: %s\nZipCode: %s\nProvince: %s\n", pr.PhoneNum, pr.AreaZone, pr.CardType, pr.City, pr.ZipCode, pr.Province)
	return _str
}

func getVersion() string {
	return string(content[0:IntLen])
}

func getTotalRecord() int32 {
	total := (int32(len(content)) - getFirstRecordOffset()) / PhoneIndexLength
	return total
}

func getFirstRecordOffset() int32 {
	var offset int32
	buffer := bytes.NewBuffer(content[IntLen : IntLen*2])
	_ = binary.Read(buffer, binary.LittleEndian, &offset)
	return offset
}

func getIndexRecord(offset int32) (phonePrefix int32, recordOffset int32, cardType byte) {
	buffer := bytes.NewBuffer(content[offset : offset+IntLen])
	_ = binary.Read(buffer, binary.LittleEndian, &phonePrefix)
	buffer = bytes.NewBuffer(content[offset+IntLen : offset+IntLen*2])
	_ = binary.Read(buffer, binary.LittleEndian, &recordOffset)
	buffer = bytes.NewBuffer(content[offset+IntLen*2 : offset+IntLen*2+CharLen])
	_ = binary.Read(buffer, binary.LittleEndian, &cardType)
	return
}

func getOpCompany(cardtype byte) string {
	var card_str = ""
	switch cardtype {
	case '1':
		card_str = "移动"
	case '2':
		card_str = "联通"
	case '3':
		card_str = "电信"
	case '4':
		card_str = "电信虚拟运营商"
	case '5':
		card_str = "联通虚拟运营商"
	default:
		card_str = "移动虚拟运营商"
	}
	return card_str
}

func Find(phoneNum string) (pr *PhoneRecord, err error) {
	err = nil
	if len(phoneNum) < 7 || len(phoneNum) > 11 {
		return nil, errors.New("illegal phone length")
	}

	var left int32 = 0
	phone_seven_int, _ := strconv.ParseInt(phoneNum[0:7], 10, 32)
	phone_seven_int32 := int32(phone_seven_int)
	total_len := int32(len(content))
	right := getTotalRecord()
	firstPhoneRecordOffset := getFirstRecordOffset()
	for {
		if left > right {
			break
		}
		mid := (left + right) / 2
		current_offset := firstPhoneRecordOffset + mid*PhoneIndexLength

		if current_offset >= total_len {
			break
		}
		cur_phone, record_offset, card_type := getIndexRecord(current_offset)
		if cur_phone > phone_seven_int32 {
			right = mid - 1
		} else if cur_phone < phone_seven_int32 {
			left = mid + 1
		} else {
			s := record_offset
			e := record_offset + int32(strings.Index(string(content[record_offset:record_offset+CHUNK]), "\000"))
			record_content := string(content[s:e])
			_tmp := strings.Split(record_content, "|")
			card_str := getOpCompany(card_type)
			pr = &PhoneRecord{
				PhoneNum: phoneNum,
				Province: _tmp[0],
				City:     _tmp[1],
				ZipCode:  _tmp[2],
				AreaZone: _tmp[3],
				CardType: card_str,
			}
			return
		}
	}
	return nil, errors.New("num not found")
}

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
	if err := Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer Close()
	pr, err := Find(numb)
	if err != nil {
		//fmt.Println(err)
		p := widgets.NewParagraph()
		p.Title = "查询结果"
		p.Text = "查询失败 按 Q 退出"
		p.TextStyle.Fg = ColorGreen
		p.BorderStyle.Fg = ColorRed
		p.SetRect(0, 0, 30, 3)
		Render(p)

		uiEvents := PollEvents()
		for {
			e := <-uiEvents
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		}
		return
	}

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
	l.TextStyle = NewStyle(ColorGreen)
	l.TitleStyle = NewStyle(ColorGreen)
	l.WrapText = false
	l.SetRect(0, 0, 40, 8)
	Render(l)
	previousKey := ""
	uiEvents := PollEvents()
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
		Render(l)
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
