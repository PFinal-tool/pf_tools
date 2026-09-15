package cmd

import (
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

var pfURLCmd = &cobra.Command{
	Use:   "pf_url",
	Short: "URL 编解码 / 拆解",
	Long: `URL 工具：
  pf_tools pf_url "https://example.com?name=张三"
  pf_tools pf_url "https://a.com?x=%E5%BC%A0" -d
  pf_tools pf_url "https://a.com/path?a=1&b=2" --parse
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		decode, _ := cmd.Flags().GetBool("decode")
		parse, _ := cmd.Flags().GetBool("parse")
		tool := pak.URLTool{Quiet: quiet, Decode: decode, Parse: parse}
		tool.Run(args[0])
	},
}

func init() {
	pfURLCmd.Flags().BoolP("quiet", "q", false, "仅输出结果，便于管道")
	pfURLCmd.Flags().BoolP("decode", "d", false, "URL 解码")
	pfURLCmd.Flags().Bool("parse", false, "拆解 URL 各部分")
	rootCmd.AddCommand(pfURLCmd)
}
