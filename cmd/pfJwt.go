package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

var pfJWTCmd = &cobra.Command{
	Use:   "pf_jwt",
	Short: "JWT 解码（不校验签名）",
	Long: `解码 JWT header/payload，不验证签名。
  pf_tools pf_jwt eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
  pf_tools pf_jwt $TOKEN --exp
  echo $TOKEN | pf_tools pf_jwt -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		quiet, _ := cmd.Flags().GetBool("quiet")
		exp, _ := cmd.Flags().GetBool("exp")
		tool := pak.JWTTool{Quiet: quiet, Exp: exp}

		token := ""
		if len(args) > 0 {
			token = args[0]
		} else {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				b, err := io.ReadAll(os.Stdin)
				if err == nil {
					token = strings.TrimSpace(string(b))
				}
			}
		}
		if token == "" {
			_ = cmd.Help()
			return
		}
		tool.Run(token)
	},
}

func init() {
	pfJWTCmd.Flags().BoolP("quiet", "q", false, "仅输出结果，便于管道")
	pfJWTCmd.Flags().Bool("exp", false, "仅检查是否过期")
	rootCmd.AddCommand(pfJWTCmd)
}
