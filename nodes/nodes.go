package nodes

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"
)

var startTime time.Time

func Init() {
	startTime = time.Now()
}

// ArraySource source formatting by array
func ArraySource(a ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()

	return out
}

// InMemSort sort in memory
func InMemSort(in <-chan int) <-chan int {
	out := make(chan int, 1024)

	go func() {
		// Read into memory
		var a []int
		for v := range in {
			a = append(a, v)
		}

		fmt.Println("Read done", time.Since(startTime).Seconds())

		// Sort
		sort.Ints(a)

		fmt.Println("Sort done", time.Since(startTime).Seconds())

		// Output
		for _, v := range a {
			out <- v
		}
		close(out)
	}()

	return out
}

// Merge ...
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int, 1024)

	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2

		for ok1 || ok2 {
			// Error Occur without ok1 condition
			if !ok2 || (ok1 && v1 < v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
		fmt.Println("Merge done", time.Since(startTime).Seconds())
	}()

	return out
}

// ReaderSource source formatting by reader
func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int, 1024)

	go func() {
		buffer := make([]byte, 8)
		var byteRead int
		for {
			n, err := reader.Read(buffer)
			byteRead += n
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil || (chunkSize != -1 && byteRead >= chunkSize) {
				break
			}
		}
		close(out)
	}()

	return out
}

// RandomSource ...
func RandomSource(count int) <-chan int {
	out := make(chan int)
	rand.Seed(time.Now().UnixNano())

	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()

	return out
}

// WriterSink ...
func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		_, _ = writer.Write(buffer)
	}
}

// MergeN ...
func MergeN(in ...<-chan int) <-chan int {
	if len(in) == 1 {
		return in[0]
	}

	m := len(in) / 2
	return Merge(
		MergeN(in[:m]...),
		MergeN(in[m:]...),
	)
}
