/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
	"os"
)

// pfCdCmd represents the pfCd command
var pfCdCmd = &cobra.Command{
	Use:   "pf_cd",
	Short: "清除目录中的.DS_Store 文件",
	Long:  `清除mac 目录中生成的.DS_Store 文件`,
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		if len(args) > 0 {
			fmt.Println(args[0])
		} else {
			path, _ = os.Getwd()
			fmt.Printf("清除目录的.DS_Store 文件: %s\n", path)
		}
		clearPath := pak.ClearPath{Path: path}
		clearPath.ClearDotDSStore()
	},
}

func init() {
	rootCmd.AddCommand(pfCdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pfCdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pfCdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
