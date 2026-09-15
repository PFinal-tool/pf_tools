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
	支持 -f 对文件做 MD5，-q 仅输出结果
`,
	Run: func(cmd *cobra.Command, args []string) {
		quiet, _ := cmd.Flags().GetBool("quiet")
		filePath, _ := cmd.Flags().GetString("file")
		md5Obj := pak.Md5{Quiet: quiet}

		if filePath != "" {
			md5Obj.EncipherFile(filePath)
			return
		}

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
	pfMd5Cmd.Flags().BoolP("quiet", "q", false, "仅输出结果，便于管道")
	pfMd5Cmd.Flags().StringP("file", "f", "", "对文件计算 MD5")
	rootCmd.AddCommand(pfMd5Cmd)
}
