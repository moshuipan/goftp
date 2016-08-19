package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
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
			if args := strings.Fields(fmt.Sprintf("%s", s)); args[0] == "ul" {
				if len(args) != 3 {
					continue
				}
				f, err := os.Open(args[1])
				if err != nil {
					log.Fatal(err)
					continue
				}
				conn.Write(s)
				_, err = io.Copy(conn, f)
				if err != nil {
					log.Fatal(err)
					continue
				}
				continue
			}
			_, err = conn.Write(s)
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
