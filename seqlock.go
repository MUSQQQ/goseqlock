package goseqlock

import (
	"sync"
	"sync/atomic"
	"time"
)

// SeqLock
// Counter increases by 1 everytime the lock gets taken or released by a writer
// sequence number is the other name for it
// inherited mutex is used only by writers
type SeqLock struct {
	Counter uint64
	sync.Mutex
}

func New() *SeqLock {
	seqlock := SeqLock{Counter: 0}
	return &seqlock
}

// RdRead reads counter using atomic operations
// used before accessing data protected by the lock
func (seq *SeqLock) RdRead() uint64 {
	return atomic.LoadUint64(&seq.Counter)
}

// RdAgain checks if data is not being modified by writer
// or if it has not been modified since rdRead func
func (seq *SeqLock) RdAgain(val uint64) bool {
	return (atomic.LoadUint64(&seq.Counter)&1) != 0 || val != seq.Counter
}

// ResetCounter resets counter to zero
func (seq *SeqLock) ResetCounter() {
	seq.Lock()
	atomic.SwapUint64(&seq.Counter, 0)
	seq.Unlock()
}

// WrLock
// counter becomes odd when writer
// starts modifying data
func (seq *SeqLock) WrLock() {
	seq.Lock()
	atomic.AddUint64(&seq.Counter, 1)
}

// WrUnlock
// counter becomes even when writer
// starts modifying data
func (seq *SeqLock) WrUnlock() {
	atomic.AddUint64(&seq.Counter, 1)
	seq.Unlock()
}

// TimeBlock locks the data for both writers and readers
// for the given amount of milliseconds
// beware that it locks once all
// writers before it are finished
func (seq *SeqLock) TimeBlock(ms int64) {
	seq.Lock()
	atomic.AddUint64(&seq.Counter, 1)
	time.Sleep(time.Duration(ms) * time.Millisecond)
	atomic.AddUint64(&seq.Counter, 1)
	seq.Unlock()
}

// TimeWriteBenchmark checks for how much time the lock
// is blocked until new writer gets
// access to it (in ms)
// better to use it multiple times
// and calculate the average
func (seq *SeqLock) TimeWriteBenchmark() uint64 {
	before := time.Now().Unix()
	seq.Lock()
	seq.Unlock()
	after := time.Now().Unix()

	return uint64(after-before) / 1000
}

// TimeReadBenchmark checks for how much time the reader
// has to keep repeating the
// reading process (in ms)
// better to use it multiple times
// and average results
func (seq *SeqLock) TimeReadBenchmark() uint64 {
	before := time.Now().Unix()
	for seq.RdAgain(seq.RdRead()) {
	}
	after := time.Now().Unix()

	return uint64(after-before) / 1000
}
