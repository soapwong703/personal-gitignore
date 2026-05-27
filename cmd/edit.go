package cmd

import "github.com/spf13/cobra"

var editCmd = &cobra.Command{
	Use:     "edit",
	Short:   "Open the ignore file in your editor",
	Long:    "Open the selected ignore file in the editor configured by $EDITOR, $VISUAL, or git var GIT_EDITOR.",
	Example: "pgi edit",
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
		return openEditor(ignoreFile, ctx.env)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
