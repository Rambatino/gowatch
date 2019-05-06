package cmd

import (
	"../app"
	"github.com/spf13/cobra"
)

var (
	Extensions []string
	Paths      []string
	Files      []string
	Recursive  bool
	Ignore     []string
)

func init() {
	runCmd.Flags().StringSliceVarP(&Extensions, "extensions", "e", []string{}, "Comma separated file extensions/types in which to search for. N.B. don't pass globs.")
	runCmd.Flags().StringSliceVarP(&Paths, "paths", "p", []string{}, "Comma separated paths (folders and files) in which to search in. N.B. don't pass globs.")
	runCmd.Flags().BoolVarP(&Recursive, "recursive", "r", false, "Whether to search recursively")
	runCmd.Flags().StringSliceVarP(&Ignore, "ignore", "i", []string{}, "What file paths to ignore. N.B. don't pass globs.")
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run custom function",
	Long:  `Runs your command with parameters`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		watcher, _ := app.NewWatcher(Extensions, Paths, Recursive, Ignore, args)
		<-watcher.WatchAndRun()
	},
}
