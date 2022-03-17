package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var (
	db *sql.DB
)

func main() {

	initDB()

	exampleInsertTx()
	exampleUpdateTx()
	exampleUpdateSelect()

}

func exampleUpdateSelect() {

	tx, err := db.Begin()

	if err != nil {
		log.Fatal("err to open tx", err)
	}

	sqlUpdate := `update books 
			set name = 'new value'`

	//update with tx
	_, errUpdate := tx.Exec(sqlUpdate)
	if errUpdate != nil {
		log.Fatal(errUpdate)
	}

	//select without tx
	qry := `select name, autor from books`
	rows, err := db.Query(qry)
	if err != nil {
		log.Fatal(err)
	}
	//defer rows.Close()

	for rows.Next() {
		var book Book
		errScan := rows.Scan(&book.Name, &book.Autor)
		if errScan != nil {
			log.Print("err query ", errScan)
		}
		log.Print("result book", book)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("err to make commit", err)
	}

}

func exampleInsertTx() {
	tx, err := db.Begin()

	if err != nil {
		log.Fatal("err to open tx", err)
	}

	sqlInsert := `insert into books(name, autor) 
			values (?, ?)`

	_, errInsert := tx.Exec(sqlInsert, "livro", "autor")
	if errInsert != nil {
		log.Fatal(errInsert)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("err to make commit", err)
	}

}

func exampleUpdateTx() {

	tx, err := db.Begin()

	if err != nil {
		log.Fatal("err to open tx", err)
	}

	sqlUpdate := `update books 
			set name = 'teste'`

	_, errUpdate := tx.Exec(sqlUpdate)
	if errUpdate != nil {
		if err := tx.Rollback(); err != nil {
			log.Fatal("err to make rollback", err)
		}
		log.Fatal(errUpdate)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("err to make commit", err)
	}

}

func initDB() {
	mysqlDB, err := sql.Open("mysql", "mysql:mysql@/bookings_db")
	if err != nil {
		panic(err)
	}
	err = mysqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	// See "Important settings" section.
	mysqlDB.SetConnMaxLifetime(time.Minute * 3)
	mysqlDB.SetMaxOpenConns(10)
	mysqlDB.SetMaxIdleConns(10)
	db = mysqlDB
}

type Book struct {
	Name  string
	Autor string
}
