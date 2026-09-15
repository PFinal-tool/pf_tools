package cmd

import (
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

var pfDNSCmd = &cobra.Command{
	Use:   "pf_dns [host]",
	Short: "DNS 查询",
	Long: `DNS 查询：
  pf_tools pf_dns example.com
  pf_tools pf_dns example.com -t MX
  pf_tools pf_dns example.com -t ALL -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		typ, _ := cmd.Flags().GetString("type")
		t := pak.DNSTool{Quiet: quiet, Type: typ}
		t.Run(args[0])
	},
}

var pfHTTPCmd = &cobra.Command{
	Use:   "pf_http [url]",
	Short: "轻量 HTTP 请求 / 探活",
	Long: `轻量 HTTP（看状态码、耗时、响应）：
  pf_tools pf_http https://example.com
  pf_tools pf_http https://httpbin.org/post -X POST -d '{"a":1}' -j
  pf_tools pf_http https://example.com --time -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		method, _ := cmd.Flags().GetString("method")
		body, _ := cmd.Flags().GetString("data")
		headers, _ := cmd.Flags().GetStringArray("header")
		jsonFmt, _ := cmd.Flags().GetBool("json")
		showTime, _ := cmd.Flags().GetBool("time")
		t := pak.HTTPTool{
			Quiet:    quiet,
			Method:   method,
			Body:     body,
			Header:   headers,
			JSONFmt:  jsonFmt,
			ShowTime: showTime,
		}
		t.Run(args[0])
	},
}

var pfCIDRCmd = &cobra.Command{
	Use:   "pf_cidr [ip] [cidr]",
	Short: "CIDR / IP 归属判断",
	Long: `CIDR 工具：
  pf_tools pf_cidr 10.0.0.8 10.0.0.0/24
  pf_tools pf_cidr 10.0.0.0/24
  pf_tools pf_cidr 10.0.0.0/24 --list -n 5 -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		quiet, _ := cmd.Flags().GetBool("quiet")
		list, _ := cmd.Flags().GetBool("list")
		count, _ := cmd.Flags().GetInt("count")
		t := pak.CIDRTool{Quiet: quiet, List: list, Count: count}
		t.Run(args)
	},
}

var pfSemverCmd = &cobra.Command{
	Use:   "pf_semver [version]",
	Short: "SemVer 版本号 bump",
	Long: `SemVer bump：
  pf_tools pf_semver 1.2.3 --patch
  pf_tools pf_semver 1.2.3 --minor
  pf_tools pf_semver v1.2.3-beta.1 --inc-prerelease -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		major, _ := cmd.Flags().GetBool("major")
		minor, _ := cmd.Flags().GetBool("minor")
		patch, _ := cmd.Flags().GetBool("patch")
		incPre, _ := cmd.Flags().GetBool("inc-prerelease")
		t := pak.SemverTool{
			Quiet:         quiet,
			Major:         major,
			Minor:         minor,
			Patch:         patch,
			IncPrerelease: incPre,
		}
		t.Run(args[0])
	},
}

func init() {
	pfDNSCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfDNSCmd.Flags().StringP("type", "t", "A", "记录类型: A/AAAA/CNAME/MX/TXT/NS/ALL")

	pfHTTPCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfHTTPCmd.Flags().StringP("method", "X", "GET", "HTTP 方法")
	pfHTTPCmd.Flags().StringP("data", "d", "", "请求体")
	pfHTTPCmd.Flags().StringArrayP("header", "H", nil, "请求头，可重复，如 -H 'Authorization: Bearer x'")
	pfHTTPCmd.Flags().BoolP("json", "j", false, "格式化 JSON 响应")
	pfHTTPCmd.Flags().Bool("time", false, "显示耗时")

	pfCIDRCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfCIDRCmd.Flags().Bool("list", false, "列出网段内 IP")
	pfCIDRCmd.Flags().IntP("count", "n", 5, "列表数量")

	pfSemverCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfSemverCmd.Flags().Bool("major", false, " bump major")
	pfSemverCmd.Flags().Bool("minor", false, " bump minor")
	pfSemverCmd.Flags().Bool("patch", false, " bump patch（默认）")
	pfSemverCmd.Flags().Bool("inc-prerelease", false, "递增预发布号")

	rootCmd.AddCommand(pfDNSCmd, pfHTTPCmd, pfCIDRCmd, pfSemverCmd)
}
