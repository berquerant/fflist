package run

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/berquerant/fflist/info"
	"github.com/berquerant/fflist/logx"
	"github.com/berquerant/fflist/meta"
	"github.com/berquerant/fflist/metric"
	"github.com/berquerant/fflist/query"
	"github.com/berquerant/fflist/walk"
	"golang.org/x/sync/errgroup"
)

const (
	queryWorkerBufferSize = 100
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
	return &Query{
		selector:       selector,
		newWalker:      newWalker,
		root:           r,
		verbose:        verbose,
		prober:         prober,
		probeWorkerNum: probeWorkerNum,
	}
}

type Query struct {
	selector       query.Selector
	newWalker      func() walk.Walker
	root           []string
	verbose        bool
	prober         meta.Prober
	probeWorkerNum int

	sync.Mutex
	err error
}

func (q *Query) Run(ctx context.Context) error {
	startTime := time.Now()

	entryC := q.walkWorker(ctx)
	dataC := q.probeWorker(ctx, entryC)

	for data := range dataC {
		if !q.selector.Select(data) {
			continue
		}
		if err := q.output(data); err != nil {
			path, _ := data.Get("path")
			slog.Error("Failed to output", slog.String("path", path), logx.Err(err))
		}
	}

	if err := q.outputMetrics(time.Since(startTime)); err != nil {
		slog.Error("Failed to output metrics", logx.Err(err))
	}

	q.Lock()
	defer q.Unlock()
	return q.err
}

func (q *Query) setErr(err error) {
	q.Lock()
	defer q.Unlock()
	q.err = err
}

func (q *Query) walkWorker(ctx context.Context) <-chan walk.Entry {
	eg, ctx := errgroup.WithContext(ctx)
	var (
		entryC = make(chan walk.Entry, queryWorkerBufferSize)
	)

	for _, root := range q.root {
		eg.Go(func() error {
			walker := q.newWalker()
			for entry := range walker.Walk(root) {
				select {
				case <-ctx.Done():
					break
				default:
					entryC <- entry
				}
			}
			return walker.Err()
		})
	}

	go func() {
		if err := eg.Wait(); err != nil {
			q.setErr(err)
		}
		close(entryC)
	}()

	return entryC
}

func (q *Query) probeWorker(ctx context.Context, entryC <-chan walk.Entry) <-chan info.Getter {
	var (
		wg      sync.WaitGroup
		resultC = make(chan info.Getter, queryWorkerBufferSize)
	)

	for range q.probeWorkerNum {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for entry := range entryC {
				resultC <- q.buildGetter(ctx, entry)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultC)
	}()

	return resultC
}

func (q *Query) buildGetter(ctx context.Context, entry walk.Entry) info.Getter {
	return BuildInfoGetter(ctx, q.prober, entry)
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
	_, err := fmt.Println(path)
	return err
}

func (q *Query) outputMetrics(duration time.Duration) error {
	if !q.verbose {
		return nil
	}

	_, err := fmt.Fprintf(os.Stderr, "%s\n", logx.Jsonify(map[string]any{
		"Duration": duration.Seconds(),
		"Metrics":  metric.Get(),
	}))
	return err
}
