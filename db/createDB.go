package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	if err != nil {
		return err
	}
	defer func() error {
		err := pg.Close()
		if err != nil {
			return err
		}
		return nil
	}()
	ok, err := checkDB(pg)
	if err != nil {
		return err
	}
	if !ok {
		_, err = pg.Exec("CREATE DATABASE " + dbname)
		if err != nil {
			return err
		}
		err = pg.Close()
		if err != nil {
			return err
		}
		dsn := fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s sslmode=%s dbname=%s",
			host, port, user, pass, sslmode, dbname)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		db.AutoMigrate(&Roles{})
		db.AutoMigrate(&Groups{})
		db.AutoMigrate(&AccessLevel{})
		db.AutoMigrate(&Users{})
		db.AutoMigrate(&UserGroups{})
		db.AutoMigrate(&RoomAccess{})
		db.AutoMigrate(&StatsRoom{})
		db.AutoMigrate(&StatsUser{})
		for _, role := range roles {
			db.Create(&role)
		}
		for _, alevel := range alevels {
			db.Create(&alevel)
		}

	}
	return nil
}
