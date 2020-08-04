package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/baronwithyou/pipeline/nodes"
)

func main() {
	inFilename, outFilename := "small.in", "small.out"

	p, err := createNetworkPipeline(inFilename, 512, 4)
	if err != nil {
		panic(err)
	}

	//for v := range p {
	//	fmt.Printf("%d\t", v)
	//}
	//
	//return

	if err := writeToFile(outFilename, p); err != nil {
		panic(err)
	}

	if err := printFile(outFilename); err != nil {
		panic(err)
	}
}

// goal: fetch numbers from file and sort them respectively
func createPipeline(filename string, fileSize, chunkCount int) (<-chan int, error) {
	chunkSize := fileSize / chunkCount
	var sortResults []<-chan int
	nodes.Init()

	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		_, _ = file.Seek(int64(chunkSize*i), 0)

		source := nodes.ReaderSource(bufio.NewReader(file), chunkSize)

		sortResults = append(sortResults, nodes.InMemSort(source))
	}

	return nodes.MergeN(sortResults...), nil
}

func createNetworkPipeline(filename string, fileSize, chunkCount int) (<-chan int, error) {
	chunkSize := fileSize / chunkCount
	var sortResults []<-chan int
	nodes.Init()

	var addrs []string
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		_, _ = file.Seek(int64(chunkSize*i), 0)

		source := nodes.ReaderSource(bufio.NewReader(file), chunkSize)

		addr := ":" + strconv.Itoa(7000+i)
		// 将内部排序好的数据通过network传输到
		nodes.NetworkSink(addr, nodes.InMemSort(source))
		addrs = append(addrs, addr)
	}

	// 开启客户端获取数据
	for _, addr := range addrs {
		sortResults = append(sortResults, nodes.NetworkSource(addr))
	}

	return nodes.MergeN(sortResults...), nil
}

func writeToFile(filename string, p <-chan int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	nodes.WriterSink(writer, p)
	defer writer.Flush()

	return nil
}

func printFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	p := nodes.ReaderSource(bufio.NewReader(file), 80)

	for v := range p {
		fmt.Println(v)
	}
	return nil
}
