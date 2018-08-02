package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"./lib"
	"github.com/mugsoft/tools"
)

var (
	loadingConsEnd = make(chan bool)
	dataEnd        = make(chan bool)
	con            net.Conn
	err            error
	listening      bool
	s              *bufio.Scanner
	scanChan       = make(chan []byte)
	data           []byte
	cons           = map[*net.Conn]*net.Conn{}
)

func startServer() {
	if listening {
		return
	}
	listening = true
	fmt.Println("listening on port", *lib.Port)
	l, err := net.Listen("tcp", ":"+*lib.Port)
	if nil != err {
		panic(err)
	}
	var bufferLock sync.Mutex
	var conBuffer []*net.Conn
	go func() {
		for {
			<-dataEnd
			bufferLock.Lock()
			for _, cn := range conBuffer {
				cons[cn] = cn
			}
			conBuffer = conBuffer[:0]
			bufferLock.Unlock()
			loadingConsEnd <- true
		}
	}()
	dataEnd <- true
	for {
		con, err := l.Accept()
		if nil != err {
			continue
		}
		bufferLock.Lock()
		conBuffer = append(conBuffer, &con)
		bufferLock.Unlock()
	}
}
func Connect() net.Conn {
	fmt.Println(*lib.Addr)
	con, err = net.Dial("tcp", *lib.Addr)
	if nil != err {
		panic(err)
	}
	return con

}
func readData() {
	for {
		<-loadingConsEnd
		var length int
		var remaining int
		var meta = make([]byte, 4)
		n, err := lib.GetConn().Read(meta)
		if nil != err {
			log.Fatalln(err)
		} else if n < 4 {
			fmt.Println("Erroorrr they sent less bytes ")
			continue
		}
		scanChan <- meta
		length = int(tools.LE2Int(meta))
		remaining = length

		if len(data) < remaining {
			data = make([]byte, remaining)
		}
	readMore:
		n, _ = lib.GetConn().Read(data[:remaining])
		remaining -= n
		scanChan <- data[:n]
		if remaining > 0 {
			goto readMore
		}
		dataEnd <- true
	}
	time.Sleep(time.Second * 11)
	log.Println("\nbetconstruct connection was interrrupted restarting")
	os.Exit(1)
}
func broadcast() {
	for {
		var dt []byte
		select {
		case dt = <-scanChan:
		case <-time.After(time.Second * 11):
			fmt.Println("no data for 11 seconds restarting")
			os.Exit(1)
		}

		for _, sock := range cons {
			_, err := (*sock).Write(dt)
			if nil != err {
				(*sock).Close()
				delete(cons, sock)
				continue
			}
		}
	}

}
