package main

import (
	"log"

	"./endpoints"
	"./lib"
)

func init() {
	log.SetFlags(0)
	// go func() {
	// 	for {
	// 		time.Sleep(time.Second)
	// 		fmt.Printf("\r\x1B[32mConnected Users (%d)\x1B[0m", len(cons))
	// 	}
	// }()
}

func main() {
	listening = false
	go endpoints.StartListening()
	lib.SetConn(Connect())
	Login(lib.GetConn())
	go startServer()
	go readData()
	broadcast()
}
