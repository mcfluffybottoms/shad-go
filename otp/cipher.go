//go:build !solution

package otp

import (
	"io"
)

type cipherReader struct {
	r    io.Reader
	prng io.Reader
	buf  []byte
}

func (cr *cipherReader) Read(p []byte) (n int, err error) {
	n, err = cr.r.Read(p)
	if n == 0 {
		return 0, err
	}

	if len(cr.buf) < n {
		cr.buf = make([]byte, n)
	}
	key := cr.buf[:n]

	if _, err := io.ReadFull(cr.prng, key); err != nil {
		return 0, err
	}

	for i := 0; i < n; i++ {
		p[i] ^= key[i]
	}

	return n, err
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &cipherReader{
		r:    r,
		prng: prng,
		buf:  make([]byte, 0),
	}
}

type cipherWriter struct {
	w    io.Writer
	prng io.Reader
	buf  []byte
}

func (cw *cipherWriter) Write(p []byte) (n int, err error) {
	if len(cw.buf) != len(p) {
		cw.buf = make([]byte, len(p))
	}
	key := cw.buf[:len(p)]

	if _, err := io.ReadFull(cw.prng, key); err != nil {
		return 0, err
	}

	enc := make([]byte, len(p))
	for i := range p {
		enc[i] = p[i] ^ key[i]
	}

	n, err = cw.w.Write(enc)
	return n, err
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &cipherWriter{
		w:    w,
		prng: prng,
		buf:  make([]byte, 0),
	}
}
