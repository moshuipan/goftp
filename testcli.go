package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var clock chan bool

type Buffer []byte

func (this *Buffer) Write(w []byte) {
	for _, v := range w {
		*this = append(*this, v)
	}
}
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
	clock = make(chan bool)
	// clock <- true
	go func() {
		for {
			scan := bufio.NewReader(os.Stdin)
			s, ok, err := scan.ReadLine()
			if err != nil {
				fmt.Println(err, ok)
			}
			if args := strings.Fields(fmt.Sprintf("%s", s)); len(args) <= 0 {
				continue
			}
			if args := strings.Fields(fmt.Sprintf("%s", s)); args[0] == "ul" {
				if len(args) != 3 {
					fmt.Println("ul dst src")
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
							break
						}
					}
					start += int64(n)
					b := buf[0:n]
					if len(b) == 1 {
						b = append(b, 0xda)
					}
					if len(b) == 0 {
						fmt.Println("read all file!")
						break
					}
					_, err = conn.Write(b)
					if err != nil {
						fmt.Println("send file error!", err)
						break
					}
				}
				_, err = conn.Write([]byte{0xda})
				if err != nil {
					fmt.Println("send file error!")
					break
				}
				f.Close()
				clock <- true
				// conn.CloseWrite()
				continue
			}
			if args := strings.Fields(fmt.Sprintf("%s", s)); args[0] == "dl" {
				if len(args) != 3 {
					fmt.Println("dl dst src")
					continue
				}
				_, filename := filepath.Split(args[2])
				name := filepath.Join(args[1], filename)
				f, err := os.Create(name)
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, err = conn.Write(s)
				if err != nil {
					fmt.Println(err)
					f.Close()
					continue
				}
				buf := make([]byte, 1024)
				for {
					n, err := conn.Read(buf)
					// _, err = io.Copy(f, conn)
					if err != nil {
						fmt.Println(err)
						break
					}
					b := buf[0:n]
					if len(b) == 1 && b[0] == 0xda {
						fmt.Println("download end!")
						break
					} else {
						if len(b) == 2 && b[1] == 0xda {
							b = b[0:0]
						}
						_, err = f.Write(b)
						if err != nil {
							fmt.Println(err)
							break
						}
					}
				}
				f.Close()
				clock <- true
				continue
			}
			_, err = conn.Write(s)
			if err != nil {
				fmt.Println(err)
			}
			clock <- true
		}
		exit <- true
	}()
	mustCopy(os.Stdout, conn)
	<-exit
}

func mustCopy(dst io.Writer, src net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		b := buf[0:n]
		os.Stdout.Write(b)
		<-clock
	}
	/*if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}*/
}
