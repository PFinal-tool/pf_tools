package pak

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2026/3/24
 * @Desc: SemVer bump
 * @Project: pf_tools
 */

type SemverTool struct {
	Quiet          bool
	Major          bool
	Minor          bool
	Patch          bool
	IncPrerelease  bool
}

var semverRe = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z.-]+))?(?:\+([0-9A-Za-z.-]+))?$`)

func (s *SemverTool) Run(version string) {
	view := ResultView{Quiet: s.Quiet}
	version = strings.TrimSpace(version)
	m := semverRe.FindStringSubmatch(version)
	if m == nil {
		view.ShowText("SemVer", "版本号无效，期望如 1.2.3 / v1.2.3-beta.1", "", false)
		return
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])
	pre := m[4]
	build := m[5]
	prefixV := strings.HasPrefix(version, "v")

	modeCount := 0
	if s.Major {
		modeCount++
	}
	if s.Minor {
		modeCount++
	}
	if s.Patch {
		modeCount++
	}
	if s.IncPrerelease {
		modeCount++
	}
	if modeCount == 0 {
		// 默认 patch bump
		s.Patch = true
	}
	if modeCount > 1 {
		view.ShowText("SemVer", "请只选择一种 bump: --major/--minor/--patch/--inc-prerelease", "", false)
		return
	}

	switch {
	case s.Major:
		major++
		minor, patch = 0, 0
		pre = ""
	case s.Minor:
		minor++
		patch = 0
		pre = ""
	case s.Patch:
		if pre != "" {
			// 有预发布时，patch bump 通常清掉预发布回到正式版
			pre = ""
		} else {
			patch++
		}
	case s.IncPrerelease:
		pre = bumpPrerelease(pre)
	}

	out := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	if pre != "" {
		out += "-" + pre
	}
	if build != "" && !s.Major && !s.Minor && !s.Patch {
		out += "+" + build
	}
	if prefixV {
		out = "v" + out
	}
	view.ShowText("SemVer", out, out, true)
}

func bumpPrerelease(pre string) string {
	if pre == "" {
		return "alpha.1"
	}
	parts := strings.Split(pre, ".")
	last := parts[len(parts)-1]
	if n, err := strconv.Atoi(last); err == nil {
		parts[len(parts)-1] = strconv.Itoa(n + 1)
		return strings.Join(parts, ".")
	}
	return pre + ".1"
}
