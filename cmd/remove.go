package cmd

import "github.com/spf13/cobra"

var removeCmd = &cobra.Command{
	Use:     "remove <pattern...>",
	Short:   "Remove a pattern",
	Long:    "Remove matching pattern lines from the selected ignore file for one or more patterns.",
	Example: "pgi remove \"*.log\"",
	Args:    cobra.MinimumNArgs(1),
	PreRunE: prepareRuntimeState,
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getRuntimeState(cmd)
		if err != nil {
			return err
		}
		toRemove := expandPatternArgs(args)
		if len(toRemove) == 0 {
			return nil
		}

		patterns, err := readPatterns(state.ignoreFile)
		if err != nil {
			return err
		}

		removeSet := make(map[string]struct{}, len(toRemove))
		for _, p := range toRemove {
			removeSet[p] = struct{}{}
		}

		updated := make([]string, 0, len(patterns))
		for _, p := range patterns {
			if _, ok := removeSet[p]; !ok {
				updated = append(updated, p)
			}
		}
		return writePatterns(state.ignoreFile, updated)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
