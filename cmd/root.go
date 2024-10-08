package cmd

import (
	"github.com/fatih/color"
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pf_tools",
	Short: "开发小工具",
	Long: `基于go开发的 小工具集合
   - pf_tools pf_wt 查询天气
   - pf_tools pf_m 手机归属地查询
   - pf_tools pf_md5 md5 小工具
   - pf_tools pf_b64 base64 小工具
   - pf_tools pf_s 查询网络词汇
   - pf_tools pf_t  获取当前时间戳
   - pf_tools pf_cd 清除目录中的.DS_Store 文件
   - pf_tools pf_cwx 微信域名拦截检测 文件
   - pf_tools pf_mk 快捷的创建时间文件夹
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetBool("version")
		if version {
			color.Green("当前版本号: " + pak.Version)
			return
		}

		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				return
			}
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Flags().BoolP("version", "v", false, "当前版本号")
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	//rootCmd.Flags().BoolP("version", "v", false, "")
}
