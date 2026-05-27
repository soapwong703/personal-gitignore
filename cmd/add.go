package cmd

import "github.com/spf13/cobra"

var addCmd = &cobra.Command{
	Use:                "add <pattern>",
	Short:              "Add a pattern",
	Long:               "Append a pattern to the selected ignore file if it is not already present.",
	Example:            "pgi add \"*.log\"",
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
		for _, p := range patterns {
			if p == args[0] {
				return nil
			}
		}
		patterns = append(patterns, args[0])
		return writePatterns(state.ignoreFile, patterns)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
