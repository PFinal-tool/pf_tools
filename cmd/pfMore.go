package cmd

import (
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

var pfCookieCmd = &cobra.Command{
	Use:   "pf_cookie [cookie]",
	Short: "解析 Cookie 字符串",
	Run: func(cmd *cobra.Command, args []string) {
		input := readMaybeStdin(args)
		if input == "" {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		t := pak.CookieTool{Quiet: quiet}; t.Run(input)
	},
}

var pfHeaderCmd = &cobra.Command{
	Use:   "pf_header",
	Short: "解析 HTTP Header 文本",
	Long:  `从参数或管道读取 Header 文本并拆解。`,
	Run: func(cmd *cobra.Command, args []string) {
		input := readMaybeStdin(args)
		if input == "" {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		t := pak.HeaderTool{Quiet: quiet}; t.Run(input)
	},
}

var pfClipCmd = &cobra.Command{
	Use:   "pf_clip [action]",
	Short: "剪贴板编解码流水线",
	Long: `读取剪贴板并处理：
  pf_tools pf_clip show
  pf_tools pf_clip md5 -q
  pf_tools pf_clip b64 --write
  pf_tools pf_clip jwt -q
  可用: show/md5/b64/b64d/url/urld/jwt/json
`,
	Run: func(cmd *cobra.Command, args []string) {
		action := "show"
		if len(args) > 0 {
			action = args[0]
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		write, _ := cmd.Flags().GetBool("write")
		t := pak.ClipTool{Quiet: quiet, Write: write}; t.Run(action)
	},
}

var pfExpandCmd = &cobra.Command{
	Use:   "pf_expand [url]",
	Short: "短链 / 重定向追踪",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		max, _ := cmd.Flags().GetInt("max")
		t := pak.ExpandTool{Quiet: quiet, Max: max}; t.Run(args[0])
	},
}

var pfMockCmd = &cobra.Command{
	Use:   "pf_mock [type]",
	Short: "快速造测试数据",
	Long: `生成测试数据：
  pf_tools pf_mock user -c 3 -q
  pf_tools pf_mock phone -c 5 -q
  pf_tools pf_mock json --fields id:uuid,name:zh,email,phone,created_at:ts -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		kind := "user"
		if len(args) > 0 {
			kind = args[0]
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		count, _ := cmd.Flags().GetInt("count")
		fields, _ := cmd.Flags().GetString("fields")
		t := pak.MockTool{Quiet: quiet, Count: count, Fields: fields}; t.Run(kind)
	},
}

var pfIPCmd = &cobra.Command{
	Use:   "pf_ip [ip]",
	Short: "IP 查询（公网/本地/信息）",
	Long: `IP 工具：
  pf_tools pf_ip
  pf_tools pf_ip -l
  pf_tools pf_ip 8.8.8.8 -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		quiet, _ := cmd.Flags().GetBool("quiet")
		local, _ := cmd.Flags().GetBool("local")
		ip := ""
		if len(args) > 0 {
			ip = args[0]
		}
		t := pak.IPTool{Quiet: quiet, Local: local}; t.Run(ip)
	},
}

func init() {
	pfCookieCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfHeaderCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfClipCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfClipCmd.Flags().Bool("write", false, "处理结果写回剪贴板")
	pfExpandCmd.Flags().BoolP("quiet", "q", false, "仅输出最终 URL")
	pfExpandCmd.Flags().Int("max", 10, "最大跳转次数")
	pfMockCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfMockCmd.Flags().IntP("count", "c", 1, "生成数量")
	pfMockCmd.Flags().String("fields", "", "自定义字段，如 id:uuid,name:zh,email")
	pfIPCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfIPCmd.Flags().BoolP("local", "l", false, "显示本地网卡 IP")

	rootCmd.AddCommand(pfCookieCmd, pfHeaderCmd, pfClipCmd, pfExpandCmd, pfMockCmd, pfIPCmd)
}
