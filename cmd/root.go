package cmd

import (
	"fmt"
	"os"

	"github.com/isayme/tox/util"
	"github.com/spf13/cobra"
)

var (
	versionFlag         bool
	enableProfilingFlag bool
)

func init() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "show version")
	rootCmd.AddCommand(localCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.PersistentFlags().BoolVarP(&enableProfilingFlag, "profiling", "", false, "enable profiling")
}

var rootCmd = &cobra.Command{
	Use: "tox",
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			util.PrintVersion()
			os.Exit(0)
		}
	},
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "run local",
	Run: func(cmd *cobra.Command, args []string) {
		util.EnableProfiling(enableProfilingFlag)
		startLocal()
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Run: func(cmd *cobra.Command, args []string) {
		util.EnableProfiling(enableProfilingFlag)
		startServer()
	},
}

// Execute run root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
