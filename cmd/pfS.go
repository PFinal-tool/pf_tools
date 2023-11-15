package cmd

import (
	"github.com/pfinal/pf_tools/pak"

	"github.com/spf13/cobra"
)

// pfSCmd represents the pfS command
var pfSCmd = &cobra.Command{
	Use:   "pf_s",
	Short: "社交拼写查询",
	Long:  `社交短语解释意思`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		speak := pak.Speak{}
		speak.GetSpeakInfo(args[0])
	},
}

func init() {
	rootCmd.AddCommand(pfSCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pfSCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pfSCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
