package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/mocsiTeam/mocsiServer/api/graph/generated"
	"github.com/mocsiTeam/mocsiServer/api/graph/model"
	"github.com/mocsiTeam/mocsiServer/auth"
	"github.com/mocsiTeam/mocsiServer/auth/jwt"
	"github.com/mocsiTeam/mocsiServer/db"
	"gorm.io/gorm"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (mod *model.Tokens, err error) {
	defer getReport(&err)
	mod = &model.Tokens{}
	var newUser = db.Users{
		Email:     input.Email,
		Nickname:  input.Nickname,
		Firstname: input.Firstname,
		Lastname:  input.Lastname,
		Pass:      input.Password,
		RoleID:    3,
	}
	err = newUser.Create(DB)
	panicIf(err)
	accessToken, err := jwt.GenerateAccessToken(newUser.Nickname, strconv.Itoa(int(newUser.ID)))
	panicIf(err)
	refreshToken, err := jwt.GenerateRefreshToken(newUser.Nickname, newUser.Email, strconv.Itoa(int(newUser.ID)))
	panicIf(err)
	newUser.RefreshToken = refreshToken
	DB.Save(&newUser)
	mod.AccessToken = accessToken
	mod.RefreshToken = refreshToken
	return mod, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (model *model.Tokens, err error) {
	defer getReport(&err)
	var user = db.Users{
		Nickname: input.Nickname,
		Pass:     input.Password,
	}
	if correct := user.Authenticate(DB); !correct {
		panic("wrong username or password")
	}
	userID := strconv.Itoa(int(user.ID))
	accessToken, err := jwt.GenerateAccessToken(user.Nickname, userID)
	panicIf(err)
	refreshToken, err := jwt.GenerateRefreshToken(user.Nickname, user.Email, userID)
	panicIf(err)
	user.RefreshToken = refreshToken
	DB.Save(&user)
	model.AccessToken = accessToken
	model.RefreshToken = refreshToken
	return
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (token string, err error) {
	defer getReport(&err)
	var user db.Users
	userID, err := jwt.ParseToken(input.Token)
	panicIf(err)
	refreshToken, err := user.GetRefreshToken(DB, userID)
	panicIf(err)
	if refreshToken != input.Token {
		panic("invalid refresh token")
	}
	token, err = jwt.GenerateAccessToken(user.Nickname, userID)
	panicIf(err)
	return
}

func (r *mutationResolver) CreateGroup(ctx context.Context, input model.NewGroup) (mod *model.Group, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	group := &db.Groups{Name: input.Name, Private: input.Private}
	err = group.Create(DB, user)
	panicIf(err)
	owner := &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.Nickname,
		Firstname: user.Firstname, Lastname: user.Lastname,
		Email: user.Email, Role: strconv.Itoa(int(user.RoleID))}
	mod = &model.Group{ID: strconv.Itoa(int(group.ID)), Name: group.Name,
		CountUsers: int(group.CountUsers), Owner: owner, Users: []*model.User{owner}}
	return
}

func (r *mutationResolver) AddUsersToGroup(ctx context.Context, input model.UsersToGroup) (status string, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	group, err := db.GetModGroup(DB, input.GroupID, user)
	panicIf(err)
	err = group.AddUsers(DB, input.UsersID, user)
	panicIf(err)
	return "users_added", nil
}

func (r *mutationResolver) AddEditorsToGroup(ctx context.Context, input model.UsersToGroup) (status string, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	group, err := db.GetModGroup(DB, input.GroupID, user)
	panicIf(err)
	err = group.AddEditors(DB, input.UsersID, user)
	panicIf(err)
	return "users_became_editors", nil
}

func (r *mutationResolver) KickUsersFromGroup(ctx context.Context, input model.UsersToGroup) (status string, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	group, err := db.GetModGroup(DB, input.GroupID, user)
	panicIf(err)
	err = group.KickUsers(DB, input.UsersID, user)
	panicIf(err)
	return "users_kicked", nil
}

func (r *mutationResolver) DeleteGroup(ctx context.Context, input string) (status string, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	group, err := db.GetModGroup(DB, input, user)
	panicIf(err)
	err = group.DeleteGroup(DB, user)
	panicIf(err)
	return "group_deleted", nil
}

func (r *mutationResolver) CreateRoom(ctx context.Context, input model.NewRoom) (mod *model.Room, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	hostname, _ := os.Hostname()
	room := &db.Rooms{
		Name: input.Name,
		Link: "https://" + hostname + "/" + input.Name,
		Pass: input.Password,
	}
	err = room.Create(DB, user)
	panicIf(err)
	owner := &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.Nickname,
		Firstname: user.Firstname, Lastname: user.Lastname,
		Email: user.Email, Role: strconv.Itoa(int(user.RoleID))}
	return &model.Room{ID: strconv.Itoa(int(room.ID)), Name: room.Name,
		Link: room.Link, Owner: owner, Users: []*model.User{owner}}, nil
}

func (r *mutationResolver) AddUsersToRoom(ctx context.Context, input *model.UsersToRoom) (status string, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	room, err := db.GetModRoom(DB, input.RoomID, user)
	panicIf(err)
	err = room.AddUsers(DB, input.UsersID, user)
	panicIf(err)
	return "users_added", nil
}

func (r *mutationResolver) AddGroupToRoom(ctx context.Context, input *model.GroupsToRoom) (status string, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	room, err := db.GetModRoom(DB, input.RoomID, user)
	panicIf(err)
	err = room.AddGroups(DB, input.GroupsID, user)
	panicIf(err)
	return "groups_added", nil
}

func (r *mutationResolver) KickUsersFromRoom(ctx context.Context, input *model.UsersToRoom) (status string, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	room, err := db.GetModRoom(DB, input.RoomID, user)
	panicIf(err)
	err = room.KickUsers(DB, input.UsersID, user)
	panicIf(err)
	return "users_kicked", nil
}

func (r *mutationResolver) KickGroupsFromRoom(ctx context.Context, input *model.GroupsToRoom) (status string, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	room, err := db.GetModRoom(DB, input.RoomID, user)
	panicIf(err)
	err = room.KickGroups(DB, input.GroupsID, user)
	panicIf(err)
	return "groups_kicked", nil
}

func (r *mutationResolver) DeleteRoom(ctx context.Context, input string) (status string, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	group, err := db.GetModRoom(DB, input, user)
	panicIf(err)
	err = group.DeleteRoom(DB, user)
	panicIf(err)
	return "room_deleted", nil
}

func (r *queryResolver) GetAuthUser(ctx context.Context) (mod *model.User, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	err = user.Get(DB)
	panicIf(err)
	return &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.Nickname,
		Firstname: user.Firstname, Lastname: user.Lastname,
		Email: user.Email, Role: strconv.Itoa(int(user.RoleID))}, nil
}

func (r *queryResolver) GetAllUsers(ctx context.Context) (mod []*model.User, err error) {
	defer getReport(&err)
	if us := auth.ForContext(ctx); us == nil {
		panic("access denied")
	}
	var (
		users    db.Users
		allUsers []*model.User
	)
	for _, user := range users.GetAll(DB) {
		allUsers = append(allUsers, &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.Nickname,
			Firstname: user.Firstname, Lastname: user.Lastname,
			Email: user.Email, Role: strconv.Itoa(int(user.RoleID))})
	}
	return allUsers, nil
}

func (r *queryResolver) GetUsers(ctx context.Context, input []string) (mod []*model.User, err error) {
	defer getReport(&err)
	if us := auth.ForContext(ctx); us == nil {
		panic("access denied")
	}
	var (
		users        db.Users
		gettingUsers []*model.User
	)
	usersMap := make(map[string]bool)
	for _, user := range input {
		usersMap[user] = false
	}
	for _, user := range users.GetUsers(DB, input) {
		usersMap[user.Nickname] = false
		gettingUsers = append(gettingUsers, &model.User{ID: strconv.Itoa(int(user.ID)),
			Nickname:  user.Nickname,
			Firstname: user.Firstname, Lastname: user.Lastname,
			Email: user.Email, Role: strconv.Itoa(int(user.RoleID))})
	}
	for _, user := range input {
		if exist := usersMap[user]; !exist {
			gettingUsers = append(gettingUsers, &model.User{ID: "0", Nickname: user,
				Firstname: "", Lastname: "",
				Email: "", Role: "0", Error: "user_not_found"})
		}
	}
	return gettingUsers, nil
}

func (r *queryResolver) GetGroups(ctx context.Context, input model.InfoGroups) (mod []*model.Group, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	if input.IsPrivate {
		return getGroups(db.GetPrivateGroups(DB, input.GroupsID, user), input.GroupsID, input.IsPrivate), nil
	}
	return getGroups(db.GetPublicGroups(DB, input.GroupsID), input.GroupsID, input.IsPrivate), nil
}

func (r *queryResolver) GetMyGroups(ctx context.Context) (mod []*model.Group, err error) {
	defer getReport(&err)
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("access denied")
	}
	return getGroups(db.GetMyGroups(DB, user), []string{}, true), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
var DB *gorm.DB = db.Connector()

func getGroups(groups []*db.Groups, groupsID []string, isPrivate bool) []*model.Group {
	var gettingGroups []*model.Group
	groupsMap := make(map[string]bool)
	for _, v := range groupsID {
		groupsMap[v] = false
	}
	for _, group := range groups {
		groupsMap[getNameOrIDFromGroup(group, isPrivate)] = true
		owner := group.GetOwner(DB)
		var gettingUsers []*model.User
		for _, user := range group.GetUsers(DB) {
			gettingUsers = append(gettingUsers, &model.User{ID: strconv.Itoa(int(user.ID)),
				Nickname:  user.Nickname,
				Firstname: user.Firstname, Lastname: user.Lastname,
				Email: user.Email, Role: strconv.Itoa(int(user.RoleID))})
		}
		gettingGroups = append(gettingGroups, &model.Group{ID: strconv.Itoa(int(group.ID)),
			Name: group.Name, CountUsers: int(group.CountUsers),
			Owner: &model.User{ID: strconv.Itoa(int(owner.ID)),
				Nickname: owner.Nickname, Email: owner.Email},
			Users: gettingUsers})
	}
	if len(groups) != len(groupsID) {
		for _, v := range groupsID {
			if exist := groupsMap[v]; !exist {
				gettingGroups = append(gettingGroups, &model.Group{ID: "0", Name: v, CountUsers: 0,
					Error: "group_not_found"})
			}
		}
	}
	return gettingGroups
}

func getNameOrIDFromGroup(group *db.Groups, isPrivate bool) string {
	if isPrivate {
		return strconv.Itoa(int(group.ID))
	}
	return group.Name
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func getReport(err *error) {
	if r := recover(); r != nil {
		fmt.Println("panica")
		switch x := r.(type) {
		case string:
			*err = errors.New(x)
		case error:
			*err = x
		default:
			fmt.Println(x)
			*err = errors.New("Unknown error")
		}
		fmt.Println(err)
	}
}
