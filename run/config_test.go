package run_test

import (
	"bytes"
	"testing"

	"github.com/berquerant/fflist/run"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	for _, tc := range []struct {
		title string
		src   string
		want  *run.Config
		err   error
	}{
		{
			title: "yaml",
			src: `{
  "root": ["ROOT"],
  "query": [["name=NAME"]]
}`,
			want: &run.Config{
				Root: []string{
					"ROOT",
				},
				Query: [][]string{
					{"name=NAME"},
				},
			},
		},
		{
			title: "yaml",
			src: `root:
- ROOT
query:
- - name=NAME`,
			want: &run.Config{
				Root: []string{
					"ROOT",
				},
				Query: [][]string{
					{"name=NAME"},
				},
			},
		},
		{
			title: "empty query",
			src: `root:
- ROOT
query:
- - name=NAME
-`,
			err: run.ErrConfig,
		},
		{
			title: "no query",
			src: `root:
- ROOT`,
			err: run.ErrConfig,
		},
		{
			title: "no root",
			src: `query:
- - name=NAME`,
			err: run.ErrConfig,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			r := bytes.NewBufferString(tc.src)
			got, err := run.ParseConfig(r)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
