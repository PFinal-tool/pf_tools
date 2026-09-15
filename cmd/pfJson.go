package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/pfinal/pf_tools/pak"
	"github.com/spf13/cobra"
)

var pfJSONCmd = &cobra.Command{
	Use:   "pf_json",
	Short: "JSON 格式化 / 压缩 / 校验",
	Long: `JSON 工具：
  pf_tools pf_json '{"a":1}'
  pf_tools pf_json file.json -f
  pf_tools pf_json '{"a":1}' -c
  pf_tools pf_json '{"a":1}' -v
  pf_tools pf_json '{"a":{"b":1}}' -p a.b
  echo '{"a":1}' | pf_tools pf_json -q
`,
	Run: func(cmd *cobra.Command, args []string) {
		quiet, _ := cmd.Flags().GetBool("quiet")
		fromFile, _ := cmd.Flags().GetBool("file")
		compact, _ := cmd.Flags().GetBool("compact")
		validate, _ := cmd.Flags().GetBool("validate")
		path, _ := cmd.Flags().GetString("path")

		tool := pak.JSONTool{
			Quiet:    quiet,
			Compact:  compact,
			Validate: validate,
			Path:     path,
		}

		input := ""
		if len(args) > 0 {
			input = args[0]
		} else {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				b, err := io.ReadAll(os.Stdin)
				if err == nil {
					input = strings.TrimSpace(string(b))
				}
			}
		}
		if input == "" {
			_ = cmd.Help()
			return
		}
		tool.Run(input, fromFile)
	},
}

func init() {
	pfJSONCmd.Flags().BoolP("quiet", "q", false, "仅输出结果，便于管道")
	pfJSONCmd.Flags().BoolP("file", "f", false, "参数视为文件路径")
	pfJSONCmd.Flags().BoolP("compact", "c", false, "压缩 JSON")
	pfJSONCmd.Flags().BoolP("validate", "v", false, "仅校验是否合法")
	pfJSONCmd.Flags().StringP("path", "p", "", "按路径取值，如 data.items.0.id")
	rootCmd.AddCommand(pfJSONCmd)
}
