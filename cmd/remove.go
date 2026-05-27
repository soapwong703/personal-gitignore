package cmd

import "github.com/spf13/cobra"

var removeCmd = &cobra.Command{
	Use:     "remove <pattern>",
	Short:   "Remove a pattern",
	Long:    "Remove matching pattern lines from the selected ignore file.",
	Example: "pgi remove \"*.log\"",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := buildCommandContext()
		if err != nil {
			return err
		}
		ignoreFile, err := resolveIgnoreFile(ctx)
		if err != nil {
			return err
		}
		if err := ensureFile(ignoreFile); err != nil {
			return err
		}
		patterns, err := readPatterns(ignoreFile)
		if err != nil {
			return err
		}
		updated := make([]string, 0, len(patterns))
		for _, p := range patterns {
			if p != args[0] {
				updated = append(updated, p)
			}
		}
		return writePatterns(ignoreFile, updated)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
