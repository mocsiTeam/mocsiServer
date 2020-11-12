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

func connectPGinDB() (*sql.DB, error) {
	conninfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s sslmode=%s dbname=%s",
		host, port, user, pass, sslmode, dbname)
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
	pg, err := connectPG()
	defer pg.Close()
	if err != nil {
		return err
	}
	ok, err := checkDB(pg)
	if err != nil {
		return err
	}
	if !ok {
		_, err = pg.Exec("CREATE DATABASE " + dbname)
		if err != nil {
			return err
		}
		pg.Close()
		db, err := connectPGinDB()
		if err != nil {
			return err
		}
		defer db.Close()
		execsDB := []string{createTableRoles, createTableGroups,
			createTableRoom, createTableAccessLevel,
			creataTableUsers, createTableRoomAccess,
			creataTableStatsRoom, createTableStatsUser}
		for _, exec := range execsDB {
			_, err = db.Exec(exec)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
	return nil
}
