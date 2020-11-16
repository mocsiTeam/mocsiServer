package db

import (
	"time"

	"gorm.io/gorm"
)

type Roles struct {
	ID   uint   `gorm: primaryKey`
	Name string `gorm: unique;not null`
}

type Groups struct {
	gorm.Model
	Name string `gorm: unique;not null`
}

type Rooms struct {
	gorm.Model
	Name  string `gorm: unique;not null`
	Link  string `gorm: unique;not null`
	Pass  []byte
	Timer time.Time
}

type AccessLevel struct {
	ID   uint   `gorm: primaryKey`
	Name string `gorm: unique;not null`
}

type Users struct {
	gorm.Model
	NickName   string  `gorm: unique;not null`
	Email      *string `gorm: unique;not null`
	Pass       []byte  `gorm: not null`
	Name       string
	SurName    string
	SecondName string
	Role       Roles `gorm:"foreignKey:Name;not null"`
}

type UserGroups struct {
	ID    uint   `gorm: primaryKey`
	User  Users  `gorm:"foreignKey:NickName;not null"`
	Group Groups `gorm:"foreignKey:Name;not null"`
}

type RoomAccess struct {
	ID    uint        `gorm: primaryKey`
	User  Users       `gorm:"foreignKey:NickName;not null"`
	Room  Rooms       `gorm:"foreignKey:Name;not null"`
	Level AccessLevel `gorm:"foreignKey:Name;not null"`
}

type StatsRoom struct {
	ID         uint  `gorm: primaryKey`
	Room       Rooms `gorm:"foreignKey:Name;not null"`
	CountUsers int
}

type StatsUser struct {
	ID        uint `gorm: primaryKey`
	CreatedAt time.Time
	Room      Rooms `gorm:"foreignKey:Name;not null"`
	TimeSpend time.Time
}
