package cmd

import (
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

// pfTCmd represents the pfT command
var pfTCmd = &cobra.Command{
	Use:   "pf_t",
	Short: "当前时间戳",
	Long:  `获取当前时间戳 如果传递时间则转化对应的时间为时间戳`,
	Run: func(cmd *cobra.Command, args []string) {
		stime := pak.STime{}
		if len(args) > 0 {
			stime.SetFormatTime(args[0])
		} else {
			stime.GetCurrentTime()
		}
	},
}

func init() {
	rootCmd.AddCommand(pfTCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pfTCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pfTCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
