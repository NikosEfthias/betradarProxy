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

var con net.Conn
var err error
var listening bool
var s *bufio.Scanner

func main() {
	listening = false
	var data = make(chan string)
begin:
	con, err = net.Dial("tcp", *addr)
	if nil != err {
		panic(err)
	}
	Login(con)
	go func() {
		if listening {
			return
		}
		listening = true
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

	s = bufio.NewScanner(con)
	go func() {
		for s.Scan() {
			data <- s.Text()
		}
	}()
	for {
		var dt string
		select {
		case dt = <-data:
		case <-time.After(time.Second * 3):
			log.Println("\nprobably the connection was lost no reply for 3000 milliseconds")
			con.Close()
			time.Sleep(time.Second)
			goto begin
		}

		lock.Lock()
		for _, sock := range cons {
			_, err := fmt.Fprintln(*sock, dt)
			if nil != err {
				(*sock).Close()
				delete(cons, sock)
				continue
			}
		}
		lock.Unlock()
	}
	log.Println("\nbetradar connection was interrrupted restarting")
	con.Close()
	time.Sleep(time.Second)
	goto begin
}
