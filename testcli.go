package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9091")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	exit := make(chan bool)
	go func() {
		for {
			scan := bufio.NewReader(os.Stdin)
			s, ok, err := scan.ReadLine()
			if err != nil {
				fmt.Println(err, ok)
			}
			_, err = conn.Write([]byte(s))
			if err != nil {
				fmt.Println(err)
			}
		}
		exit <- true
	}()
	mustCopy(os.Stdout, conn)
	<-exit
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
