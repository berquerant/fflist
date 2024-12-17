package walk

import (
	"bufio"
	"context"
	"io"
	"io/fs"
	"iter"
	"log/slog"
	"os"
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

func NewFile() *FileWalker { return &FileWalker{} }

// FileWalker walks only files under the root.
type FileWalker struct {
	err error
}

func (w FileWalker) Err() error { return w.err }

const (
	walkerBufferSize = 100
)

func (w *FileWalker) Walk(root string) iter.Seq[Entry] {
	w.err = nil

	return func(yield func(Entry) bool) {
		resultC := make(chan Entry, walkerBufferSize)
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

var (
	_ Walker = &ReaderWalker{}
)

type ReaderWalker struct {
	r          io.Reader
	fileWalker Walker
	err        error
}

func NewReader(r io.Reader, fileWalker Walker) *ReaderWalker {
	return &ReaderWalker{
		r:          r,
		fileWalker: fileWalker,
	}
}

func (w ReaderWalker) Err() error { return w.err }

func (w *ReaderWalker) Walk(_ string) iter.Seq[Entry] {
	w.err = nil

	return func(yield func(Entry) bool) {
		resultC := make(chan Entry, walkerBufferSize)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			defer close(resultC)

			select {
			case <-ctx.Done():
				return
			default:
				scanner := bufio.NewScanner(w.r)
				for scanner.Scan() {
					select {
					case <-ctx.Done():
						return
					default:
						path := scanner.Text()
						info, err := os.Stat(path)
						if os.IsNotExist(err) {
							slog.Debug("ReaderWalker", slog.String("path", path), logx.Err(err))
							continue
						}

						slog.Debug("ReaderWalker", slog.String("path", path), logx.Err(err))
						metric.IncrEntryCount()

						if err != nil {
							w.err = err
							return
						}
						if info.IsDir() {
							for x := range w.fileWalker.Walk(path) {
								resultC <- x
							}
							if err := w.fileWalker.Err(); err != nil {
								slog.Warn("ReaderWalker", slog.String("path", path), logx.Err(err))
							}
							continue
						}
						resultC <- NewEntry(path, info)
					}
				}

				if err := scanner.Err(); err != nil {
					w.err = err
				}
			}
		}()

		for x := range resultC {
			if !yield(x) {
				return
			}
		}
	}
}
