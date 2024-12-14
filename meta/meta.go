package meta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/berquerant/fflist/metric"
)

type Prober interface {
	Probe(ctx context.Context, path string) (*Data, error)
}

var (
	_ Prober = &FFProber{}
)

func NewProber(cmd string) *FFProber {
	return &FFProber{
		cmd: cmd,
	}
}

// FFProber reads file using ffprobe and returns a metadata.
type FFProber struct {
	cmd string
}

var (
	ErrProbe = errors.New("Probe")
)

func (p FFProber) Probe(ctx context.Context, path string) (*Data, error) {
	metric.IncrProbeCount()

	b, err := p.probe(ctx, path)
	if err != nil {
		metric.IncrProbeFailedCount()
		return nil, fmt.Errorf("%w: path %s", err, path)
	}

	d, err := p.formatData(b, path)
	if err != nil {
		metric.IncrProbeFailedCount()
		return nil, fmt.Errorf("%w: path %s", err, path)
	}

	metric.IncrProbeSuccessCount()
	return d, nil
}

func (p FFProber) probe(ctx context.Context, path string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, p.cmd,
		"-v", "error", // log level
		"-hide_banner",
		"-show_entries", "format", // display file format
		"-of", "json=c=1", // as compact json
		path,
	)
	x, err := cmd.Output()
	if err != nil {
		return nil, errors.Join(ErrProbe, err)
	}
	return x, nil
}

func (FFProber) formatData(b []byte, path string) (*Data, error) {
	d := map[string]any{}
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, errors.Join(ErrProbe, err)
	}

	formatRaw, ok := d["format"]
	if !ok {
		return nil, fmt.Errorf("%w: ffprobe: format is not found", ErrProbe)
	}
	format, ok := formatRaw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w: ffprobe: format is not an object", ErrProbe)
	}

	var (
		r = map[string]string{}
		w = func(key string, value any) {
			if v, exist := r[key]; exist {
				slog.Warn("FFProber duplicated metadata",
					slog.String("path", path),
					slog.String("key", key),
					slog.Any("value", v),
					slog.Any("newValue", value),
				)
			}
			r[key] = fmt.Sprint(value)
		}
	)

	for k, v := range format {
		switch k {
		case "tags":
			// flatten tags
			tags, ok := v.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("%w: ffprobe: tags is not an object", ErrProbe)
			}
			for tk, tv := range tags {
				w(tk, tv)
			}
		default:
			w(k, v)
		}
	}

	return NewData(r), nil
}
