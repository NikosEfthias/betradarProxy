package endpoints

import (
	"net"
	"net/http"
	"time"

	"../lib"
	"github.com/k0kubun/pp"
	"github.com/mugsoft/tools"
)

func sendToBetradar() *http.ServeMux {
	mux := http.NewServeMux()
	var con net.Conn
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		var (
			data string = r.Form.Get("data")
		)
		if data == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("empty data"))
			return
		}
	checkAgain:
		con = lib.GetConn()
		if nil == con {
			time.Sleep(time.Millisecond * 100)
			goto checkAgain
		}
		length := tools.Int2LE(uint(len([]byte(data))))
		pp.Println(length, data)
		con.Write(length[:])
		con.Write([]byte(data))
	})
	return mux
}
