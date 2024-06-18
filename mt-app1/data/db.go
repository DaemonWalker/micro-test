package data

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDb() {
	dsn := "root:root@tcp(127.0.0.1:3306)/boss"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(100)
}

func Begin() (*sql.Tx, error) {
	return db.Begin()
}

func query[T any](scaner func(*sql.Rows) T, sql string, args ...any) ([]T, error) {
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		log.Println("rows closed")
		rows.Close()
	}()
	arr := make([]T, 0)
	for rows.Next() {
		arr = append(arr, scaner(rows))
	}

	return arr, nil
}

func queryRow[T any](scaner func(*sql.Row) (T, error), sql string, args ...any) (T, error) {
	rows := db.QueryRow(sql, args...)
	return scaner(rows)
}

func exec(executor executor, sql string, args ...any) {
	executor.Exec(sql, args...)
}

type executor interface {
	Exec(query string, args ...any) (sql.Result, error)
}
