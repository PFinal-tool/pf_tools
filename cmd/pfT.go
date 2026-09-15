package cmd

import (
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

// pfTCmd represents the pfT command
var pfTCmd = &cobra.Command{
	Use:   "pf_t",
	Short: "当前时间戳",
	Long: `获取当前时间戳；传入时间字符串则转为时间戳；传入纯数字时间戳则转为时间。
支持 -m 毫秒、-z 时区、-f 格式、-q 仅输出结果`,
	Run: func(cmd *cobra.Command, args []string) {
		quiet, _ := cmd.Flags().GetBool("quiet")
		millis, _ := cmd.Flags().GetBool("milli")
		zone, _ := cmd.Flags().GetString("zone")
		layout, _ := cmd.Flags().GetString("format")
		stime := pak.STime{Quiet: quiet, Millis: millis, Zone: zone, Layout: layout}
		if len(args) > 0 {
			stime.Convert(args[0])
		} else {
			stime.GetCurrentTime()
		}
	},
}

func init() {
	pfTCmd.Flags().BoolP("quiet", "q", false, "仅输出结果，便于管道")
	pfTCmd.Flags().BoolP("milli", "m", false, "使用毫秒时间戳")
	pfTCmd.Flags().StringP("zone", "z", "", "时区，如 Asia/Shanghai")
	pfTCmd.Flags().StringP("format", "f", "", "时间格式，默认 2006-01-02 15:04:05")
	rootCmd.AddCommand(pfTCmd)
}
