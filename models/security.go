package models

import "github.com/jinzhu/gorm"
import (
	"../lib"
)

var Model *gorm.DB

type Origin struct {
	ID       string `gorm:"column:id;primary_key"json:"-"`
	Password string `gorm:"default:null"json:"-"`
	Admin    bool `gorm:"column:admin"json:"-"`
}

func (a *Origin) TableName() string {
	return "Origins"
}

func init() {
	Model = lib.GetDb().Model(&Origin{})
	Model.AutoMigrate(&Origin{})
}
func CheckOk(org string, key string) bool {
	_ = key //not implemented yet
	out := new(Origin)
	Model.Where(&Origin{ID: org, Password: key, Admin: true}).First(out)
	if out.ID != org {
		return false
	}
	return true
}
