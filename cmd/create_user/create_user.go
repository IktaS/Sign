package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/IktaS/sign/service"
	"github.com/IktaS/sign/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func CreateUser(username string, password string, fullname string) (int, error) {
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		return -1, err
	}

	defer db.Close()

	signRepo, err := sqlite.NewSQLiteDB(db)
	if err != nil {
		return -1, err
	}

	signService := service.NewSignService(signRepo, os.Getenv("VERIFY_PATH"))

	return signService.CreateUser(context.Background(), username, password, fullname)
}

func main() {
	username := flag.String("username", "", "username of user")
	password := flag.String("password", "", "password of user")
	fullname := flag.String("fullname", "", "full name of user")

	flag.Parse()

	if username == nil || len(*username) <= 0 {
		panic("username cannot be empty")
	}

	if password == nil || len(*password) <= 0 {
		panic("password cannot be empty")
	}

	if fullname == nil || len(*fullname) <= 0 {
		panic("fullname cannot be empty")
	}

	id, err := CreateUser(*username, *password, *fullname)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created user with \nID: %d\nusername: %s\nfullname: %s\n", id, *username, *fullname)
}
