package db

import (
	"time"

	"gorm.io/gorm"
)

type Roles struct {
	ID       uint   `gorm:"primaryKey;"`
	RoleName string `gorm:"unique; not null;"`
}

type Groups struct {
	gorm.Model
	Name       string `gorm:"unique; not null;"`
	CountUsers uint
	Private    bool
}

type Rooms struct {
	gorm.Model
	Name  string `gorm:"unique; not null;"`
	Link  string `gorm:"unique; not null;"`
	Pass  string
	Timer time.Time
}

type AccessLevel struct {
	ID   uint   `gorm:"primaryKey;"`
	Name string `gorm:"unique; not null;"`
}

type Users struct {
	gorm.Model
	Nickname     string `gorm:"unique; not null;" json:"name"`
	Email        string `gorm:"unique; not null;"`
	Pass         string `gorm:"not null;"`
	Firstname    string
	Lastname     string
	RefreshToken string
	RoleID       uint
	Role         Roles `gorm:"foreignKey:RoleID;"`
}

type GroupAccess struct {
	ID      uint `gorm:"primaryKey;"`
	UserID  uint
	GroupID uint
	LevelID uint
	User    Users       `gorm:"foreignKey:UserID; not null;"`
	Group   Groups      `gorm:"foreignKey:GroupID; not null;"`
	Level   AccessLevel `gorm:"foreignKey:LevelID; not null;"`
}

type RoomAccess struct {
	ID      uint `gorm:"primaryKey;"`
	UserID  uint
	RoomID  uint
	LevelID uint
	User    Users       `gorm:"foreignKey:UserID; not null;"`
	Room    Rooms       `gorm:"foreignKey:RoomID; not null;"`
	Level   AccessLevel `gorm:"foreignKey:LevelID; not null;"`
}

type StatsRoom struct {
	ID         uint `gorm:"primaryKey;"`
	RoomID     uint
	Room       Rooms `gorm:"foreignKey:RoomID; not null;"`
	CountUsers int
}

type StatsUser struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	RoomID    uint
	Room      Rooms `gorm:"foreignKey:RoomID; not null;"`
	UserID    uint
	User      Users `gorm:"foreignKey:UserID; not null;"`
	TimeSpend time.Time
}
