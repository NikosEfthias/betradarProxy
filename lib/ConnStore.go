package lib

import (
	"net"
	"sync"
)

var con net.Conn
var l sync.Mutex

func GetConn() net.Conn {
	l.Lock()
	defer l.Unlock()
	return con
}
func SetConn(conn net.Conn) {
	l.Lock()
	defer l.Unlock()
	con = conn
}
