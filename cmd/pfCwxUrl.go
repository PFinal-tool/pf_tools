/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

// pfCwxUrlCmd represents the pfCwxUrl command
var pfCwxUrlCmd = &cobra.Command{
	Use:   "pf_cwx",
	Short: "微信域名拦截检测",
	Long:  `检测域名是否被微信拦截`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		if pak.CheckUrl(args[0]) == false {
			_ = cmd.Help()
			return
		}
		pak.GetWxUrlInfo(args[0])
	},
}

func init() {
	rootCmd.AddCommand(pfCwxUrlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pfCwxUrlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pfCwxUrlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
