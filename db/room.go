package db

import (
	"errors"

	"gorm.io/gorm"
)

func (room *Rooms) Create(db *gorm.DB, user *Users) error {
	if err := db.Select("name").Where("name = ?", room.Name).First(&room).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return &NameAlredyExists{}
	} else if err := user.Check(db); err != nil {
		return &UserNotFound{}
	} else if result := db.Create(&room); result.Error != nil {
		return err
	}
	var roomAccess = &RoomAccess{
		UserID:  user.ID,
		RoomID:  room.ID,
		LevelID: 1,
	}
	if err := db.Create(&roomAccess).Error; err != nil {
		return err
	}
	return nil
}
