package cmd

import (
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

var pfCertCmd = &cobra.Command{
	Use:   "pf_cert [host]",
	Short: "查询域名 TLS 证书信息",
	Long: `查询 TLS 证书：
  pf_tools pf_cert example.com
  pf_tools pf_cert example.com:443 --days -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		days, _ := cmd.Flags().GetBool("days")
		tool := pak.CertTool{Quiet: quiet, Days: days}
		tool.Run(args[0])
	},
}

func init() {
	pfCertCmd.Flags().BoolP("quiet", "q", false, "仅输出结果，便于管道")
	pfCertCmd.Flags().Bool("days", false, "仅关注剩余天数")
	rootCmd.AddCommand(pfCertCmd)
}
