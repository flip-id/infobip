package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func Connect(dbUser string, dbPass string, dbName string, dbHost string) (*gorm.DB, error) {
	var connString = fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)

	var db, err = gorm.Open("mysql", connString)

	// Can't connect to db
	if err != nil {
		log.Fatalln(err)
	}

	return db, err
}
