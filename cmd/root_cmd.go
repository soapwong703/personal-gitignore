package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:           "pgi",
	Short:         "Manage per-repo and global gitignore patterns",
	Long:          "pgi manages per-repo and global gitignore patterns.\n\nYou can operate on the current repository's `.git/info/exclude` file (local)\nor the global core.excludesfile (global). By default the local scope is used.\nUse --global or -g to operate on the global excludesfile.",
	Example:       "\n  pgi list\n  pgi list \"*.log\"\n  pgi add \"*.log\"\n  pgi --global add \"*.env\"\n  pgi edit\n",
	SilenceUsage:  true,
	SilenceErrors: true,
}

var flagGlobal bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&flagGlobal, "global", "g", false, "use global core.excludesfile")
}
