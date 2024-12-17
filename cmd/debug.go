package main

import (
	"fmt"
	"os"

	"github.com/berquerant/fflist/logx"
	"github.com/berquerant/fflist/meta"
	"github.com/berquerant/fflist/worker"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(debugCmd)
}

var debugCmd = &cobra.Command{
	Use:   "debug [DIR...]",
	Short: `Walk directories and show metadata`,
	Long: `Walk directories and show metadata.

# walk ~/Music
fflint debug ~/Music
# read paths from stdin
fflint debug -`,
	RunE: func(cmd *cobra.Command, args []string) error {
		roots := []string{"."}
		if len(args) > 0 {
			roots = args
		}

		newWalker, err := newWalkerFactory(roots)
		if err != nil {
			return err
		}

		var (
			walker = newWalker()
			prober = meta.NewProber(getProbe(cmd))
		)

		for _, root := range roots {
			for x := range walker.Walk(root) {
				e := logx.Jsonify(worker.BuildInfoGetter(cmd.Context(), prober, x))
				fmt.Printf("%s\n", e)
			}
			if err := walker.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		}
		return nil
	},
}
