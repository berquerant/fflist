package run

import (
	"context"
	"errors"
	"log/slog"

	"github.com/berquerant/fflist/info"
	"github.com/berquerant/fflist/logx"
	"github.com/berquerant/fflist/meta"
	"github.com/berquerant/fflist/walk"
)

func BuildInfoGetter(ctx context.Context, prober meta.Prober, e walk.Entry) info.Getter {
	r := []*meta.Data{
		info.NewMetadataFromEntry(e),
	}

	data, err := prober.Probe(ctx, e.Path())
	if err != nil && !errors.Is(err, context.Canceled) {
		slog.Warn("Failed to probe", logx.Err(err))
	} else {
		r = append(r, data)
	}

	return info.New(r...)
}
