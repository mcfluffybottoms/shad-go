//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"io"
	"os"
	"sort"
)

type item struct {
	line     string
	readerId int
}

type lineHeap []item

func (h lineHeap) Len() int {
	return len(h)
}

func (h lineHeap) Less(i, j int) bool {
	return h[i].line < h[j].line
}

func (h lineHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *lineHeap) Push(x any) {
	*h = append(*h, x.(item))
}

func (h *lineHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

type Reader struct {
	sc *bufio.Reader
}

func (r *Reader) ReadLine() (string, error) {
	line, err := r.sc.ReadString('\n')
	if err == io.EOF {
		if len(line) > 0 {
			return line, nil
		} else {
			return "", io.EOF
		}
	}

	if err != nil {
		return "", err
	}

	return line[:len(line)-1], nil
}

func NewReader(r io.Reader) LineReader {
	return &Reader{sc: bufio.NewReader(r)}
}

type Writer struct {
	w io.Writer
}

func (w *Writer) Write(l string) error {
	_, err := io.WriteString(w.w, l+"\n")
	return err
}

func NewWriter(w io.Writer) LineWriter {
	return &Writer{w: w}
}

func readLineToHeap(h *lineHeap, r LineReader, readerId int) error {
	line, err := r.ReadLine()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	heap.Push(h, item{readerId: readerId, line: line})
	return nil
}

func Merge(w LineWriter, readers ...LineReader) error {
	h := &lineHeap{}
	heap.Init(h)

	for i, r := range readers {
		if err := readLineToHeap(h, r, i); err != nil {
			return err
		}
	}

	for h.Len() > 0 {
		it := heap.Pop(h).(item)
		if err := w.Write(it.line); err != nil {
			return err
		}
		if err := readLineToHeap(h, readers[it.readerId], it.readerId); err != nil {
			return err
		}
	}

	return nil
}

func Sort(w io.Writer, in ...string) error {
	readers := make([]LineReader, len(in))
	for i, p := range in {
		// read and sort
		file1, err := os.Open(p)
		if err != nil {
			return err
		}
		lr := NewReader(file1)
		lines := make([]string, 0)
		for {
			line, err := lr.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				file1.Close()
				return err
			}
			lines = append(lines, line)
		}
		if err := file1.Close(); err != nil {
			return err
		}
		sort.Strings(lines)

		// write to file
		file2, err := os.Create(p)
		if err != nil {
			return err
		}

		lw := NewWriter(file2)
		for _, line := range lines {
			if err := lw.Write(line); err != nil {
				file2.Close()
				return err
			}
		}
		if err := file2.Close(); err != nil {
			return err
		}
		// add reader to readers
		file3, err := os.Open(p)
		if err != nil {
			return err
		}
		defer file3.Close()
		readers[i] = NewReader(file3)
	}
	return Merge(NewWriter(w), readers...)
}
