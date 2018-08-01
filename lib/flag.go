package lib

import "flag"

var (
	Addr    *string
	Port    *string
	Key     *string
	Pass    *string
	ApiPort *string
)

func init() {
	Addr = flag.String("h", "odds-stream-test.betconstruct.com:8077", "host [betconstruct url]")
	Port = flag.String("p", "1111", "socket port to listen")
	ApiPort = flag.String("ap", "2222", "api port to listen")
	Key = flag.String("k", "", "login name")
	Pass = flag.String("pass", "", "betconstruct Password")
	if !flag.Parsed() {
		flag.Parse()
	}
}
