package async

import (
	"sync"
	"sync/atomic"
)

// Promise represents a value that will be resolved in the future
type Promise[T any] struct {
	value  T
	err    error
	done   chan struct{}
	once   sync.Once
}

// NewPromise creates a new promise
func NewPromise[T any]() *Promise[T] {
	return &Promise[T]{
		done: make(chan struct{}),
	}
}

// Resolve resolves the promise with a value
func (p *Promise[T]) Resolve(value T) {
	p.once.Do(func() {
		p.value = value
		close(p.done)
	})
}

// Reject rejects the promise with an error
func (p *Promise[T]) Reject(err error) {
	p.once.Do(func() {
		p.err = err
		close(p.done)
	})
}

// Await waits for the promise to be resolved and returns the value
func (p *Promise[T]) Await() (T, error) {
	<-p.done
	return p.value, p.err
}

// Done returns a channel that is closed when the promise is resolved
func (p *Promise[T]) Done() <-chan struct{} {
	return p.done
}

// Future represents a read-only promise
type Future[T any] interface {
	Await() (T, error)
	Done() <-chan struct{}
}

// Async runs a function asynchronously and returns a promise
func Async[T any](fn func() (T, error)) *Promise[T] {
	p := NewPromise[T]()
	go func() {
		value, err := fn()
		if err != nil {
			p.Reject(err)
		} else {
			p.Resolve(value)
		}
	}()
	return p
}

// All waits for all promises to be resolved
func All[T any](promises ...*Promise[T]) *Promise[[]T] {
	p := NewPromise[[]T]()
	go func() {
		values := make([]T, len(promises))
		for i, promise := range promises {
			value, err := promise.Await()
			if err != nil {
				p.Reject(err)
				return
			}
			values[i] = value
		}
		p.Resolve(values)
	}()
	return p
}

// Race returns the first resolved promise
func Race[T any](promises ...*Promise[T]) *Promise[T] {
	p := NewPromise[T]()
	go func() {
		for _, promise := range promises {
			go func(pr *Promise[T]) {
				value, err := pr.Await()
				if err != nil {
					p.Reject(err)
				} else {
					p.Resolve(value)
				}
			}(promise)
		}
	}()
	return p
}

// Mutex is a sync.Mutex wrapper with TryLock
type Mutex struct {
	mu sync.Mutex
}

// Lock locks the mutex
func (m *Mutex) Lock() {
	m.mu.Lock()
}

// Unlock unlocks the mutex
func (m *Mutex) Unlock() {
	m.mu.Unlock()
}

// TryLock tries to lock the mutex without blocking
func (m *Mutex) TryLock() bool {
	return m.mu.TryLock()
}

// RWMutex is a sync.RWMutex wrapper with TryLock
type RWMutex struct {
	mu sync.RWMutex
}

// Lock locks for writing
func (m *RWMutex) Lock() {
	m.mu.Lock()
}

// Unlock unlocks for writing
func (m *RWMutex) Unlock() {
	m.mu.Unlock()
}

// RLock locks for reading
func (m *RWMutex) RLock() {
	m.mu.RLock()
}

// RUnlock unlocks for reading
func (m *RWMutex) RUnlock() {
	m.mu.RUnlock()
}

// TryLock tries to lock for writing without blocking
func (m *RWMutex) TryLock() bool {
	return m.mu.TryLock()
}

// TryRLock tries to lock for reading without blocking
func (m *RWMutex) TryRLock() bool {
	return m.mu.TryRLock()
}

// AtomicBool is an atomic boolean
type AtomicBool struct {
	value int32
}

// NewAtomicBool creates a new atomic bool
func NewAtomicBool(value bool) *AtomicBool {
	var v int32
	if value {
		v = 1
	}
	return &AtomicBool{value: v}
}

// Get returns the current value
func (b *AtomicBool) Get() bool {
	return atomic.LoadInt32(&b.value) == 1
}

// Set sets the value
func (b *AtomicBool) Set(value bool) {
	var v int32
	if value {
		v = 1
	}
	atomic.StoreInt32(&b.value, v)
}

// CompareAndSwap compares and swaps
func (b *AtomicBool) CompareAndSwap(old, new bool) bool {
	var oldV, newV int32
	if old {
		oldV = 1
	}
	if new {
		newV = 1
	}
	return atomic.CompareAndSwapInt32(&b.value, oldV, newV)
}

// AtomicInt is an atomic int
type AtomicInt struct {
	value int64
}

// NewAtomicInt creates a new atomic int
func NewAtomicInt(value int) *AtomicInt {
	return &AtomicInt{value: int64(value)}
}

// Get returns the current value
func (i *AtomicInt) Get() int {
	return int(atomic.LoadInt64(&i.value))
}

// Set sets the value
func (i *AtomicInt) Set(value int) {
	atomic.StoreInt64(&i.value, int64(value))
}

// Add adds a value and returns the new value
func (i *AtomicInt) Add(delta int) int {
	return int(atomic.AddInt64(&i.value, int64(delta)))
}

// Increment increments by 1
func (i *AtomicInt) Increment() int {
	return i.Add(1)
}

// Decrement decrements by 1
func (i *AtomicInt) Decrement() int {
	return i.Add(-1)
}

// CompareAndSwap compares and swaps
func (i *AtomicInt) CompareAndSwap(old, new int) bool {
	return atomic.CompareAndSwapInt64(&i.value, int64(old), int64(new))
}

// Once is a sync.Once wrapper
type Once struct {
	once sync.Once
	done int32
}

// Do executes the function once
func (o *Once) Do(fn func()) {
	o.once.Do(fn)
}

// Done returns true if the function has been executed
func (o *Once) Done() bool {
	return atomic.LoadInt32(&o.done) == 1
}

// WaitGroup is a sync.WaitGroup wrapper
type WaitGroup struct {
	wg sync.WaitGroup
}

// Add adds delta to the counter
func (w *WaitGroup) Add(delta int) {
	w.wg.Add(delta)
}

// Done decrements the counter
func (w *WaitGroup) Done() {
	w.wg.Done()
}

// Wait waits for the counter to reach zero
func (w *WaitGroup) Wait() {
	w.wg.Wait()
}

// Semaphore is a counting semaphore
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore creates a new semaphore with the given capacity
func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, capacity),
	}
}

// Acquire acquires a permit
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// TryAcquire tries to acquire a permit without blocking
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release releases a permit
func (s *Semaphore) Release() {
	<-s.ch
}

// Available returns the number of available permits
func (s *Semaphore) Available() int {
	return cap(s.ch) - len(s.ch)
}
