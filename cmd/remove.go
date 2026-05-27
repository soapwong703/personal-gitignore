package cmd

import "github.com/spf13/cobra"

var removeCmd = &cobra.Command{
	Use:                "remove <pattern>",
	Short:              "Remove a pattern",
	Long:               "Remove matching pattern lines from the selected ignore file.",
	Example:            "pgi remove \"*.log\"",
	Args:               cobra.ExactArgs(1),
	DisableFlagParsing: true,
	PreRunE:            prepareRuntimeState,
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getRuntimeState(cmd)
		if err != nil {
			return err
		}
		patterns, err := readPatterns(state.ignoreFile)
		if err != nil {
			return err
		}
		updated := make([]string, 0, len(patterns))
		for _, p := range patterns {
			if p != args[0] {
				updated = append(updated, p)
			}
		}
		return writePatterns(state.ignoreFile, updated)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
