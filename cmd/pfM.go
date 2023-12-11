package cmd

import (
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

// pfMCmd represents the pfM command
var pfMCmd = &cobra.Command{
	Use:   "pf_m",
	Short: "查询手机归属地",
	Long:  `查询对应手机的归属地`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		if pak.CheckMobile(args[0]) == false {
			_ = cmd.Help()
			return
		}
		m := pak.Mobile{}
		m.GetInfo(args[0])
	},
}

func init() {
	rootCmd.AddCommand(pfMCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pfMCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pfMCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
