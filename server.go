package main

import (
	"fmt"
	"io/ioutil"
	"net"
	// "time"
	"os"
	"path/filepath"
	"strings"
)

const (
	CD = "cd"
	LS = "ls"
)

type Buffer []byte

func (this *Buffer) Write(w []byte) {
	for _, v := range w {
		*this = append(*this, v)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":9091")
	if err != nil {
		fmt.Println(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle one connection at a time
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	b := make([]byte, 512)
	var out Buffer
	currdir := filepath.Join(".")
	for {
		conn.Write([]byte(currdir + "#"))
		n, err := conn.Read(b)
		if err != nil {
			fmt.Println(err)
			break
		}
		body := b[0:n]
		fmt.Println(body)
		s := fmt.Sprintf("%s", body)
		ss := strings.Split(s, " ")
		switch ss[0] {
		case LS:
			out = ls(ss, currdir)
		default:
			out.Write([]byte("unknow commond!\n"))
		}
		conn.Write(out)
		out = nil
	}
}

func ls(args []string, currdir string) (out Buffer) {
	//three args
	//ls [-l]  [dir]
	var f []os.FileInfo
	var err error
	if len(args) > 1 {
		if len(args) == 3 {
			f, err = ioutil.ReadDir(args[2])
			if err != nil {
				out.Write([]byte("read dir error!"))
				return
			}
		} else if args[1] == "-l" {
			f, err = ioutil.ReadDir(currdir)
			if err != nil {
				out.Write([]byte("read dir error!"))
				return
			}
		} else {
			out.Write([]byte("unknow args!"))
			return
		}
	} else {
		f, err = ioutil.ReadDir(currdir)
		if err != nil {
			out.Write([]byte("read dir error!"))
			return
		}
	}
	if len(args) >= 2 && args[1] == "-l" {
		for _, v := range f {
			out.Write([]byte(fmt.Sprint(v.Mode()) + "\t" + fmt.Sprint(v.Size()) + "\t" + v.Name() + "\n"))
		}
	} else {
		for _, v := range f {
			// out.Write([]byte(fmt.Sprint(v.Mode()) + "\t" + fmt.Sprint(v.Size()) + "\t" + v.Name() + "\n"))
			out.Write([]byte(v.Name() + "\t"))
		}
	}
	return
}

/*func handleConn(c net.Conn) {
	defer c.Close()
	for {
		b := make([]byte, 512)
		n, err := c.Read(b)
		if err != nil {
			fmt.Println("read:", err)
			c.Close()
			break
		}
		b = b[0:n]
		s := fmt.Sprintf("%s", b)
		ss := strings.Split(s, " ")
		fmt.Println(ss)
		if ss[0] != "ls" {
			fmt.Printf("%x\n%x", ss[0], "ls")
		}
		cmd := exec.Command(ss[0])
		cmd.Path = "/bin/sh/"
		for i := 1; i < len(ss); i++ {
			cmd.Args = append(cmd.Args, ss[i])
		}
		out, err := cmd.Output()
		if err != nil {
			fmt.Println("out:", err)
		}
		c.Write(out)
	}
}*/
