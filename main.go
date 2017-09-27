package main

import (
	"flag"
	"net"
	"bufio"
	"sync"
	"time"
	"fmt"
	"log"
)

var (
	addr *string
	port *string
	key  *string
	id   *string
)
var cons = map[*net.Conn]*net.Conn{}
var lock sync.Mutex

func init() {
	log.SetFlags(0)
	go func() {
		for {
			time.Sleep(time.Second / 2)
			fmt.Printf("\r\x1B[32mConnected Users (%d)\x1B[0m", len(cons))
		}
	}()
	addr = flag.String("h", "", "host [betradar url]")
	port = flag.String("p", "1111", "port to listen")
	key = flag.String("k", "", "betradar key")
	id = flag.String("id", "", "betradar id")
	flag.Parse()
}
func main() {
	con, err := net.Dial("tcp", *addr)
	if nil != err {
		panic(err)
	}
	Login(con)
	go func() {
		fmt.Println("listening on port", *port)
		l, err := net.Listen("tcp", ":" + *port)
		if nil != err {
			panic(err)
		}
		for {
			con, err := l.Accept()
			if nil != err {
				continue
			}
			lock.Lock()
			cons[&con] = &con
			lock.Unlock()
		}
	}()
	s := bufio.NewScanner(con)
	for s.Scan() {
		lock.Lock()
		for _, sock := range cons {
			_, err := fmt.Fprintln(*sock, s.Text())
			if nil != err {
				(*sock).Close()
				delete(cons, sock)
				continue
			}
		}
		lock.Unlock()
	}
	log.Fatalln("\nbetradar connection was interrrupted")
}
