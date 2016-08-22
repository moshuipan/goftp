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
	rmoteaddr := &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 9091,
	}
	conn, err := net.DialTCP("tcp", nil, rmoteaddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	exit := make(chan bool)
	go func() {
	again:
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
				f, err := os.Open(args[2])
				if err != nil {
					fmt.Println(err)
					continue
				}
				conn.Write(s)
				buf := make([]byte, 1024)
				var start int64 = 0
				for {
					n, err := f.ReadAt(buf, start)
					if err != nil {
						if err != io.EOF {
							fmt.Println("read file error!", err)
							f.Close()
							continue again
						}
					}
					start += int64(n)
					b := buf[0:n]
					if len(b) == 1 {
						b = append(b, 0xda)
					}
					if len(b) == 0 {
						fmt.Println("read all file!")
						f.Close()
						break
					}
					_, err = conn.Write(b)
					if err != nil {
						fmt.Println("send file error!", err)
						f.Close()
						continue again
					}
				}
				_, err = conn.Write([]byte{0xda})
				if err != nil {
					fmt.Println("send file error!")
					f.Close()
					continue again
				}
				f.Close()
				// conn.CloseWrite()
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
