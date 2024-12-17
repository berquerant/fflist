package worker

import (
	"context"
	"log/slog"

	"github.com/berquerant/fflist/walk"
	"golang.org/x/sync/errgroup"
)

const (
	walkWorkerBufferSize = 100
)

type Walker struct {
	newWalker func() walk.Walker
	err       error
}

func NewWalker(newWalker func() walk.Walker) *Walker {
	return &Walker{
		newWalker: newWalker,
	}
}

func (w Walker) Err() error { return w.err }

func (w *Walker) Start(ctx context.Context, root ...string) <-chan walk.Entry {
	w.err = nil
	eg, ctx := errgroup.WithContext(ctx)
	entryC := make(chan walk.Entry, walkWorkerBufferSize)

	for i, r := range root {
		eg.Go(func() error {
			slog.Debug("Walker Start", slog.Int("n", i), slog.String("root", r))

			walker := w.newWalker()
			for entry := range walker.Walk(r) {
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
			w.err = err
		}
		close(entryC)
		slog.Debug("Walker Stop")
	}()

	return entryC
}
