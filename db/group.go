package db

import (
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type qgroup struct {
	groups        []*Groups
	groupAccesses []GroupAccess
	groupAccess   GroupAccess
	group         Groups
	users         []*Users
	user          Users
}

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
	var qg qgroup
	db.Where(map[string]interface{}{"name": names, "private": false}).Find(&qg.groups)
	return qg.groups
}

func GetPrivateGroups(db *gorm.DB, id []string, user *Users) []*Groups {
	var qg qgroup
	db.Joins("Group").Where(map[string]interface{}{"group_id": id, "user_id": user.ID, "level_id": 1}).Or(map[string]interface{}{"group_id": id, "user_id": user.ID, "level_id": 2}).Find(&qg.groupAccesses)
	for _, group := range qg.groupAccesses {
		if group.Group.Private {
			qg.groups = append(qg.groups, &group.Group)
		}
	}
	return qg.groups
}

func GetMyGroups(db *gorm.DB, user *Users) []*Groups {
	var qg qgroup
	db.Joins("Group").Where("user_id = ?", user.ID).Find(&qg.groupAccesses)
	for _, group := range qg.groupAccesses {
		if (group.LevelID == 3 && !group.Group.Private) || group.LevelID > 3 {
			g := group.Group
			qg.groups = append(qg.groups, &g)
		}
	}
	return qg.groups
}

func GetModGroup(db *gorm.DB, id string, user *Users) (*Groups, error) {
	var qg qgroup
	if err := db.Joins("Group").Where("user_id = ? AND level_id = ?", user.ID, 1).Or("user_id = ? AND level_id = ?", user.ID, 2).First(&qg.groupAccess).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("group_not_found")
	}
	return &qg.groupAccess.Group, nil
}

func (group *Groups) GetOwner(db *gorm.DB) *Users {
	var qg qgroup
	db.Joins("User").Where("group_id = ? AND level_id = ?", group.ID, 1).First(&qg.groupAccess)
	return &qg.groupAccess.User
}

func (group *Groups) GetUsers(db *gorm.DB) []*Users {
	var qg qgroup
	db.Joins("User").Where("group_id = ?", group.ID).Find(&qg.groupAccesses)
	for _, user := range qg.groupAccesses {
		qg.users = append(qg.users, &user.User)
	}
	return qg.users
}

func (group *Groups) GetEditors(db *gorm.DB) []*Users {
	var qg qgroup
	db.Joins("User").Where("group_id = ? AND level_id = ?", group.ID, 2).First(&qg.groupAccesses)
	for _, user := range qg.groupAccesses {
		qg.users = append(qg.users, &user.User)
	}
	return qg.users
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
	var qg qgroup
	if err := user.Check(db); err != nil {
		return &UserNotFound{}
	} else if err := db.Where("group_id = ? AND level_id = ? AND user_id = ?", group.ID, 1, user.ID).First(&qg.groupAccess).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("not_owner")
	}
	return nil
}

func (group *Groups) checkOwnerOrEditor(db *gorm.DB, user *Users) error {
	var qg qgroup
	if err := user.Check(db); err != nil {
		return &UserNotFound{}
	} else if err := db.Where("group_id = ? AND level_id = ? AND user_id = ?", group.ID, 1, user.ID).Or("group_id = ? AND level_id = ? AND user_id = ?", group.ID, 2, user.ID).First(&qg.groupAccess).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("not_owner_or_editor")
	}
	return nil
}
