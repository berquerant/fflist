package walk

import (
	"context"
	"io/fs"
	"iter"
	"log/slog"
	"path/filepath"

	"github.com/berquerant/fflist/logx"
	"github.com/berquerant/fflist/metric"
)

type Walker interface {
	Walk(root string) iter.Seq[Entry]
	Err() error
}

//go:generate go run github.com/berquerant/dataclass -type Entry -field "Path string|Info fs.FileInfo" -output entry_dataclass_generated.go

var (
	_ Walker = &FileWalker{}
)

func New() *FileWalker { return &FileWalker{} }

// FileWalker walks only files under the root.
type FileWalker struct {
	err error
}

func (w FileWalker) Err() error { return w.err }

const (
	fileWalkerBufferSize = 100
)

func (w *FileWalker) Walk(root string) iter.Seq[Entry] {
	w.err = nil

	return func(yield func(Entry) bool) {
		resultC := make(chan Entry, fileWalkerBufferSize)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			defer close(resultC)

			_ = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
				slog.Debug("FileWalker", slog.String("path", path), logx.Err(err))
				metric.IncrEntryCount()

				select {
				case <-ctx.Done():
					return filepath.SkipAll
				default:
					if err != nil {
						w.err = err
						return err
					}
					if info.IsDir() {
						// skip dir
						return nil
					}
					resultC <- NewEntry(path, info)
					return nil
				}
			})
		}()

		for x := range resultC {
			if !yield(x) {
				return
			}
		}
	}
}
