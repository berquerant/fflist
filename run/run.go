package run

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/berquerant/fflist/info"
	"github.com/berquerant/fflist/logx"
	"github.com/berquerant/fflist/meta"
	"github.com/berquerant/fflist/worker"
)

func NewQuery(
	root []string,
	walkWorker *worker.Walker,
	probeWorker *worker.Prober,
	writer *Writer,
) *Query {
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

func NewIndexQuery(
	r io.Reader,
	writer *Writer,
) *IndexQuery {
	return &IndexQuery{
		r:      r,
		writer: writer,
	}
}

type IndexQuery struct {
	r      io.Reader
	writer *Writer
}

func (q *IndexQuery) Run(ctx context.Context) error {
	startTime := time.Now()

	scanner := bufio.NewScanner(q.r)
	for scanner.Scan() {
		d := map[string]string{}
		if err := json.Unmarshal(scanner.Bytes(), &d); err != nil {
			slog.Warn("IndexQuery", logx.Err(err))
			continue
		}
		md := info.New(meta.NewData(d))
		if err := q.writer.Write(ctx, md); err != nil {
			slog.Warn("IndexQuery", logx.Err(err))
		}
	}

	q.writer.WriteMetrics(time.Since(startTime))
	return nil
}
