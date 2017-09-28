package lib

import _ "github.com/jinzhu/gorm/dialects/mysql"
import (
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

func GetDb() *gorm.DB {
	if nil != db {
		return db
	}
	db, err = gorm.Open("mysql", *Db)
	if nil != err {
		panic(err)
	}
	err = db.DB().Ping()
	if nil != err {
		panic(err)
	}
	return db
}
