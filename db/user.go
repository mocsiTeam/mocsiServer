package db

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
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

func (user *Users) Create(db *gorm.DB) error {
	var err error
	user.Pass, err = HashPassword(user.Pass)
	if err != nil {
		return err
	}
	if err := db.Select("nick_name").Where("nick_name = ?", user.NickName).First(&user).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return &NameAlredyExists{}
	} else if err := db.Select("email").Where("email = ?", user.Email).First(&user).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return &EmailAlredyExists{}
	} else if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (user *Users) GetAll(db *gorm.DB) []Users {
	users := []Users{}
	db.Find(&users)
	return users
}

func (user *Users) Check(db *gorm.DB) error {
	err := db.Select("ID", "nick_name").Where("nick_name = ?", user.NickName).First(&user).Error
	return err
}
func (user *Users) Authenticate(db *gorm.DB) bool {
	pass := user.Pass
	if result := db.Select("Pass").Where("nick_name = ?", user.NickName).First(&user); result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return CheckPasswordHash(pass, user.Pass)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type WrongUsernameOrPasswordError struct{}

func (m *WrongUsernameOrPasswordError) Error() string {
	return "wrong username or password"
}

type NameAlredyExists struct{}

func (m *NameAlredyExists) Error() string {
	return "name alredy exists"
}

type EmailAlredyExists struct{}

func (m *EmailAlredyExists) Error() string {
	return "email alredy exists"
}
