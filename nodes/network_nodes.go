package nodes

import (
	"bufio"
	"net"
)

func NetworkSink(addr string, in <-chan int) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	go func() {
		defer listener.Close()

		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		writer := bufio.NewWriter(conn)
		defer writer.Flush()
		WriterSink(writer, in)
	}()
}

func NetworkSource(addr string) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		r := ReaderSource(bufio.NewReader(conn), -1)
		for v := range r {
			out <- v
		}
	}()

	return out
}
