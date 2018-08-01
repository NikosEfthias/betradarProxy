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
	"github.com/mugsoft/tools"
)

var cons = map[*net.Conn]*net.Conn{}
var lock sync.Mutex

func init() {
	log.SetFlags(0)
	go func() {
		for {
			time.Sleep(time.Second)
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
	fmt.Println(*lib.Addr)
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
	var scanChan = make(chan []byte)
	go func() {
		totalN := 0
		for {
			var totalData []byte
			var meta = make([]byte, 4)
			n, err := lib.GetConn().Read(meta)
			var length = int(tools.LE2Int(meta))
			if nil != err {
				log.Fatalln(err)
			} else if n < 4 {
				fmt.Println("Erroorrr they sent less bytes ")
				continue
			}

			var data = make([]byte, length)
		readMore:
			n, _ = lib.GetConn().Read(data)
			totalN += n
			totalData = append(totalData, data[:n]...)
			if totalN < length {
				goto readMore
			}
			scanChan <- totalData[:totalN]
			totalN = 0
		}
		time.Sleep(time.Second * 15)
		log.Println("\nbetconstruct connection was interrrupted restarting")
		os.Exit(1)
	}()
	for {
		var dt []byte
		select {
		case dt = <-scanChan:
		case <-time.After(time.Minute):
			fmt.Println("no data for a minute restarting")
			os.Exit(1)
		}
		lock.Lock()
		for _, sock := range cons {
			_, err := fmt.Fprintln(*sock, tools.Int2LE(uint(len(dt))))
			_, err = fmt.Fprintln(*sock, dt)
			if nil != err {
				(*sock).Close()
				delete(cons, sock)
				continue
			}
		}
		lock.Unlock()
	}
}
