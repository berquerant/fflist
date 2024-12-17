package worker

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/berquerant/fflist/info"
	"github.com/berquerant/fflist/logx"
	"github.com/berquerant/fflist/meta"
	"github.com/berquerant/fflist/walk"
)

const (
	probeWorkerBufferSize = 100
)

type Prober struct {
	prober    meta.Prober
	workerNum int
}

func NewProbe(prober meta.Prober, workerNum int) *Prober {
	if workerNum < 1 {
		workerNum = 1
	}
	return &Prober{
		prober:    prober,
		workerNum: workerNum,
	}
}

func (w *Prober) Start(ctx context.Context, entryC <-chan walk.Entry) <-chan info.Getter {
	var (
		wg      sync.WaitGroup
		resultC = make(chan info.Getter, probeWorkerBufferSize)
	)

	for i := range w.workerNum {
		wg.Add(1)
		go func() {
			slog.Debug("Prober Start", slog.Int("n", i))
			defer wg.Done()

			for entry := range entryC {
				resultC <- w.buildInfoGetter(ctx, entry)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultC)
		slog.Debug("Prober Stop")
	}()

	return resultC
}

func (w *Prober) buildInfoGetter(ctx context.Context, entry walk.Entry) info.Getter {
	return BuildInfoGetter(ctx, w.prober, entry)
}

func BuildInfoGetter(ctx context.Context, prober meta.Prober, entry walk.Entry) info.Getter {
	r := []*meta.Data{
		info.NewMetadataFromEntry(entry),
	}

	data, err := prober.Probe(ctx, entry.Path())
	if err != nil && !errors.Is(err, context.Canceled) {
		slog.Warn("Failed to probe", logx.Err(err))
	} else {
		r = append(r, data)
	}

	return info.New(r...)
}
