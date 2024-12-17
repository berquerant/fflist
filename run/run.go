package run

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/berquerant/fflist/logx"
	"github.com/berquerant/fflist/meta"
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
	walkWorker := worker.NewWalker(newWalker)
	probeWorker := worker.NewProbe(prober, probeWorkerNum)
	writer := NewWriter(os.Stdout, selector, verbose)

	return &Query{
		root:        ExpandEnvAll(root...),
		walkWorker:  walkWorker,
		probeWorker: probeWorker,
		writer:      writer,
	}
}

type Query struct {
	root        []string
	walkWorker  *worker.Walker
	probeWorker *worker.Prober
	writer      *Writer
}

func (q *Query) Run(ctx context.Context) error {
	startTime := time.Now()

	entryC := q.walkWorker.Start(ctx, q.root...)
	dataC := q.probeWorker.Start(ctx, entryC)

	for data := range dataC {
		if err := q.writer.Write(ctx, data); err != nil {
			path, _ := data.Get("path")
			slog.Error("Failed to output", slog.String("path", path), logx.Err(err))
		}
	}

	q.writer.WriteMetrics(time.Since(startTime))
	return q.walkWorker.Err()
}
