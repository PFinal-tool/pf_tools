package cmd

import (
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

var pfUUIDCmd = &cobra.Command{
	Use:   "pf_uuid",
	Short: "生成 UUID",
	Long: `生成 UUID v4：
  pf_tools pf_uuid
  pf_tools pf_uuid -c 5
  pf_tools pf_uuid -n -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		quiet, _ := cmd.Flags().GetBool("quiet")
		count, _ := cmd.Flags().GetInt("count")
		noHyphen, _ := cmd.Flags().GetBool("no-hyphen")
		tool := pak.UUIDTool{Quiet: quiet, Count: count, NoHyphen: noHyphen}
		tool.Run()
	},
}

func init() {
	pfUUIDCmd.Flags().BoolP("quiet", "q", false, "仅输出结果，便于管道")
	pfUUIDCmd.Flags().IntP("count", "c", 1, "生成数量")
	pfUUIDCmd.Flags().BoolP("no-hyphen", "n", false, "去掉横杠")
	rootCmd.AddCommand(pfUUIDCmd)
}
