package lib

import "flag"

var (
	Addr    *string
	Port    *string
	Key     *string
	Id      *string
	Db      *string
	ApiPort *string
)

func init() {
	Addr = flag.String("h", "", "host [betradar url]")
	Port = flag.String("p", "1111", "socket port to listen")
	ApiPort = flag.String("ap", "2222", "api port to listen")
	Key = flag.String("k", "", "betradar key")
	Id = flag.String("id", "", "betradar id")
	Db = flag.String("db", "root:@tcp(127.0.0.1:3306)/test", "Database Address to use")
	if !flag.Parsed() {
		flag.Parse()
	}
}
