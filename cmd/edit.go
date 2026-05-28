package cmd

import "github.com/spf13/cobra"

var editCmd = &cobra.Command{
	Use:     "edit",
	Short:   "Open the ignore file in your editor",
	Long:    "Open the selected ignore file in the editor configured by $EDITOR, $VISUAL, or git var GIT_EDITOR.",
	Example: "  pgi edit\n  pgi edit --editor \"code --wait\"",
	Args:    cobra.NoArgs,
	PreRunE: prepareRuntimeState,
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getRuntimeState(cmd)
		if err != nil {
			return err
		}
		editor, _ := cmd.Flags().GetString("editor")
		return openEditor(state.ignoreFile, state.ctx.env, editor)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringP("editor", "e", "", "Editor command to use (overrides $EDITOR/$VISUAL)")
}
