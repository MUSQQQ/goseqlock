package goseqlock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	MockSeqLock := New()
	MockSeqLock.WrLock()
	MockSeqLock.WrUnlock()

	require.Equal(t, uint64(2), MockSeqLock.Counter)
}

func TestReset(t *testing.T) {
	MockSeqLock := New()
	MockSeqLock.Counter = 30

	require.Equal(t, uint64(30), MockSeqLock.Counter)

	MockSeqLock.ResetCounter()

	require.Equal(t, uint64(0), MockSeqLock.Counter)
}

func TestRead(t *testing.T) {
	MockSeqLock := New()
	MockSeqLock.Counter = 30

	require.Equal(t, uint64(30), MockSeqLock.RdRead())
}

func TestAgain(t *testing.T) {
	MockSeqLock := New()
	MockSeqLock.Counter = 1

	require.Equal(t, true, MockSeqLock.RdAgain(1))

	MockSeqLock.Counter = 2

	require.Equal(t, false, MockSeqLock.RdAgain(2))
}

func TestReadWriteIntegration(t *testing.T) {
	MockSeqLock := New()
	MockSeqLock.WrLock()

	require.Equal(t, uint64(1), MockSeqLock.RdRead())
	require.Equal(t, true, MockSeqLock.RdAgain(1))

	MockSeqLock.WrUnlock()

	require.Equal(t, uint64(2), MockSeqLock.RdRead())
	require.Equal(t, false, MockSeqLock.RdAgain(2))
}

func TestTimeBlocking(t *testing.T) {
	MockSeqLock := New()
	go MockSeqLock.TimeBlock(1500)
	time.Sleep(500 * time.Millisecond)

	require.Equal(t, uint64(1), MockSeqLock.RdRead())

	time.Sleep(500 * time.Millisecond)

	require.Equal(t, uint64(1), MockSeqLock.RdRead())

	time.Sleep(1000 * time.Millisecond)

	require.Equal(t, uint64(2), MockSeqLock.RdRead())
}
