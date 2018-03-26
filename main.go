package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"./endpoints"
	"./lib"
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
	go endpoints.StartListening()
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
		l, err := net.Listen("tcp", ":"+*lib.Port)
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
	var scanChan = make(chan string)
	go func() {
		for s.Scan() {
			scanChan <- s.Text()
		}
		time.Sleep(time.Second * 5)
		log.Println("\nbetradar connection was interrrupted restarting")
	}()
	for {
		var dt string
		select {
		case dt = <-scanChan:
		case <-time.After(time.Minute):
			fmt.Println("no data for a minute restarting")
			os.Exit(0)
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
}
