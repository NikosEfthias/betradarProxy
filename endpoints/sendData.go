package endpoints

import (
	"net/http"
	"../lib"
	"net"
	"time"
	"../models"
	"fmt"
)

func sendToBetradar() *http.ServeMux {
	mux := http.NewServeMux()
	var con net.Conn
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		var (
			origin string = r.Form.Get("origin")
			key    string = r.Form.Get("key")
			data   string = r.Form.Get("data")
		)
		if origin == "" || key == "" || data == "" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("missing key origin or data fields"))
			return
		}
		if !models.CheckOk(origin, key) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("unauthorized"))
			return
		}
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
		fmt.Fprintln(con, data)
	})
	return mux
}
