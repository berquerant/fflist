package slicesx_test

import (
	"testing"

	"github.com/berquerant/fflist/slicesx"
	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	for _, tc := range []struct {
		title string
		s     []string
		p     []string
		want  [][]string
	}{
		{
			title: "pivot hit last",
			s:     []string{"1", "2", "|"},
			p:     []string{"|"},
			want: [][]string{
				{"1", "2"},
				{},
			},
		},
		{
			title: "pivot hit first",
			s:     []string{"|", "1", "2"},
			p:     []string{"|"},
			want: [][]string{
				{},
				{"1", "2"},
			},
		},
		{
			title: "pivot hit",
			s:     []string{"1", "|", "2"},
			p:     []string{"|"},
			want: [][]string{
				{"1"},
				{"2"},
			},
		},
		{
			title: "pivot not hit",
			s:     []string{"1", "|", "2"},
			p:     []string{"&"},
			want: [][]string{
				{"1", "|", "2"},
			},
		},
		{
			title: "empty pivot",
			s:     []string{"1", "2"},
			want: [][]string{
				{"1", "2"},
			},
		},
		{
			title: "empty s",
			p:     []string{"|"},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got := slicesx.Chunk(tc.s, tc.p...)
			assert.Equal(t, tc.want, got)
		})
	}
}
