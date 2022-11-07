package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Init() {
	fmt.Printf("Initializing connection to %v", os.Getenv("DATABASE")+"\n")
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(mariadb:3306)/"+os.Getenv("DATABASE"))
	if err != nil {
		fmt.Println(err.Error())
	}

	if err = db.Ping(); err != nil {
		fmt.Println(err.Error())
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func GetConnection() *sql.DB {
	return db
}
