package main

import (
	"./lib"
	"net"
	"bufio"
	"sync"
	"time"
	"fmt"
	"log"
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
}

var con net.Conn
var err error
var listening bool
var s *bufio.Scanner

func main() {
	listening = false
	var data = make(chan string)
begin:
	con, err = net.Dial("tcp", *lib.Addr)
	if nil != err {
		panic(err)
	}
	lib.SetConn(con)
	Login(lib.GetConn())
	go func() {
		if listening {
			return
		}
		listening = true
		fmt.Println("listening on port", *lib.Port)
		l, err := net.Listen("tcp", ":" + *lib.Port)
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

	s = bufio.NewScanner(lib.GetConn())
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
			lib.GetConn().Close()
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
	lib.GetConn().Close()
	time.Sleep(time.Second)
	goto begin
}
