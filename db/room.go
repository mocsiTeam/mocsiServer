package db

import (
	"errors"
	"strconv"

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

func (room *Rooms) AddUsers(db *gorm.DB, usersID []string, user *Users) error {
	if err := room.checkOwnerOrEditor(db, user); err != nil {
		return err
	}
	for _, id := range usersID {
		if err := db.Exec("SELECT * FROM room_accesses WHERE room_id = ? AND level_id = ? AND user_id = ?", room.ID, 3, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			userID, _ := strconv.Atoi(id)
			roomAccess := &RoomAccess{
				UserID:  uint(userID),
				RoomID:  room.ID,
				LevelID: 3,
			}
			db.Create(&roomAccess)
		}
	}
	return nil
}

func (room *Rooms) AddGroups(db *gorm.DB, groupsID []string, user *Users) error {
	if err := room.checkOwnerOrEditor(db, user); err != nil {
		return err
	}
	for _, id := range groupsID {
		var groupAccess []GroupAccess
		db.Joins("User").Where("group_id = ?", id).Find(&groupAccess)
		for _, id := range groupAccess {
			if err := db.Exec("SELECT * FROM room_accesses WHERE room_id = ? AND level_id = ? AND user_id = ?", room.ID, 3, id.User.ID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				roomAccess := &RoomAccess{
					UserID:  id.User.ID,
					RoomID:  room.ID,
					LevelID: 3,
				}
				db.Create(&roomAccess)
			}
		}
	}
	return nil
}

func (room *Rooms) KickUsers(db *gorm.DB, usersID []string, user *Users) error {
	if err := room.checkOwnerOrEditor(db, user); err != nil {
		return err
	}
	for _, id := range usersID {
		db.Exec("DELETE FROM room_accesses WHERE room_id = ? AND level_id = ? AND user_id = ?", room.ID, 3, id)
	}
	return nil
}

func (room *Rooms) KickGroups(db *gorm.DB, groupsID []string, user *Users) error {
	if err := room.checkOwnerOrEditor(db, user); err != nil {
		return err
	}
	for _, id := range groupsID {
		var groupAccess []GroupAccess
		db.Joins("User").Where("group_id = ?", id).Find(&groupAccess)
		for _, id := range groupAccess {
			db.Exec("DELETE FROM room_accesses WHERE room_id = ? AND level_id = ? AND user_id = ?", room.ID, 3, id.User.ID)
		}
	}
	return nil
}

func GetModRoom(db *gorm.DB, id string, user *Users) (*Rooms, error) {
	var roomAccess RoomAccess
	if err := db.Joins("Room").Where("user_id = ? AND level_id = ?", user.ID, 1).Or("user_id = ? AND level_id = ?", user.ID, 2).First(&roomAccess).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("group_not_found")
	}
	return &roomAccess.Room, nil
}

func (room *Rooms) DeleteRoom(db *gorm.DB, user *Users) error {
	if err := room.checkOwner(db, user); err != nil {
		return err
	} else if err := db.Delete(&room).Error; err != nil {
		return err
	} else if err := db.Exec("DELETE FROM room_accesses WHERE room_id = ?", room.ID).Error; err != nil {
		return nil
	}
	return nil
}

func (room *Rooms) checkOwner(db *gorm.DB, user *Users) error {
	var roomAccess RoomAccess
	if err := user.Check(db); err != nil {
		return &UserNotFound{}
	} else if err := db.Where("room_id = ? AND level_id = ? AND user_id = ?", room.ID, 1, user.ID).First(&roomAccess).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("not_owner")
	}
	return nil
}

func (room *Rooms) checkOwnerOrEditor(db *gorm.DB, user *Users) error {
	var roomAccess RoomAccess
	if err := user.Check(db); err != nil {
		return &UserNotFound{}
	} else if err := db.Where("room_id = ? AND level_id = ? AND user_id = ?", room.ID, 1, user.ID).Or("room_id = ? AND level_id = ? AND user_id = ?", room.ID, 2, user.ID).First(&roomAccess).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("not_owner_or_editor")
	}
	return nil
}
