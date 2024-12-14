package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/berquerant/fflist/logx"
	"github.com/berquerant/fflist/run"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logs")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet logs")
	rootCmd.PersistentFlags().StringP("probe", "p", "ffprobe", "Media analyzer command")
}

func getProbe(cmd *cobra.Command) string {
	x, _ := cmd.Flags().GetString("probe")
	return x
}

var rootCmd = &cobra.Command{
	Use:   "fflist",
	Short: `Select media file resources`,
	Long: `Select media file resources.

Requirements:
- ffprobe 7.1 https://ffmpeg.org/ffprobe.html`,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		logLevel := slog.LevelInfo
		if debugEnabled, _ := cmd.Flags().GetBool("debug"); debugEnabled {
			logLevel = slog.LevelDebug
		}
		if quiet, _ := cmd.Flags().GetBool("quiet"); quiet {
			logLevel = slog.LevelError
		}
		logx.Setup(os.Stderr, logLevel)
	},
}

func rootFlag(cmd *cobra.Command) {
	cmd.Flags().StringSliceP("root", "r", []string{"."}, "Root directories")
}

func getRoot(cmd *cobra.Command) []string {
	x, _ := cmd.Flags().GetStringSlice("root")
	return x
}

func verboseFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}

func getVerbose(cmd *cobra.Command) bool {
	x, _ := cmd.Flags().GetBool("verbose")
	return x
}

func probeWorkerNumFlag(cmd *cobra.Command) {
	cmd.Flags().IntP("worker", "w", 8, "Probe worker num")
}

func getProbeWorkerNum(cmd *cobra.Command) int {
	x, _ := cmd.Flags().GetInt("worker")
	if x < 1 {
		return 1
	}
	return x
}

var (
	errNoConfig = errors.New("NoConfig")
)

func configFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("config", "c", "", "Query config file")
}

func getConfig(cmd *cobra.Command) (*run.Config, error) {
	file, _ := cmd.Flags().GetString("config")
	if file == "" {
		return nil, errNoConfig
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return run.ParseConfig(f)
}
