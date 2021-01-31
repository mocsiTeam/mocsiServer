package db

import (
	"errors"
	"strconv"

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

func GetPublicGroups(db *gorm.DB, names []string) []*Groups {
	var groups []*Groups
	db.Where(map[string]interface{}{"name": names, "private": false}).Find(&groups)
	return groups
}

func GetPrivateGroups(db *gorm.DB, id []string, user *Users) []*Groups {
	var groupAccess []GroupAccess
	var groups []*Groups
	db.Joins("Group").Where(map[string]interface{}{"group_id": id, "user_id": user.ID, "level_id": 1}).Or(map[string]interface{}{"group_id": id, "user_id": user.ID, "level_id": 2}).Find(&groupAccess)
	for _, group := range groupAccess {
		if group.Group.Private {
			groups = append(groups, &group.Group)
		}
	}
	return groups
}

func GetMyGroups(db *gorm.DB, user *Users) []*Groups {
	var groupAccess []GroupAccess
	var groups []*Groups
	db.Joins("Group").Where("user_id = ?", user.ID).Find(&groupAccess)
	for _, group := range groupAccess {
		if (group.LevelID == 3 && !group.Group.Private) || group.LevelID > 3 {
			g := group.Group
			groups = append(groups, &g)
		}
	}
	return groups
}

func GetModGroup(db *gorm.DB, id string, user *Users) (*Groups, error) {
	var groupAccess GroupAccess
	if err := db.Joins("Group").Where("user_id = ? AND level_id = ?", user.ID, 1).Or("user_id = ? AND level_id = ?", user.ID, 2).First(&groupAccess).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("group_not_found")
	}
	return &groupAccess.Group, nil
}

func (group *Groups) GetOwner(db *gorm.DB) *Users {
	var groupAccess GroupAccess
	db.Joins("User").Where("group_id = ? AND level_id = ?", group.ID, 1).First(&groupAccess)
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

func (group *Groups) GetEditors(db *gorm.DB) []Users {
	var groupAccess []GroupAccess
	var users []Users
	db.Joins("User").Where("group_id = ? AND level_id = ?", group.ID, 2).First(&groupAccess)
	for _, user := range groupAccess {
		users = append(users, user.User)
	}
	return users
}

func (group *Groups) AddUsers(db *gorm.DB, usersID []string, user *Users) error {
	if err := group.checkOwnerOrEditor(db, user); err != nil {
		return err
	}
	for _, id := range usersID {
		if err := db.Exec("SELECT * FROM group_accesses WHERE group_id = ? AND level_id = ? AND user_id = ?", group.ID, 3, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			group.CountUsers++
			userID, _ := strconv.Atoi(id)
			groupAccess := &GroupAccess{
				UserID:  uint(userID),
				GroupID: group.ID,
				LevelID: 3,
			}
			db.Create(&groupAccess)
			db.Save(&group)
		}
	}
	return nil
}

func (group *Groups) AddEditors(db *gorm.DB, usersID []string, user *Users) error {
	if err := group.checkOwner(db, user); err != nil {
		return err
	}
	for _, id := range usersID {
		if err := db.Exec("SELECT * FROM group_accesses WHERE group_id = ? AND level_id = ? AND user_id = ?", group.ID, 2, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			group.CountUsers++
			userID, _ := strconv.Atoi(id)
			groupAccess := &GroupAccess{
				UserID:  uint(userID),
				GroupID: group.ID,
				LevelID: 3,
			}
			db.Create(&groupAccess)
			db.Save(&group)
		}
	}
	return nil
}

func (group *Groups) KickUsers(db *gorm.DB, usersID []string, user *Users) error {
	if err := group.checkOwnerOrEditor(db, user); err != nil {
		return err
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
	if err := group.checkOwner(db, user); err != nil {
		return err
	} else if err := db.Delete(&group).Error; err != nil {
		return err
	} else if err := db.Exec("DELETE FROM group_accesses WHERE group_id = ?", group.ID).Error; err != nil {
		return nil
	}
	return nil
}

func (group *Groups) checkOwner(db *gorm.DB, user *Users) error {
	var groupAccess GroupAccess
	if err := user.Check(db); err != nil {
		return &UserNotFound{}
	} else if err := db.Where("group_id = ? AND level_id = ? AND user_id = ?", group.ID, 1, user.ID).First(&groupAccess).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("not_owner")
	}
	return nil
}

func (group *Groups) checkOwnerOrEditor(db *gorm.DB, user *Users) error {
	var groupAccess GroupAccess
	if err := user.Check(db); err != nil {
		return &UserNotFound{}
	} else if err := db.Where("group_id = ? AND level_id = ? AND user_id = ?", group.ID, 1, user.ID).Or("group_id = ? AND level_id = ? AND user_id = ?", group.ID, 2, user.ID).First(&groupAccess).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("not_owner_or_editor")
	}
	return nil
}
