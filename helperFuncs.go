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
	"github.com/k0kubun/pp"
	"github.com/mugsoft/tools"
	"github.com/mugsoft/tools/bytesize"
)

type conStruct struct {
	sync.Mutex
	con             *net.Conn
	lastSendSuccess time.Time
}

var (
	loadingConsEnd = make(chan bool)
	dataEnd        = make(chan bool)
	con            net.Conn
	err            error
	listening      bool
	s              *bufio.Scanner
	scanChan       = make(chan []byte)
	cons           = map[*net.Conn]*conStruct{}
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
				cons[cn] = &conStruct{con: cn}
				go handlePing(cn)
			}
			conBuffer = conBuffer[:0]
			bufferLock.Unlock()
			loadingConsEnd <- true
		}
	}()
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

var buffering = false

var store_lock sync.Mutex

func readData() {
	const bufsize = 1024
	var buf = make([]byte, bufsize)
	var firstTime = true
	var fullData = []byte{}
	var parsed = map[string]interface{}{}
	_, _ = fullData, parsed
	for {
		if !firstTime {
			<-loadingConsEnd
		} else {
			firstTime = false
		}
		var length int
		var remaining int
		var meta = make([]byte, 4)
		n, err := lib.GetConn().Read(meta)
		if nil != err {
			pp.Println(err)
			break
		} else if n < 4 {
			fmt.Println("Erroorrr they sent less bytes ")
			continue
		}
		// scanChan <- meta
		length = int(tools.LE2Int(meta))
		remaining = length
		for remaining > 0 {
			buffering = true
			if uint64(length) > bytesize.MB*30 {
				pp.Println(">>", remaining, len(buf))
			}
			if remaining < len(buf) {
				buf = buf[:remaining]
			} else if remaining > bufsize && len(buf) < bufsize {
				buf = make([]byte, bufsize)
			}
			n, _ = lib.GetConn().Read(buf)
			remaining -= n
			// scanChan <- buf[:n]
			fullData = append(fullData, buf[:n]...)
		}
		buffering = false
		scanChan <- fullData
		scanChan <- []byte{'\n'}
		dataEnd <- true

		// if err := json.Unmarshal(fullData, &parsed); nil != err {
		// 	pp.Println("errored json", err.Error())
		// }
		fullData = fullData[:0]
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
			if buffering {
				continue
			}
			fmt.Println("no data for 11 seconds restarting")
			os.Exit(1)
		}
		var cpy = make([]byte, len(dt))
		copy(cpy, dt)
		for _, sock := range cons {
			go func(data []byte, sock *conStruct) {
				// if sock.lastSendSuccess.Unix() != 0 && sock.lastSendSuccess.Unix() < time.Now().Unix()-int64(time.Minute) {
				// 	(*sock.con).Close()
				// 	store_lock.Lock()
				// 	delete(cons, sock.con)
				// 	store_lock.Unlock()
				// 	return
				// }
				_, err := (*sock.con).Write(data)
				if nil != err {
					store_lock.Lock()
					delete(cons, sock.con)
					store_lock.Unlock()
					(*sock.con).Close()
				}
			}(cpy, sock)
		}
	}
}
func handlePing(cn *net.Conn) {
	dataCH := make(chan bool)
	data := make([]byte, 20)
	pp.Println(len(cons))
	go func() {
		defer func() {
			store_lock.Lock()
			delete(cons, cn)
			store_lock.Unlock()
			(*cn).Close()
			pp.Println(len(cons))
		}()
		for {
			select {
			case <-dataCH:
				continue
			case <-time.After(time.Second * 15):
				return
			}
		}
		return
	}()

	for {
		n, err := (*cn).Read(data)
		if n == 0 || nil != err {
			store_lock.Lock()
			delete(cons, cn)
			store_lock.Unlock()
			(*cn).Close()
			return
		}
		dataCH <- true
	}
}
