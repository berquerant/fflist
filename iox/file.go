package iox

import "os"

func Open(file ...string) ([]*os.File, error) {
	fs := make([]*os.File, len(file))

	for i, x := range file {
		f, err := os.Open(x)
		if err != nil {
			for j := 0; j < i; j++ {
				fs[j].Close()
			}
			return nil, err
		}
		fs[i] = f
	}

	return fs, nil
}
