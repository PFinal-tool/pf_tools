package cmd

import (
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
	"strconv"
)

// pfMd5Cmd represents the pfMd5 command
var pfMd5Cmd = &cobra.Command{
	Use:   "pf_md5",
	Short: "字符串加密成md5",
	Long: `命令行加密md5 默认加密123456 方便调试
	1. 第一个参数是要加密的字符串
	2. 第二个参数是盐
	3. 第三个参数 0或者1 0 是盐在左边 1 是盐在右边
`,
	Run: func(cmd *cobra.Command, args []string) {
		md5Obj := pak.Md5{}
		if len(args) > 3 {
			_ = cmd.Help()
			return
		}
		if len(args) == 0 {
			md5Obj.Encipher("", pak.Options{Salt: "", Position: 0})
			return
		}
		if len(args) == 1 {
			md5Obj.Encipher(args[0], pak.Options{Salt: "", Position: 0})
			return
		}
		if len(args) == 2 {
			md5Obj.Encipher(args[0], pak.Options{Salt: args[1], Position: 0})
			return
		}
		if len(args) == 3 {
			position, _ := strconv.Atoi(args[2])
			set := map[int]string{
				0: "",
				1: "",
			}
			if _, ok := set[position]; !ok {
				_ = cmd.Help()
				return
			}
			md5Obj.Encipher(args[0], pak.Options{Salt: args[1], Position: position})
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pfMd5Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pfMd5Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pfMd5Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
