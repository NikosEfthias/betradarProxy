package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"./lib"
	"github.com/mugsoft/tools"
)

type Command struct {
	Command string                   `json:"Command"`
	Params  []map[string]interface{} `json:"Params,omitempty"`
	Objects []map[string]interface{} `json:"Objects,omitempty"`
}

func (l *Command) Send(sock net.Conn) error {
	d, err := json.Marshal(l)
	if nil != err {
		return err
	}
	fmt.Println(string(d))
	ln := tools.Int2LE(uint(len(d)))
	_, err = sock.Write(ln[:])
	if nil != err {
		return err
	}
	_, err = sock.Write(d)
	return err
}

func LoginWithValues(uname string, pass string) *Command {
	return &Command{
		Command: "Login",
		Params:  []map[string]interface{}{{"UserName": uname, "Password": pass}},
	}
}
func Login(sock net.Conn) {
	LoginWithValues(*lib.Key, *lib.Pass).Send(sock)
	var cmds = []*Command{
		{
			Command: "GetSports",
		},
		{
			Command: "GetRegions",
		},
		{
			Command: "GetCompetitions",
		},
		{
			Command: "GetMarketTypes",
			Params:  []map[string]interface{}{{}},
		},
	}
	for _, cmd := range cmds {
		// time.Sleep(time.Second)
		cmd.Send(sock)
	}
	go func() {
		for {
			(&Command{Command: "HeartBeat"}).Send(lib.GetConn())
			time.Sleep(time.Second * 3)
		}
	}()
}
