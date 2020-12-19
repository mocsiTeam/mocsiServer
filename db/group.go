package db

import (
	"errors"

	"gorm.io/gorm"
)

func (group *Groups) Create(db *gorm.DB, user *Users) error {
	group.CountUsers = 1
	if err := db.Select("name").Where("name = ?", group.Name).First(&group).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return &NameAlredyExists{}
	} else if err := user.Check(db); err != nil {
		return &UserNotFound{}
	} else if result := db.Create(&group); result.Error != nil {
		return err
	}
	var groupAccess = &GroupAccess{
		UserID:  user.ID,
		LevelID: 1,
		GroupID: group.ID,
	}
	if err := db.Create(&groupAccess).Error; err != nil {
		return err
	}
	return nil
}

func GetGroups(db *gorm.DB, names []string) []*Groups {
	//var group *Groups
	var groups []*Groups
	db.Where(map[string]interface{}{"name": names}).Find(&groups)
	return groups
}

func (group *Groups) GetOwner(db *gorm.DB) *Users {
	var groupAccess GroupAccess
	db.Joins("User").Where("group_id = ? AND level_id = ?", int(group.ID), 1).First(&groupAccess)
	return &groupAccess.User
}

func (group *Groups) GetUsers(db *gorm.DB) []*Users {
	var groupAccess []*GroupAccess
	var users []*Users
	db.Joins("User").Where("group_id = ?", group.ID).Find(&groupAccess)
	for _, user := range groupAccess {
		users = append(users, &user.User)
	}
	return users
}

func (group *Groups) AddUsers(db *gorm.DB, usersID []int, user *Users) error {
	var groupAccess GroupAccess
	if err := user.Check(db); err != nil {
		return &UserNotFound{}
	} else if err := db.Where("group_id = ? AND level_id = ? AND user_id = ?", group.ID, 1, user.ID).First(&groupAccess).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("not_owner")
	}
	for _, id := range usersID {
		if err := db.Exec("SELECT * FROM group_accesses WHERE group_id = ? AND level_id = ? AND user_id = ?", group.ID, 3, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			group.CountUsers++
			groupAccess := &GroupAccess{
				UserID:  uint(id),
				GroupID: group.ID,
				LevelID: 3,
			}
			db.Create(&groupAccess)
			db.Save(&group)
		}
	}
	return nil
}

func (group *Groups) KickUsers(db *gorm.DB, usersID []int, user *Users) error {
	var groupAccess GroupAccess
	if err := user.Check(db); err != nil {
		return &UserNotFound{}
	} else if err := db.Where("group_id = ? AND level_id = ? AND user_id = ?", group.ID, 1, user.ID).First(&groupAccess).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("not_owner")
	}
	for _, id := range usersID {
		if err := db.Exec("DELETE FROM group_accesses WHERE group_id = ? AND level_id = ? AND user_id = ?", group.ID, 3, id).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
			group.CountUsers--
			db.Save(&group)
		}
	}
	return nil
}

func (group *Groups) DeleteGroup(db *gorm.DB, user *Users) error {
	var groupAccess GroupAccess
	if err := user.Check(db); err != nil {
		return &UserNotFound{}
	} else if err := db.Where("group_id = ? AND level_id = ? AND user_id = ?", group.ID, 1, user.ID).First(&groupAccess).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("not_owner")
	} else if err := db.Delete(&group).Error; err != nil {
		return err
	} else if err := db.Where("group_id = ?", group.ID).Delete(&groupAccess).Error; err != nil {
		return nil
	}
	return nil
}
