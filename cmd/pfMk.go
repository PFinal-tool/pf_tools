/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/pfinal/pf_tools/pak"
	"os"

	"github.com/spf13/cobra"
)

// pfMkCmd represents the pfMk command
var pfMkCmd = &cobra.Command{
	Use:   "pf_mk",
	Short: "快捷的创建时间文件夹",
	Long:  `方便快捷的创建名字为时间的文件夹 -r 递归创建文件夹 -l 参数 指定分隔字符`,
	Run: func(cmd *cobra.Command, args []string) {
		recursion, _ := cmd.Flags().GetBool("recursion")
		separation, _ := cmd.Flags().GetString("separation")
		path, _ := os.Getwd()
		mkdir := pak.Mkdir{Path: path, Recursion: recursion, Seperation: separation}
		mkdir.Done()
	},
}

func init() {
	pfMkCmd.Flags().BoolP("recursion", "r", false, "递归创建")
	pfMkCmd.Flags().StringP("separation", "l", "", "分隔符")
	rootCmd.AddCommand(pfMkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pfMkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pfMkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
