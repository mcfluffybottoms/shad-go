//go:build !solution

package once

type Once struct {
	ch chan struct{}
}

func New() *Once {
	return &Once{
		ch: make(chan struct{}, 1),
	}
}

func (o *Once) Do(f func()) {
	select {
	case o.ch <- struct{}{}:
		f()
	default:
	}
}
