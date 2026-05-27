package cmd

import "github.com/spf13/cobra"

var editCmd = &cobra.Command{
	Use:     "edit",
	Short:   "Open the ignore file in your editor",
	Long:    "Open the selected ignore file in the editor configured by $EDITOR, $VISUAL, or git var GIT_EDITOR.",
	Example: "pgi edit",
	Args:    cobra.NoArgs,
	PreRunE: prepareRuntimeState,
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getRuntimeState(cmd)
		if err != nil {
			return err
		}
		return openEditor(state.ignoreFile, state.ctx.env)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
