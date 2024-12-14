package slicesx

// Chunk divides the slice using elements that match the pivot as boundaries.
// This matching element is not included in the resulting segments.
func Chunk[S ~[]E, E comparable](s S, pivot ...E) []S {
	var (
		r   []S
		acc = []E{}
		p   = map[E]bool{}
	)
	for _, x := range pivot {
		p[x] = true
	}

	for i, x := range s {
		if p[x] {
			r = append(r, acc)
			acc = []E{}

			if i == len(s)-1 {
				r = append(r, acc)
			}
			continue
		}
		acc = append(acc, x)
	}
	if len(acc) > 0 {
		r = append(r, acc)
	}

	return r
}
