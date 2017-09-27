package main

import "fmt"
import (
	"net"
	"encoding/xml"
)

type Ask struct {
	XMLName     xml.Name `xml:"BookmakerStatus"`
	Bookmakerid string   `xml:"bookmakerid,attr"`
	Type        string   `xml:"type,attr"`
	Timestamp   int64    `xml:"timestamp,attr"`
	Key         string   `xml:"key,attr"`
}

func (ask *Ask) Send(sock net.Conn) error {
	d, err := xml.Marshal(ask)
	if nil != err {
		return err
	}
	d = append(d, '\n')
	fmt.Print(string(d))
	_, err = sock.Write(d)
	return err
}

func AskWithValues(id string, tp string, ts int64, key string) *Ask {
	return &Ask{
		Bookmakerid: id,
		Type:        tp,
		Timestamp:   ts,
		Key:         key,
	}
}
func Login(sock net.Conn) {
	AskWithValues(*id, "login", 0, *key).Send(sock)
	sock.Write([]byte("<StartAuto/>\n"))
}
