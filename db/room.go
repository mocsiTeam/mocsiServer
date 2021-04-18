package db

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type qroom struct {
	users    []*Users
	rooms    []*Rooms
	accesses []RoomAccess
	owner    *Users
	room     Rooms
	access   RoomAccess
}

func (room *Rooms) Create(db *gorm.DB, user *Users) error {
	if err := db.Select("name").Where("name = ?", room.Name).First(&room).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("31")
	} else if err := user.Check(db); err != nil {
		return fmt.Errorf("32")
	} else if result := db.Create(&room); result.Error != nil {
		return fmt.Errorf("37")
	}
	var roomAccess = &RoomAccess{
		UserID:  user.ID,
		RoomID:  room.ID,
		LevelID: 1,
	}
	if err := db.Create(&roomAccess).Error; err != nil {
		return fmt.Errorf("38")
	}
	return nil
}

func (event *Events) Create(db *gorm.DB) error {
	if result := db.Create(&event); result.Error != nil {
		return result.Error
	}
	return nil
}

func GetEventsMonth(db *gorm.DB, datetime string) ([]Events, error) {
	var events []Events
	dt, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		return []Events{}, err
	}
	dt2 := dt.AddDate(0, 1, 0)
	db.Joins("Room").Where("date_time between ? and ?", dt, dt2).Find(&events)
	return events, nil
}

func (room *Rooms) AddUsers(db *gorm.DB, usersID []string, user *Users) error {
	if err := room.checkOwnerOrEditor(db, user); err != nil {
		return fmt.Errorf("33")
	}
	for _, id := range usersID {
		if err := db.Where("room_id = ? AND level_id = ? AND user_id = ?", room.ID, 3, id).First(&RoomAccess{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			userID, _ := strconv.Atoi(id)
			roomAccess := &RoomAccess{
				UserID:  uint(userID),
				RoomID:  room.ID,
				LevelID: 3,
			}
			db.Create(&roomAccess)
			db.Save(&room)
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
		return fmt.Errorf("33")
	}
	for _, id := range usersID {
		db.Exec("DELETE FROM room_accesses WHERE room_id = ? AND level_id = ? AND user_id = ?", room.ID, 3, id)
	}
	return nil
}

func (room *Rooms) KickGroups(db *gorm.DB, groupsID []string, user *Users) error {
	if err := room.checkOwnerOrEditor(db, user); err != nil {
		return fmt.Errorf("33")
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
	var qr qroom
	if err := db.Joins("Room").Where("user_id = ? AND level_id = ?", user.ID, 1).Or("user_id = ? AND level_id = ?", user.ID, 2).First(&qr.access).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("34")
	}
	return &qr.access.Room, nil
}

func (room *Rooms) DeleteRoom(db *gorm.DB, user *Users) error {
	if err := room.checkOwner(db, user); err != nil {
		return fmt.Errorf("33")
	} else if err := db.Delete(&room).Error; err != nil {
		return fmt.Errorf("35")
	} else if err := db.Exec("DELETE FROM room_accesses WHERE room_id = ?", room.ID).Error; err != nil {
		return fmt.Errorf("36")
	}
	return nil
}

func (room *Rooms) checkOwner(db *gorm.DB, user *Users) error {
	var qr qroom
	if err := user.Check(db); err != nil {
		return fmt.Errorf("32")
	} else if err := db.Where("room_id = ? AND level_id = ? AND user_id = ?", room.ID, 1, user.ID).First(&qr.access).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("33")
	}
	return nil
}

func (room *Rooms) checkOwnerOrEditor(db *gorm.DB, user *Users) error {
	var qr qroom
	if err := user.Check(db); err != nil {
		return fmt.Errorf("32")
	} else if err := db.Where("room_id = ? AND level_id = ? AND user_id = ?", room.ID, 1, user.ID).Or("room_id = ? AND level_id = ? AND user_id = ?", room.ID, 2, user.ID).First(&qr.access).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("33")
	}
	return nil
}

func GetMyRooms(db *gorm.DB, user *Users) []*Rooms {
	var qr qroom
	db.Joins("Room").Where("user_id = ?", user.ID).Find(&qr.accesses)
	for _, room := range qr.accesses {
		tmp := room.Room
		qr.rooms = append(qr.rooms, &tmp)
	}
	return qr.rooms
}

func GetRooms(db *gorm.DB, id []string, user *Users) []*Rooms {
	var qr qroom
	db.Joins("Room").Where(map[string]interface{}{"room_id": id, "user_id": user.ID}).Find(&qr.accesses)
	for _, room := range qr.accesses {
		qr.rooms = append(qr.rooms, &room.Room)
	}
	return qr.rooms
}

func (room *Rooms) GetUsers(db *gorm.DB) []*Users {
	var qr qroom
	db.Joins("User").Where("room_id = ?", room.ID).Find(&qr.accesses)
	for _, user := range qr.accesses {
		tmp := user.User
		qr.users = append(qr.users, &tmp)
	}
	return qr.users
}

func (room *Rooms) GetOwner(db *gorm.DB) *Users {
	var qr qroom
	db.Joins("User").Where("room_id = ? AND level_id = ?", room.ID, 1).Find(&qr.access)
	return &qr.access.User
}

func (room *Rooms) GetEditors(db *gorm.DB) []*Users {
	var qr qroom
	db.Joins("User").Where("room_id = ? AND level_id = ?", room.ID, 2).Find(&qr.accesses)
	for _, editor := range qr.accesses {
		qr.users = append(qr.users, &editor.User)
	}
	return qr.users
}

func (room *Rooms) AddEditors(db *gorm.DB, usersID []string, user *Users) error {
	if err := room.checkOwner(db, user); err != nil {
		return fmt.Errorf("33")
	}
	for _, id := range usersID {
		if err := db.Exec("SELECT * FROM room_accesses WHERE room_id = ? AND level_id = ? AND user_id = ?", room.ID, 2, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			// TODO room.CountUsers++
			userID, _ := strconv.Atoi(id)
			roomAccess := &RoomAccess{
				UserID:  uint(userID),
				RoomID:  room.ID,
				LevelID: 3,
			}
			db.Create(&roomAccess)
			db.Save(&room)
		}
	}
	return nil
}

func (room *Rooms) KickEditors(db *gorm.DB, usersID []string, user *Users) error {
	if err := room.checkOwner(db, user); err != nil {
		return fmt.Errorf("33")
	}
	for _, id := range usersID {
		idint, _ := strconv.Atoi(id)
		ra := RoomAccess{UserID: uint(idint), RoomID: room.ID, LevelID: 3}
		db.Save(&ra)
	}
	return nil
}
