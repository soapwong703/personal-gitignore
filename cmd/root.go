package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:                   "pgi [--global] [--help] <command> [pattern]",
	Short:                 "Manage per-repo and global gitignore patterns",
	Long:                  "pgi manages per-repo and global gitignore patterns.\n\nYou can operate on the current repository's `.git/info/exclude` file (local)\nor the global core.excludesfile (global). By default the local scope is used.\nUse --global or -g to operate on the global excludesfile.",
	Example:               "\n  pgi list\n  pgi list \"*.log\"\n  pgi add \"*.log\"\n  pgi --global add \"*.env\"\n  pgi edit\n",
	DisableFlagsInUseLine: true,
	TraverseChildren:      true,
	SilenceUsage:          true,
	SilenceErrors:         true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

var flagGlobal bool

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		err = normalizeCLIError(err)
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func normalizeCLIError(err error) error {
	if strings.HasPrefix(err.Error(), "unknown flag:") {
		flag := strings.TrimSpace(strings.TrimPrefix(err.Error(), "unknown flag:"))
		return fmt.Errorf("unknown argument: %s", flag)
	}
	return err
}

func init() {
	flagErr := func(_ *cobra.Command, err error) error {
		return normalizeCLIError(err)
	}

	rootCmd.SetFlagErrorFunc(flagErr)
	for _, child := range rootCmd.Commands() {
		child.SetFlagErrorFunc(flagErr)
	}

	rootCmd.PersistentFlags().BoolVarP(&flagGlobal, "global", "g", false, "use global core.excludesfile")
}
