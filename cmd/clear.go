package cmd

import "github.com/spf13/cobra"

var clearCmd = &cobra.Command{
	Use:     "clear",
	Short:   "Remove all patterns",
	Long:    "Remove all non-comment patterns from the selected ignore file.",
	Example: "pgi clear",
	Args:    cobra.NoArgs,
	PreRunE: prepareRuntimeState,
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getRuntimeState(cmd)
		if err != nil {
			return err
		}
		return writePatterns(state.ignoreFile, []string{})
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
