package cmd

import (
	"fmt"
	"github.com/pfinal/pf_tools/pak"

	"github.com/spf13/cobra"
)

// pfB64Cmd represents the pfB64 command
var pfB64Cmd = &cobra.Command{
	Use:   "pf_b64",
	Short: "字符串base64格式",
	Long:  `字符串base64格式, 默认加密, 第二个参数-d 则是解密,方便快捷使用`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		cry := pak.Cryption{}
		decode, _ := cmd.Flags().GetBool("decode")
		fmt.Println(decode)
		if decode {
			cry.DecodeString(args[0])
		} else {
			cry.EncodeToString(args[0])
		}
	},
}

func init() {
	pfB64Cmd.Flags().BoolP("decode", "d", false, "解密")
	rootCmd.AddCommand(pfB64Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pfB64Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pfB64Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
