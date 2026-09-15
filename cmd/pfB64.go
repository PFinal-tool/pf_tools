package cmd

import (
	"github.com/pfinal/pf_tools/pak"

	"github.com/spf13/cobra"
)

// pfB64Cmd represents the pfB64 command
var pfB64Cmd = &cobra.Command{
	Use:   "pf_b64",
	Short: "字符串base64格式",
	Long:  `字符串base64格式, 默认加密, -d 解密, -u 使用 Base64URL, -q 仅输出结果`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		decode, _ := cmd.Flags().GetBool("decode")
		quiet, _ := cmd.Flags().GetBool("quiet")
		urlSafe, _ := cmd.Flags().GetBool("url")
		cry := pak.Cryption{Quiet: quiet, URLSafe: urlSafe}
		if decode {
			cry.DecodeString(args[0])
		} else {
			cry.EncodeToString(args[0])
		}
	},
}

func init() {
	pfB64Cmd.Flags().BoolP("decode", "d", false, "解密")
	pfB64Cmd.Flags().BoolP("quiet", "q", false, "仅输出结果，便于管道")
	pfB64Cmd.Flags().BoolP("url", "u", false, "使用 Base64URL")
	rootCmd.AddCommand(pfB64Cmd)
}
