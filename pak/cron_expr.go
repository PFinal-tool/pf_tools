package pak

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: Cron 表达式解读
 * @Project: pf_tools
 */

type CronTool struct {
	Quiet bool
	Next  int
}

func (c *CronTool) Run(expr string) {
	view := ResultView{Quiet: c.Quiet}
	expr = strings.TrimSpace(expr)
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		view.ShowText("Cron", "需要 5 段表达式: 分 时 日 月 周\n示例: 0 */5 * * *", "", false)
		return
	}

	desc, err := describeCron(fields)
	if err != nil {
		view.ShowText("Cron", err.Error(), "", false)
		return
	}

	n := c.Next
	if n <= 0 {
		n = 5
	}
	if n > 20 {
		n = 20
	}
	nexts, err := nextCronTimes(fields, time.Now(), n)
	if err != nil {
		view.ShowText("Cron", err.Error(), "", false)
		return
	}

	var b strings.Builder
	b.WriteString("表达式: " + expr + "\n")
	b.WriteString("含义: " + desc + "\n")
	b.WriteString("接下来触发:\n")
	for i, t := range nexts {
		b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, t.Format("2006-01-02 15:04:05 Mon")))
	}
	text := strings.TrimRight(b.String(), "\n")
	if c.Quiet {
		view.Print(text)
		return
	}
	copyText := ""
	if len(nexts) > 0 {
		copyText = nexts[0].Format("2006-01-02 15:04:05")
	}
	view.ShowText("Cron", text, copyText, true)
}

func describeCron(fields []string) (string, error) {
	names := []string{"分", "时", "日", "月", "周"}
	parts := make([]string, 0, 5)
	for i, f := range fields {
		parts = append(parts, fmt.Sprintf("%s=%s", names[i], humanField(f, i)))
	}
	return strings.Join(parts, "，"), nil
}

func humanField(f string, idx int) string {
	if f == "*" {
		return "每" + []string{"分钟", "小时", "天", "月", "周"}[idx]
	}
	if strings.HasPrefix(f, "*/") {
		return "每 " + strings.TrimPrefix(f, "*/") + " " + []string{"分钟", "小时", "天", "月", "周"}[idx]
	}
	return f
}

func nextCronTimes(fields []string, from time.Time, n int) ([]time.Time, error) {
	matchers := make([]func(int) bool, 5)
	bounds := [][2]int{{0, 59}, {0, 23}, {1, 31}, {1, 12}, {0, 6}}
	for i, f := range fields {
		m, err := buildMatcher(f, bounds[i][0], bounds[i][1])
		if err != nil {
			return nil, fmt.Errorf("第 %d 段无效: %v", i+1, err)
		}
		matchers[i] = m
	}

	t := from.Truncate(time.Minute).Add(time.Minute)
	out := make([]time.Time, 0, n)
	// 最多扫描约 2 年分钟数，避免死循环
	limit := 366 * 24 * 60 * 2
	for i := 0; i < limit && len(out) < n; i++ {
		dow := int(t.Weekday()) // 0=Sunday
		if matchers[0](t.Minute()) &&
			matchers[1](t.Hour()) &&
			matchers[2](t.Day()) &&
			matchers[3](int(t.Month())) &&
			matchers[4](dow) {
			out = append(out, t)
		}
		t = t.Add(time.Minute)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("未找到未来触发时间，请检查表达式")
	}
	return out, nil
}

func buildMatcher(field string, min, max int) (func(int) bool, error) {
	if field == "*" {
		return func(int) bool { return true }, nil
	}
	allowed := map[int]bool{}
	for _, part := range strings.Split(field, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		step := 1
		body := part
		if strings.Contains(part, "/") {
			bits := strings.SplitN(part, "/", 2)
			body = bits[0]
			s, err := strconv.Atoi(bits[1])
			if err != nil || s <= 0 {
				return nil, fmt.Errorf("步长无效: %s", part)
			}
			step = s
		}
		if body == "*" {
			for v := min; v <= max; v += step {
				allowed[v] = true
			}
			continue
		}
		if strings.Contains(body, "-") {
			rg := strings.SplitN(body, "-", 2)
			a, err1 := strconv.Atoi(rg[0])
			b, err2 := strconv.Atoi(rg[1])
			if err1 != nil || err2 != nil || a > b {
				return nil, fmt.Errorf("范围无效: %s", part)
			}
			for v := a; v <= b; v += step {
				if v >= min && v <= max {
					allowed[v] = true
				}
			}
			continue
		}
		v, err := strconv.Atoi(body)
		if err != nil {
			return nil, fmt.Errorf("数值无效: %s", part)
		}
		if step == 1 {
			if v < min || v > max {
				return nil, fmt.Errorf("越界: %d", v)
			}
			allowed[v] = true
		} else {
			for x := v; x <= max; x += step {
				if x >= min {
					allowed[x] = true
				}
			}
		}
	}
	return func(v int) bool { return allowed[v] }, nil
}
