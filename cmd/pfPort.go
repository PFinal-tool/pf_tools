package cmd

import (
	"strconv"

	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

var pfWhichPortCmd = &cobra.Command{
	Use:   "pf_whichport [port]",
	Short: "查看谁占用了端口",
	Long: `查看占用指定端口的进程：
  pf_tools pf_whichport 8080
  pf_tools pf_whichport 3306 -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		port, err := strconv.Atoi(args[0])
		if err != nil {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		tool := pak.PortTool{Quiet: quiet}
		tool.Which(port)
	},
}

var pfKillPortCmd = &cobra.Command{
	Use:   "pf_killport [port]",
	Short: "按端口结束进程",
	Long: `结束占用指定端口的进程：
  pf_tools pf_killport 8080
  pf_tools pf_killport 8080 -f
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		port, err := strconv.Atoi(args[0])
		if err != nil {
			_ = cmd.Help()
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		force, _ := cmd.Flags().GetBool("force")
		tool := pak.PortTool{Quiet: quiet, Force: force}
		tool.Kill(port)
	},
}

func init() {
	pfWhichPortCmd.Flags().BoolP("quiet", "q", false, "仅输出结果，便于管道")
	pfKillPortCmd.Flags().BoolP("quiet", "q", false, "仅输出结果，便于管道")
	pfKillPortCmd.Flags().BoolP("force", "f", false, "强制结束 (SIGKILL)")
	rootCmd.AddCommand(pfWhichPortCmd)
	rootCmd.AddCommand(pfKillPortCmd)
}
