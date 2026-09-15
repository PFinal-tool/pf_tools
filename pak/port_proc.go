package pak

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: 端口占用查询 / 按端口杀进程
 * @Project: pf_tools
 */

type PortTool struct {
	Quiet bool
	Force bool
}

type portProc struct {
	PID  string
	Name string
	Cmd  string
}

func (p *PortTool) Which(port int) {
	view := ResultView{Quiet: p.Quiet}
	if port <= 0 || port > 65535 {
		view.ShowText("端口", "端口无效", "", false)
		return
	}
	list, err := findProcsByPort(port)
	if err != nil {
		view.ShowText("端口", err.Error(), "", false)
		return
	}
	if len(list) == 0 {
		view.ShowText("端口", fmt.Sprintf("端口 %d 未被占用", port), "", true)
		return
	}
	rows := [][]string{{"PID", "进程", "命令"}}
	pids := make([]string, 0, len(list))
	for _, item := range list {
		rows = append(rows, []string{item.PID, item.Name, truncate(item.Cmd, 60)})
		pids = append(pids, item.PID)
	}
	view.ShowTable(fmt.Sprintf("端口 %d", port), rows, strings.Join(pids, " "))
}

func (p *PortTool) Kill(port int) {
	view := ResultView{Quiet: p.Quiet}
	if port <= 0 || port > 65535 {
		view.ShowText("KillPort", "端口无效", "", false)
		return
	}
	list, err := findProcsByPort(port)
	if err != nil {
		view.ShowText("KillPort", err.Error(), "", false)
		return
	}
	if len(list) == 0 {
		view.ShowText("KillPort", fmt.Sprintf("端口 %d 未被占用", port), "", true)
		return
	}

	sig := "-15"
	if p.Force {
		sig = "-9"
	}
	killed := make([]string, 0, len(list))
	failed := make([]string, 0)
	seen := map[string]bool{}
	for _, item := range list {
		if seen[item.PID] {
			continue
		}
		seen[item.PID] = true
		cmd := exec.Command("kill", sig, item.PID)
		if err := cmd.Run(); err != nil {
			failed = append(failed, fmt.Sprintf("%s(%s)", item.PID, item.Name))
			continue
		}
		killed = append(killed, fmt.Sprintf("%s(%s)", item.PID, item.Name))
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("端口 %d\n", port))
	if len(killed) > 0 {
		b.WriteString("已结束: " + strings.Join(killed, ", ") + "\n")
	}
	if len(failed) > 0 {
		b.WriteString("失败: " + strings.Join(failed, ", "))
		view.ShowText("KillPort", strings.TrimSpace(b.String()), strings.Join(killed, " "), false)
		return
	}
	view.ShowText("KillPort", strings.TrimSpace(b.String()), strings.Join(killed, " "), true)
}

func findProcsByPort(port int) ([]portProc, error) {
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		return nil, fmt.Errorf("当前系统暂不支持: %s", runtime.GOOS)
	}
	// -nP 避免反查，输出更稳
	out, err := exec.Command("lsof", "-nP", fmt.Sprintf("-iTCP:%d", port), "-sTCP:LISTEN").Output()
	if err != nil {
		// lsof 无匹配时退出码非 0
		if len(out) == 0 {
			if _, ok := err.(*exec.ExitError); ok {
				// 再试一遍不限 LISTEN（兼容部分场景）
				out2, err2 := exec.Command("lsof", "-nP", fmt.Sprintf("-i:%d", port)).Output()
				if err2 != nil {
					if len(out2) == 0 {
						if _, ok2 := err2.(*exec.ExitError); ok2 {
							return nil, nil
						}
					}
					return nil, fmt.Errorf("查询失败，请确认已安装 lsof")
				}
				return parseLsof(string(out2)), nil
			}
		}
		return nil, fmt.Errorf("查询失败，请确认已安装 lsof")
	}
	return parseLsof(string(out)), nil
}

func parseLsof(out string) []portProc {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) <= 1 {
		return nil
	}
	result := make([]portProc, 0)
	seen := map[string]bool{}
	for i, line := range lines {
		if i == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name := fields[0]
		pid := fields[1]
		if _, err := strconv.Atoi(pid); err != nil {
			continue
		}
		key := pid + ":" + name
		if seen[key] {
			continue
		}
		seen[key] = true
		cmd := name
		if full := processCommand(pid); full != "" {
			cmd = full
		}
		result = append(result, portProc{PID: pid, Name: name, Cmd: cmd})
	}
	return result
}

func processCommand(pid string) string {
	out, err := exec.Command("ps", "-p", pid, "-o", "command=").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
