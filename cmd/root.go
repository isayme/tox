package cmd

import (
	"fmt"
	"os"

	"github.com/isayme/go-logger"
	"github.com/isayme/tox/util"
	"github.com/spf13/cobra"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
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
		if enableProfilingFlag {
			logger.Info("profiling enabled")

			err := profiler.Start(
				profiler.WithVersion(util.Version),
				profiler.WithProfileTypes(
					profiler.CPUProfile,
					profiler.HeapProfile,
					profiler.BlockProfile,
					profiler.GoroutineProfile,
					profiler.MutexProfile,
				),
			)
			if err != nil {
				logger.Panic(err)
			}
			defer profiler.Stop()
		}

		startLocal()
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Run: func(cmd *cobra.Command, args []string) {
		if enableProfilingFlag {
			logger.Info("profiling enabled")

			err := profiler.Start(
				profiler.WithVersion(util.Version),
				profiler.WithProfileTypes(
					profiler.CPUProfile,
					profiler.HeapProfile,
					profiler.BlockProfile,
					profiler.GoroutineProfile,
					profiler.MutexProfile,
				),
			)
			if err != nil {
				logger.Panic(err)
			}
			defer profiler.Stop()
		}

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
