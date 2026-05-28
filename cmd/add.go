package cmd

import "github.com/spf13/cobra"

var addCmd = &cobra.Command{
	Use:     "add <pattern...>",
	Short:   "Add a pattern",
	Long:    "Append one or more patterns to the selected ignore file if they are not already present.",
	Example: "  pgi add \"*.log\"",
	Args:    cobra.MinimumNArgs(1),
	PreRunE: prepareRuntimeState,
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getRuntimeState(cmd)
		if err != nil {
			return err
		}
		toAdd := expandPatternArgs(args)
		if len(toAdd) == 0 {
			return nil
		}
		patterns, err := readPatterns(state.ignoreFile)
		if err != nil {
			return err
		}

		existing := make(map[string]struct{}, len(patterns))
		for _, p := range patterns {
			existing[p] = struct{}{}
		}

		for _, p := range toAdd {
			if _, ok := existing[p]; ok {
				continue
			}
			patterns = append(patterns, p)
			existing[p] = struct{}{}
		}

		return writePatterns(state.ignoreFile, patterns)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
