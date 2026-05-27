package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list [glob]",
	Short:   "Show the current ignore patterns",
	Long:    "List patterns in the selected ignore file, optionally filtered by a glob.",
	Example: "pgi list \"*.log\"",
	Args:    cobra.MaximumNArgs(1),
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
		visible := make([]string, 0, len(patterns))
		for _, pattern := range patterns {
			if isCommentLine(pattern) {
				continue
			}
			visible = append(visible, pattern)
		}
		if len(args) == 1 {
			visible, err = filterPatternsByGlob(visible, args[0])
			if err != nil {
				return err
			}
		}
		out := cmd.OutOrStdout()
		for _, p := range visible {
			fmt.Fprintln(out, p)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
