package main

import (
	"fmt"
	"os"

	"github.com/berquerant/fflist/logx"
	"github.com/berquerant/fflist/meta"
	"github.com/berquerant/fflist/run"
	"github.com/berquerant/fflist/walk"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(debugCmd)
}

var debugCmd = &cobra.Command{
	Use:   "debug [DIR...]",
	Short: `Walk directories and show metadata`,
	RunE: func(cmd *cobra.Command, args []string) error {
		roots := []string{"."}
		if len(args) > 0 {
			roots = args
		}

		var (
			walker = walk.New()
			prober = meta.NewProber(getProbe(cmd))
		)

		for _, root := range roots {
			for x := range walker.Walk(root) {
				e := logx.Jsonify(run.BuildInfoGetter(cmd.Context(), prober, x))
				fmt.Printf("%s\n", e)
			}
			if err := walker.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		}
		return nil
	},
}
