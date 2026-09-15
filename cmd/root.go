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
	Long: `基于 go 开发的小工具集合

基础:
   - pf_wt / pf_m / pf_s          天气、手机归属地、网络词汇
   - pf_cd / pf_mk / pf_cwx       清 DS_Store、建日期目录、微信链接检测

编解码:
   - pf_md5 / pf_b64 / pf_t       MD5、Base64、时间戳
   - pf_json / pf_url / pf_jwt    JSON、URL、JWT
   - pf_hex / pf_case / pf_unicode  Hex、命名风格、Unicode

生成器:
   - pf_uuid / pf_rand / pf_pass  UUID、随机串、密码
   - pf_mock                      测试数据

排障:
   - pf_whichport / pf_killport   端口占用查询 / 按端口杀进程
   - pf_cert / pf_ip / pf_dns     证书、IP、DNS
   - pf_http / pf_expand          HTTP 探活、短链展开
   - pf_cidr / pf_semver          网段判断、版本 bump
   - pf_cron / pf_status          Cron 解读、HTTP 状态码
   - pf_cookie / pf_header / pf_clip  Cookie、Header、剪贴板流水线

多数命令支持 -q 仅输出结果，便于管道。
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
