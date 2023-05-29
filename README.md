# goseqlock
Basic implementation of a sequence lock in golang.
# How to use
```
go get github.com/MUSQQQ/goseqlock
```
```go

// func that shows the process of reading data
func ReadingData(seq *goseqlock.SeqLock, wg *sync.WaitGroup) {
	tmp := uint32(0)
	for {
		tmp = seq.RdRead()
		/*

			reading data

		*/
		if !seq.RdAgain(tmp) {
			break
		}
	}
	defer wg.Done()
}

// func that shows the process of writing data
func WritingData(seq *goseqlock.SeqLock, wg *sync.WaitGroup) {
	seq.WrLock()
	/*

		writing data

	*/
	seq.WrUnlock()
	defer wg.Done()
}
```
enjoy and criticise so I can improve it if needed :)
