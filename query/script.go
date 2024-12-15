package query

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"

	"github.com/berquerant/execx"
	"github.com/berquerant/fflist/info"
	"github.com/berquerant/fflist/logx"
	"github.com/berquerant/fflist/metric"
)

// ScriptSelector returns true if script exit with 0.
// The stdin of script is the json of metadata.
type ScriptSelector struct {
	shell  string
	script *execx.Script
}

func NewScriptSelector(q Query) *ScriptSelector {
	s := execx.NewScript(q.Value(), q.Key())
	s.KeepScriptFile = true
	s.Env = execx.EnvFromEnviron()
	return &ScriptSelector{
		shell:  q.Key(),
		script: s,
	}
}

func (s *ScriptSelector) Select(ctx context.Context, data info.Getter) bool {
	metric.IncrSelectCount()

	err := s.script.Runner(func(cmd *execx.Cmd) error {
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		cmd.Stdin = bytes.NewBuffer(b)
		_, err = cmd.Run(ctx)
		return err
	})

	r := err == nil
	slog.Debug("ScriptSelector", logx.Err(err))

	if r {
		metric.IncrSelectSuccessCount()
	} else {
		metric.IncrSelectFailedCount()
	}
	return r
}
