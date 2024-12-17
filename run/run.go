package run

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/berquerant/fflist/info"
	"github.com/berquerant/fflist/logx"
	"github.com/berquerant/fflist/meta"
	"github.com/berquerant/fflist/metric"
	"github.com/berquerant/fflist/query"
	"github.com/berquerant/fflist/walk"
	"github.com/berquerant/fflist/worker"
)

func NewQuery(
	prober meta.Prober,
	selector query.Selector,
	newWalker func() walk.Walker,
	root []string,
	verbose bool,
	probeWorkerNum int,
) *Query {
	r := make([]string, len(root))
	for i := range root {
		r[i] = os.ExpandEnv(root[i])
	}

	walkWorker := worker.NewWalker(newWalker)
	probeWorker := worker.NewProbe(prober, probeWorkerNum)

	return &Query{
		selector: selector,
		root:     r,
		verbose:  verbose,

		walkWorker:  walkWorker,
		probeWorker: probeWorker,
	}
}

type Query struct {
	selector query.Selector
	root     []string
	verbose  bool

	walkWorker  *worker.Walker
	probeWorker *worker.Prober
}

func (q *Query) Run(ctx context.Context) error {
	startTime := time.Now()

	entryC := q.walkWorker.Start(ctx, q.root...)
	dataC := q.probeWorker.Start(ctx, entryC)

	for data := range dataC {
		if !q.selector.Select(ctx, data) {
			continue
		}
		if err := q.output(data); err != nil {
			path, _ := data.Get("path")
			slog.Error("Failed to output", slog.String("path", path), logx.Err(err))
		}
	}

	q.outputMetrics(time.Since(startTime))
	return q.walkWorker.Err()
}

func (q *Query) output(data info.Getter) error {
	metric.IncrAcceptCount()

	if q.verbose {
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		_, err = fmt.Printf("%s\n", b)
		return err
	}

	path, _ := data.Get("path")
	fmt.Println(path)
	return nil
}

func (q *Query) outputMetrics(duration time.Duration) {
	if !q.verbose {
		return
	}

	fmt.Fprintf(os.Stderr, "%s\n", logx.Jsonify(map[string]any{
		"Duration": duration.Seconds(),
		"Metrics":  metric.Get(),
	}))
}
