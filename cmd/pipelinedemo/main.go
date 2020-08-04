package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/baronwithyou/pipeline/nodes"
)

func main() {
	const filename = "large.in"
	const n = 100000000

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	p := nodes.RandomSource(n)
	//p := randDemo()
	writer := bufio.NewWriter(file)
	nodes.WriterSink(writer, p)
	writer.Flush()

	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	p = nodes.ReaderSource(bufio.NewReader(file), 80)
	for v := range p {
		fmt.Println(v)
	}
}

//func randDemo() <-chan int {
//	return nodes.MergeN(
//		nodes.InMemSort(nodes.RandomSource(10)),
//		nodes.InMemSort(nodes.RandomSource(10)),
//	)
//}

//
//func arrayDemo() <-chan int {
//	return nodes.Merge(
//		nodes.InMemSort(nodes.ArraySource([]int{54, 12, 0, 123, 564, 23}...)),
//		nodes.InMemSort(nodes.ArraySource([]int{9, 1, 5, 13, 98, 6, 9, 12}...)),
//	)
//}
