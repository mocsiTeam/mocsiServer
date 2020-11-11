package main

import (
	"log"

	"github.com/mocsiTeam/mocsiServer/db"
)

func init() {
	err := db.CreateDB()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {

}
