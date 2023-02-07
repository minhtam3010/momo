package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "quynhnhu2010"
	dbName   = "momo"
)

type ConnectDB struct {
	tx *sql.Tx
}

// Connect to mySql
func Connect() *sql.DB {
	// Open a database connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", username, password, dbName))
	if err != nil {
		fmt.Println(err)
		return db
	}

	// Perform a simple query to check the connection
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return db
	}

	return db
}

func BeginTx() *sql.Tx {
	db := Connect()
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return tx
	}

	defer db.Close()

	return tx
}

func NewDB() *ConnectDB {
	return &ConnectDB{tx: BeginTx()}
}
