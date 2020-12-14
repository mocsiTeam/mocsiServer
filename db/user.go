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
	if err := db.Select("nickname").Where("nickname = ?", user.Nickname).First(&user).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
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
	err := db.Select("ID", "nickname").Where("nickname = ?", user.Nickname).First(&user).Error
	return err
}

func (user *Users) Authenticate(db *gorm.DB) bool {
	pass := user.Pass
	if result := db.Where("nickname = ?", user.Nickname).First(&user); result.Error != nil {
		log.Println(result.Error)
		return false
	}
	return CheckPasswordHash(pass, user.Pass)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (user *Users) GetRefreshToken(db *gorm.DB, id string) (string, error) {
	if err := db.Select("refresh_token", "nickname").Where("id = ?", id).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return "", &UserNotFound{}
	}
	return user.RefreshToken, nil
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

type UserNotFound struct{}

func (m *UserNotFound) Error() string {
	return "user not found"
}
