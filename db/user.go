package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connector() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s sslmode=%s dbname=%s",
		host, port, user, pass, sslmode, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil
	}
	return db
}

func (user *Users) Create(db *gorm.DB) (*gorm.DB, error) {
	if user.Check(db) == nil {
		return nil, &NameAlredyExists{}
	}
	//TODO: check email
	return db.Create(&user), nil
}

func (user *Users) GetAll(db *gorm.DB) []Users {
	users := []Users{}
	db.Find(&users)
	return users
}

func (user *Users) Check(db *gorm.DB) error {
	result := db.Select("ID", "nick_name").Find(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (user *Users) Authenticate(db *gorm.DB) bool {
	pass := user.Pass
	result := db.Select("Pass").Find(&user)
	if result.Error != nil {
		log.Println(result.Error)
		return false
	}
	if pass != user.Pass {
		return false
	}
	return true
}

type WrongUsernameOrPasswordError struct{}

func (m *WrongUsernameOrPasswordError) Error() string {
	return "wrong username or password"
}

type NameAlredyExists struct{}

func (m *NameAlredyExists) Error() string {
	return "name alredy exists"
}
