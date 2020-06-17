package db

import (
	"database/sql"
	"fmt"
	"os"

	// package used to read the .env file
	"github.com/golang/glog"
)

func CreatePGClient() *sql.DB {

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		glog.Error("create pg connection err:", err)
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		glog.Error("ping postgres err:", err)
		panic(err)
	}

	fmt.Println("Successfully connected postgres!")

	return db
}
