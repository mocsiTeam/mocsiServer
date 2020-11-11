package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type database struct {
	name string
}

const (
	user    = "postgres"
	pass    = "123"
	host    = "localhost"
	port    = "5433"
	sslmode = "disable"
	dbname  = "mocsidb"
)

func connectPG() (*sql.DB, error) {
	conninfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s sslmode=%s",
		host, port, user, pass, sslmode)
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		return db, err
	}
	return db, nil
}

func checkDB(db *sql.DB) (bool, error) {
	rows, err := db.Query(`SELECT datname FROM pg_database;`)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	databases := []database{}

	for rows.Next() {
		d := database{}
		err := rows.Scan(&d.name)
		if err != nil {
			log.Println(err)
			continue
		}
		databases = append(databases, d)
	}
	for _, d := range databases {
		if d.name == dbname {
			return true, nil
		}
	}
	return false, nil
}

func CreateDB() error {
	db, err := connectPG()
	if err != nil {
		return err
	}
	if ok, err := checkDB(db); ok && err != nil {
		//TODO: Write code to create db and tables
	} else {
		return err
	}
	return nil
}
