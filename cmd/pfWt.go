package cmd

import (
	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

// pfWtCmd represents the pfWt command
var pfWtCmd = &cobra.Command{
	Use:   "pf_wt",
	Short: "查询天气",
	Long:  `查询对应参数的天气`,
	Run: func(cmd *cobra.Command, args []string) {
		weather := pak.Weather{}
		if len(args) == 0 {
			weather.GetWeather("")
		} else {
			weather.GetWeather(args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(pfWtCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pfWtCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pfWtCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
