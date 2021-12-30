package cmd

import (
	"fmt"
	"os"

	"github.com/isayme/go-toh2/util"
	"github.com/spf13/cobra"
)

var versionFlag bool

func init() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "show version")
	rootCmd.AddCommand(localCmd)
	rootCmd.AddCommand(serverCmd)
}

var rootCmd = &cobra.Command{
	Use: "toh2",
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
		startLocal()
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Run: func(cmd *cobra.Command, args []string) {
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
