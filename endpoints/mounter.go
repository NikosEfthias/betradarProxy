package endpoints

import "net/http"
import (
	"../lib"
	"fmt"
)

func StartListening() {
	var mainMux = http.NewServeMux()
	mainMux.Handle("/send/", http.StripPrefix("/send", sendToBetradar()))
	fmt.Println("http server listening on port:", *lib.ApiPort)
	panic(http.ListenAndServe(":" + *lib.ApiPort, mainMux))
}
