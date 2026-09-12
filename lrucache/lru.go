//go:build !solution

package lrucache

import "container/list"

type item struct {
	key   int
	value int
}

type Lrucache struct {
	capacity int
	time     *list.List
	storage  map[int]*list.Element
}

// Clear implements [Cache].
func (l *Lrucache) Clear() {
	l.time.Init()
	for k := range l.storage {
		delete(l.storage, k)
	}
}

// Get implements [Cache].
func (l *Lrucache) Get(key int) (int, bool) {
	val, ok := l.storage[key]
	if !ok {
		return 0, false
	}
	l.time.MoveToFront(val)
	return val.Value.(*item).value, true
}

// Range implements [Cache].
func (l *Lrucache) Range(f func(key int, value int) bool) {
	for i := l.time.Back(); i != nil; i = i.Prev() {
		// l.time.MoveToFront(i)
		it := i.Value.(*item)
		if !f(it.key, it.value) {
			return
		}
	}
}

// Set implements [Cache].
func (l *Lrucache) Set(key int, value int) {
	if l.capacity == 0 {
		return
	}
	if el, ok := l.storage[key]; ok {
		el.Value.(*item).value = value
		l.time.MoveToFront(el)
		return
	}
	if len(l.storage) >= l.capacity {
		back := l.time.Back()
		l.time.Remove(back)
		delete(l.storage, back.Value.(*item).key)
	}
	node := l.time.PushFront(&item{key: key, value: value})
	l.storage[key] = node
}

func New(cap int) Cache {
	return &Lrucache{
		capacity: cap,
		storage:  make(map[int]*list.Element),
		time:     list.New(),
	}
}
