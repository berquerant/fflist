package run

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/berquerant/fflist/info"
	"github.com/berquerant/fflist/logx"
	"github.com/berquerant/fflist/metric"
	"github.com/berquerant/fflist/query"
)

type Writer struct {
	w        io.Writer
	selector query.Selector
	verbose  bool
}

func NewWriter(w io.Writer, selector query.Selector, verbose bool) *Writer {
	return &Writer{
		w:        w,
		selector: selector,
		verbose:  verbose,
	}
}

func (w *Writer) Write(ctx context.Context, data info.Getter) error {
	if !w.selector.Select(ctx, data) {
		return nil
	}
	return w.write(data)
}

func (w *Writer) write(data info.Getter) error {
	metric.IncrAcceptCount()

	if w.verbose {
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w.w, "%s\n", b)
		return err
	}

	path, _ := data.Get("path")
	fmt.Fprintln(w.w, path)
	return nil
}

func (w *Writer) WriteMetrics(duration time.Duration) {
	if !w.verbose {
		return
	}

	fmt.Fprintf(os.Stderr, "%s\n", logx.Jsonify(map[string]any{
		"Duration": duration.Seconds(),
		"Metrics":  metric.Get(),
	}))
}
