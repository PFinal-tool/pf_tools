package cmd

import (
	"strconv"

	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

var pfRandCmd = &cobra.Command{
	Use:   "pf_rand [length]",
	Short: "生成随机串（nonce/token）",
	Long: `生成随机串：
  pf_tools pf_rand
  pf_tools pf_rand 32
  pf_tools pf_rand 16 -t digit
  pf_tools pf_rand 24 -t base58 -c 3 -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		quiet, _ := cmd.Flags().GetBool("quiet")
		typ, _ := cmd.Flags().GetString("type")
		count, _ := cmd.Flags().GetInt("count")
		length := 32
		if len(args) > 0 {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				_ = cmd.Help()
				return
			}
			length = n
		}
		tool := pak.RandTool{Quiet: quiet, Len: length, Type: typ, Count: count}
		tool.Run()
	},
}

func init() {
	pfRandCmd.Flags().BoolP("quiet", "q", false, "仅输出结果，便于管道")
	pfRandCmd.Flags().StringP("type", "t", "alphanum", "类型: alphanum/digit/hex/base58")
	pfRandCmd.Flags().IntP("count", "c", 1, "生成数量")
	rootCmd.AddCommand(pfRandCmd)
}
