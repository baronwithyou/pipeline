package main

import (
	"fmt"
	"os"

	"github.com/baronwithyou/pipeline/nodes"
)

func main() {
	file, err := os.Open("large.out")
	if err != nil {
		panic(err)
	}

	p := nodes.ReaderSource(file, 80)

	for v := range p {
		fmt.Println(v)
	}
}
