package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

var pfCronCmd = &cobra.Command{
	Use:   "pf_cron [expr]",
	Short: "Cron 表达式解读",
	Long: `解读 5 段 Cron，并给出接下来几次触发时间：
  pf_tools pf_cron "0 */5 * * *"
  pf_tools pf_cron "0 9 * * 1-5" -n 5 -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		next, _ := cmd.Flags().GetInt("next")
		tool := pak.CronTool{Quiet: quiet, Next: next}
		tool.Run(args[0])
	},
}

var pfStatusCmd = &cobra.Command{
	Use:   "pf_status [code]",
	Short: "HTTP 状态码速查",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		t := pak.StatusTool{Quiet: quiet}; t.Run(args[0])
	},
}

var pfHexCmd = &cobra.Command{
	Use:   "pf_hex [text]",
	Short: "十六进制编解码",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		decode, _ := cmd.Flags().GetBool("decode")
		t := pak.HexTool{Quiet: quiet, Decode: decode}; t.Run(args[0])
	},
}

var pfCaseCmd = &cobra.Command{
	Use:   "pf_case [ident]",
	Short: "命名风格转换",
	Long: `camel / pascal / snake / kebab / export：
  pf_tools pf_case user_id --to camel
  pf_tools pf_case UserId --to snake -q
  pf_tools pf_case user-id
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		to, _ := cmd.Flags().GetString("to")
		t := pak.CaseTool{Quiet: quiet, To: to}; t.Run(args[0])
	},
}

var pfUnicodeCmd = &cobra.Command{
	Use:   "pf_unicode [text]",
	Short: "Unicode 码点 / 转义",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		decode, _ := cmd.Flags().GetBool("decode")
		t := pak.UnicodeTool{Quiet: quiet, Decode: decode}; t.Run(args[0])
	},
}

func init() {
	pfCronCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfCronCmd.Flags().IntP("next", "n", 5, "显示接下来 N 次触发")
	pfStatusCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfHexCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfHexCmd.Flags().BoolP("decode", "d", false, "十六进制解码")
	pfCaseCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfCaseCmd.Flags().String("to", "", "目标风格: camel/pascal/snake/kebab/export")
	pfUnicodeCmd.Flags().BoolP("quiet", "q", false, "仅输出结果")
	pfUnicodeCmd.Flags().BoolP("decode", "d", false, "解码 \\uXXXX / U+XXXX")

	rootCmd.AddCommand(pfCronCmd, pfStatusCmd, pfHexCmd, pfCaseCmd, pfUnicodeCmd)
}

// 供同文件外其它命令复用 stdin 读取
func readMaybeStdin(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		b, err := io.ReadAll(os.Stdin)
		if err == nil {
			return strings.TrimSpace(string(b))
		}
	}
	return ""
}
