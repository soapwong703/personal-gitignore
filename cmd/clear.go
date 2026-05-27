package cmd

import "github.com/spf13/cobra"

var clearCmd = &cobra.Command{
	Use:     "clear",
	Short:   "Remove all patterns",
	Long:    "Remove all non-comment patterns from the selected ignore file.",
	Example: "pgi clear",
	Args:    cobra.NoArgs,
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
		return writePatterns(ignoreFile, []string{})
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
