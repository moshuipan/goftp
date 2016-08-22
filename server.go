package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	// "time"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	CD = "cd"
	LS = "ls"
	CP = "cp"
	UL = "ul"
)

var Root string

type Buffer []byte

func (this *Buffer) Write(w []byte) {
	for _, v := range w {
		*this = append(*this, v)
	}
}

func init() {
	var err error
	Root, err = filepath.Abs(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func main() {
	listenaddr := &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 9091,
	}
	listener, err := net.ListenTCP("tcp", listenaddr)
	if err != nil {
		fmt.Println(err)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err) // e.g., connection aborted
			continue
		}
		go handleConn(*conn) // handle one connection at a time
	}
}

func handleConn(conn net.TCPConn) {
	defer conn.Close()
	b := make([]byte, 512)
	var out Buffer
	currdir := "."
	for {
		conn.Write([]byte(currdir + "#"))
		n, err := conn.Read(b)
		if err != nil {
			fmt.Println(err)
			break
		}
		body := b[0:n]
		fmt.Printf("%s\n", body)
		s := fmt.Sprintf("%s", body)
		ss := strings.Fields(s)
		switch ss[0] {
		case LS:
			out = ls(ss, currdir)
		case CD:
			err := cd(ss, &currdir)
			if err != nil {
				out.Write([]byte(err.Error()))
			}
		case CP:
			err := cp(ss)
			if err != nil {
				out.Write([]byte(err.Error()))
			}
		case UL:
			err := upload(ss, conn)
			if err != nil {
				out.Write([]byte(err.Error()))
			}
		default:
			out.Write([]byte("unknow commond!\n"))
		}
		conn.Write(out)
		out = nil
	}
}
func upload(args []string, conn net.TCPConn) error {
	//ul dst src
	if len(args) != 3 {
		return errors.New("ul dst src\n")
	}
	_, filename := filepath.Split(args[2])
	name := filepath.Join(args[1], filename)
	f, err := os.Create(name)
	if err != nil {
		return errors.New(err.Error() + "\n")
	}
	defer f.Close()
	// conn.SetReadBuffer(1024 * 10000)
	//end:0xda
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		// _, err = io.Copy(f, conn)
		if err != nil {
			return errors.New(err.Error() + "\n")
		}
		b := buf[0:n]
		if len(b) == 1 {
			fmt.Println("upload end!")
			break
		} else {
			if len(b) == 2 && b[1] == 0xda {
				b = b[0:0]
			}
			_, err = f.Write(b)
			if err != nil {
				return errors.New(err.Error() + "\n")
			}
		}
	}
	return nil
}
func cp(args []string) error {
	//cp dstdir+dstfilename src
	if len(args) != 3 {
		return errors.New("cp dstdir+dstfilename src\n")
	}
	src, err := os.Open(args[2])
	if err != nil {
		return errors.New(err.Error() + "\n")
	}
	dst, err := os.Create(args[1])
	if err != nil {
		return errors.New(err.Error() + "\n")
	}
	_, err = io.Copy(dst, src)
	if err != nil {
		return errors.New(err.Error() + "\n")
	}
	return nil
}
func cd(args []string, currdir *string) error {
	//cd ..判断cd后的目录权限
	path := filepath.Join(*currdir, args[1])
	p, err := filepath.Abs(path)
	if err != nil {
		return errors.New(err.Error() + "\n")
	}
	if strings.Contains(p, Root) {
		*currdir = path
		return nil
	}
	return errors.New("路径权限不够!\n")
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
				out.Write([]byte("read dir error!\n"))
				return
			}
		} else if args[1] == "-l" {
			f, err = ioutil.ReadDir(currdir)
			if err != nil {
				out.Write([]byte("read dir error!\n"))
				return
			}
		} else {
			f, err = ioutil.ReadDir(args[1])
			if err != nil {
				out.Write([]byte("read dir error!\n"))
				return
			}
		}
	} else {
		f, err = ioutil.ReadDir(currdir)
		if err != nil {
			out.Write([]byte("read dir error!\n"))
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
		out.Write([]byte("\n"))
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
