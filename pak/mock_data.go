package pak

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: 快速造测试数据
 * @Project: pf_tools
 */

type MockTool struct {
	Quiet  bool
	Count  int
	Fields string
}

var mockFirstNames = []string{"张", "王", "李", "赵", "刘", "陈", "杨", "黄", "周", "吴"}
var mockGivenNames = []string{"伟", "芳", "娜", "敏", "静", "强", "磊", "洋", "勇", "艳", "杰", "涛", "明", "超", "霞"}
var mockDomains = []string{"example.com", "test.local", "demo.dev"}

func (m *MockTool) Run(kind string) {
	view := ResultView{Quiet: m.Quiet}
	kind = strings.ToLower(strings.TrimSpace(kind))
	n := m.Count
	if n <= 0 {
		n = 1
	}
	if n > 100 {
		view.ShowText("Mock", "一次最多 100 条", "", false)
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	switch kind {
	case "", "user", "json":
		list := make([]map[string]any, 0, n)
		for i := 0; i < n; i++ {
			item := map[string]any{}
			if m.Fields != "" {
				item = mockByFields(r, m.Fields)
			} else {
				id, _ := newUUIDv4()
				name := mockFirstNames[r.Intn(len(mockFirstNames))] + mockGivenNames[r.Intn(len(mockGivenNames))]
				item = map[string]any{
					"id":         id,
					"name":       name,
					"email":      fmt.Sprintf("user%d@%s", r.Intn(10000), mockDomains[r.Intn(len(mockDomains))]),
					"phone":      mockPhone(r),
					"created_at": time.Now().Unix(),
				}
			}
			list = append(list, item)
		}
		var raw any = list
		if n == 1 {
			raw = list[0]
		}
		b, _ := json.MarshalIndent(raw, "", "  ")
		text := string(b)
		view.ShowText("Mock", text, text, true)
	case "phone":
		lines := make([]string, 0, n)
		for i := 0; i < n; i++ {
			lines = append(lines, mockPhone(r))
		}
		text := strings.Join(lines, "\n")
		view.ShowText("Mock手机号", text, text, true)
	case "name", "zh":
		lines := make([]string, 0, n)
		for i := 0; i < n; i++ {
			lines = append(lines, mockFirstNames[r.Intn(len(mockFirstNames))]+mockGivenNames[r.Intn(len(mockGivenNames))])
		}
		text := strings.Join(lines, "\n")
		view.ShowText("Mock姓名", text, text, true)
	default:
		view.ShowText("Mock", "未知类型: "+kind+"\n可用: user/json/phone/name", "", false)
	}
}

func mockByFields(r *rand.Rand, fields string) map[string]any {
	item := map[string]any{}
	for _, f := range strings.Split(fields, ",") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		name, typ, _ := strings.Cut(f, ":")
		name = strings.TrimSpace(name)
		typ = strings.TrimSpace(typ)
		if typ == "" {
			typ = name
		}
		switch strings.ToLower(typ) {
		case "uuid", "id":
			id, _ := newUUIDv4()
			item[name] = id
		case "zh", "name":
			item[name] = mockFirstNames[r.Intn(len(mockFirstNames))] + mockGivenNames[r.Intn(len(mockGivenNames))]
		case "email":
			item[name] = fmt.Sprintf("user%d@%s", r.Intn(10000), mockDomains[r.Intn(len(mockDomains))])
		case "phone":
			item[name] = mockPhone(r)
		case "ts", "time", "created_at":
			item[name] = time.Now().Unix()
		case "int":
			item[name] = r.Intn(1000)
		case "bool":
			item[name] = r.Intn(2) == 1
		default:
			item[name] = typ
		}
	}
	return item
}

func mockPhone(r *rand.Rand) string {
	prefixes := []string{"130", "131", "132", "133", "135", "136", "137", "138", "139", "150", "151", "152", "155", "156", "157", "158", "159", "186", "187", "188", "189"}
	p := prefixes[r.Intn(len(prefixes))]
	return fmt.Sprintf("%s%08d", p, r.Intn(100000000))
}
