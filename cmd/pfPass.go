package cmd

import (
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

var pfPassCmd = &cobra.Command{
	Use:   "pf_pass",
	Short: "生成随机密码",
	Long: `生成可读随机密码：
  pf_tools pf_pass
  pf_tools pf_pass -l 32
  pf_tools pf_pass -l 20 -s
  pf_tools pf_pass -l 12 -n 5 -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		quiet, _ := cmd.Flags().GetBool("quiet")
		length, _ := cmd.Flags().GetInt("length")
		special, _ := cmd.Flags().GetBool("special")
		count, _ := cmd.Flags().GetInt("count")
		tool := pak.PassTool{Quiet: quiet, Len: length, Special: special, Count: count}
		tool.Run()
	},
}

func init() {
	pfPassCmd.Flags().BoolP("quiet", "q", false, "仅输出结果，便于管道")
	pfPassCmd.Flags().IntP("length", "l", 16, "密码长度")
	pfPassCmd.Flags().BoolP("special", "s", false, "包含特殊字符")
	pfPassCmd.Flags().IntP("count", "n", 1, "生成数量")
	rootCmd.AddCommand(pfPassCmd)
}
