//go:build !solution

package rwmutex

// A RWMutex is a reader/writer mutual exclusion lock.
// The lock can be held by an arbitrary number of readers or a single writer.
// The zero value for a RWMutex is an unlocked mutex.
//
// If a goroutine holds a RWMutex for reading and another goroutine might
// call Lock, no goroutine should expect to be able to acquire a read lock
// until the initial read lock is released. In particular, this prohibits
// recursive read locking. This is to ensure that the lock eventually becomes
// available; a blocked Lock call excludes new readers from acquiring the
// lock.

const writeBit = 1 << 31

type RWMutex struct {
	mu             chan struct{}
	state          uint32
	wakeAll        chan struct{}
	waitingWriters int
}

// New creates *RWMutex.
func New() *RWMutex {
	rw := &RWMutex{
		mu:      make(chan struct{}, 1),
		wakeAll: make(chan struct{}),
	}
	rw.mu <- struct{}{}
	return rw
}

// RLock locks rw for reading.
//
// It should not be used for recursive read locking; a blocked Lock
// call excludes new readers from acquiring the lock. See the
// documentation on the RWMutex type.
func (rw *RWMutex) RLock() {
	for {
		<-rw.mu
		if rw.state&writeBit != 0 {
			ch := rw.wakeAll
			rw.mu <- struct{}{}
			<-ch
			continue
		}

		rw.state++
		rw.mu <- struct{}{}
		return
	}
}

// RUnlock undoes a single RLock call;
// it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading
// on entry to RUnlock.
func (rw *RWMutex) RUnlock() {
	<-rw.mu
	defer func() { rw.mu <- struct{}{} }()

	if rw.state&^writeBit == 0 {
		panic("RUnlock of unlocked RWMutex")
	}

	rw.state--

	if rw.state&^writeBit == 0 && rw.state&writeBit != 0 {
		close(rw.wakeAll)
		rw.wakeAll = make(chan struct{})
	}
}

// Lock locks rw for writing.
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
func (rw *RWMutex) Lock() {
	<-rw.mu
	rw.state |= writeBit

	for rw.state&^writeBit != 0 {
		ch := rw.wakeAll
		rw.mu <- struct{}{}
		<-ch
		<-rw.mu
	}

	rw.mu <- struct{}{}
}

// Unlock unlocks rw for writing. It is a run-time error if rw is
// not locked for writing on entry to Unlock.
//
// As with Mutexes, a locked RWMutex is not associated with a particular
// goroutine. One goroutine may RLock (Lock) a RWMutex and then
// arrange for another goroutine to RUnlock (Unlock) it.
func (rw *RWMutex) Unlock() {
	<-rw.mu
	defer func() { rw.mu <- struct{}{} }()

	if rw.state&writeBit == 0 {
		panic("Unlock of unlocked RWMutex")
	}

	rw.state &^= writeBit
	close(rw.wakeAll)
	rw.wakeAll = make(chan struct{})
}
