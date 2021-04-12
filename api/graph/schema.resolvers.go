package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mocsiTeam/mocsiServer/api/graph/generated"
	"github.com/mocsiTeam/mocsiServer/api/graph/model"
	"github.com/mocsiTeam/mocsiServer/auth"
	"github.com/mocsiTeam/mocsiServer/auth/jwt"
	"github.com/mocsiTeam/mocsiServer/db"
	"gorm.io/gorm"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.Tokens, error) {
	var newUser = db.Users{
		Email:     input.Email,
		Nickname:  input.Nickname,
		Firstname: input.Firstname,
		Lastname:  input.Lastname,
		Pass:      input.Password,
		RoleID:    3,
	}
	err := newUser.Create(DB)
	panicIf(err)
	accessToken, err := jwt.GenerateAccessToken(newUser.Nickname, strconv.Itoa(int(newUser.ID)))
	if err != nil {
		panic("0")
	}
	refreshToken, err := jwt.GenerateRefreshToken(newUser.Nickname, newUser.Email, strconv.Itoa(int(newUser.ID)))
	if err != nil {
		panic("0")
	}
	newUser.RefreshToken = refreshToken
	DB.Save(&newUser)
	return &model.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (*model.Tokens, error) {
	var user = db.Users{
		Nickname: input.Nickname,
		Pass:     input.Password,
	}
	if correct := user.Authenticate(DB); !correct {
		panic("21")
	}
	userID := strconv.Itoa(int(user.ID))
	accessToken, err := jwt.GenerateAccessToken(user.Nickname, userID)
	if err != nil {
		panic("0")
	}
	refreshToken, err := jwt.GenerateRefreshToken(user.Nickname, user.Email, userID)
	if err != nil {
		panic("0")
	}
	user.RefreshToken = refreshToken
	DB.Save(&user)
	return &model.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, token model.RefreshTokenInput) (string, error) {
	var user db.Users
	userID, err := jwt.ParseToken(token.Token)
	if err != nil {
		panic("0")
	}
	refreshToken, err := user.GetRefreshToken(DB, userID)
	panicIf(err)
	if refreshToken != token.Token {
		panic("23")
	}
	atoken, err := jwt.GenerateAccessToken(user.Nickname, userID)
	if err != nil {
		panic("0")
	}
	return atoken, nil
}

func (r *mutationResolver) CreateGroup(ctx context.Context, input model.NewGroup) (*model.Group, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	group := &db.Groups{Name: input.Name, Private: input.Private}
	err := group.Create(DB, user)
	panicIf(err)
	owner := &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.Nickname,
		Firstname: user.Firstname, Lastname: user.Lastname,
		Email: user.Email, Role: strconv.Itoa(int(user.RoleID))}
	return &model.Group{ID: strconv.Itoa(int(group.ID)), Name: group.Name,
		CountUsers: int(group.CountUsers), Owner: owner, Users: []*model.User{owner}, Editors: []*model.User{owner}}, nil
}

func (r *mutationResolver) AddUsersToGroup(ctx context.Context, input model.UsersToGroup) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	group, err := db.GetModGroup(DB, input.GroupID, user)
	panicIf(err)
	err = group.AddUsers(DB, input.UsersID, user)
	panicIf(err)
	return "users_added", nil
}

func (r *mutationResolver) AddEditorsToGroup(ctx context.Context, input model.UsersToGroup) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	group, err := db.GetModGroup(DB, input.GroupID, user)
	panicIf(err)
	err = group.AddEditors(DB, input.UsersID, user)
	panicIf(err)
	return "users_became_editors", nil
}

func (r *mutationResolver) KickUsersFromGroup(ctx context.Context, input model.UsersToGroup) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	group, err := db.GetModGroup(DB, input.GroupID, user)
	panicIf(err)
	err = group.KickUsers(DB, input.UsersID, user)
	panicIf(err)
	return "users_kicked", nil
}

func (r *mutationResolver) DeleteGroup(ctx context.Context, id string) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	group, err := db.GetModGroup(DB, id, user)
	panicIf(err)
	err = group.DeleteGroup(DB, user)
	panicIf(err)
	return "group_deleted", nil
}

func (r *mutationResolver) CreateRoom(ctx context.Context, input model.NewRoom) (*model.Room, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	hostname, _ := os.Hostname()
	room := &db.Rooms{
		Name:       input.Name,
		UniqueName: input.UniqueName,
		Link:       hostname + "/" + input.UniqueName,
		Pass:       input.Password,
	}
	err := room.Create(DB, user)
	panicIf(err)
	owner := &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.Nickname,
		Firstname: user.Firstname, Lastname: user.Lastname,
		Email: user.Email, Role: strconv.Itoa(int(user.RoleID))}
	return &model.Room{ID: strconv.Itoa(int(room.ID)), Name: room.Name,
		Link: room.Link, Owner: owner, Users: []*model.User{owner}, Editors: []*model.User{owner}}, nil
}

func (r *mutationResolver) CreateEvent(ctx context.Context, input *model.NewEvent) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	dt, _ := time.Parse(time.RFC3339, input.Datetime)
	if dt.Before(time.Now()) {
		panic("44")
	}
	room := db.GetRooms(DB, []string{input.IDRoom}, user)
	event := &db.Events{
		DateTime: dt,
		Room:     *room[0],
	}
	err := event.Create(DB)
	panicIf(err)
	return "event_created", nil
}

func (r *mutationResolver) AddUsersToRoom(ctx context.Context, input model.UsersToRoom) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	room, err := db.GetModRoom(DB, input.RoomID, user)
	panicIf(err)
	err = room.AddUsers(DB, input.UsersID, user)
	panicIf(err)
	return "users_added", nil
}

func (r *mutationResolver) AddGroupToRoom(ctx context.Context, input model.GroupsToRoom) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	room, err := db.GetModRoom(DB, input.RoomID, user)
	panicIf(err)
	err = room.AddGroups(DB, input.GroupsID, user)
	panicIf(err)
	return "groups_added", nil
}

func (r *mutationResolver) AddEditorsToRoom(ctx context.Context, input model.UsersToRoom) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	room, err := db.GetModRoom(DB, input.RoomID, user)
	panicIf(err)
	err = room.AddEditors(DB, input.UsersID, user)
	panicIf(err)
	return "users_became_editors", nil
}

func (r *mutationResolver) KickUsersFromRoom(ctx context.Context, input model.UsersToRoom) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	room, err := db.GetModRoom(DB, input.RoomID, user)
	panicIf(err)
	err = room.KickUsers(DB, input.UsersID, user)
	panicIf(err)
	return "users_kicked", nil
}

func (r *mutationResolver) KickGroupsFromRoom(ctx context.Context, input model.GroupsToRoom) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	room, err := db.GetModRoom(DB, input.RoomID, user)
	panicIf(err)
	err = room.KickGroups(DB, input.GroupsID, user)
	panicIf(err)
	return "groups_kicked", nil
}

func (r *mutationResolver) KickEditorsFromRoom(ctx context.Context, input model.UsersToRoom) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	room, err := db.GetModRoom(DB, input.RoomID, user)
	panicIf(err)
	err = room.KickEditors(DB, input.UsersID, user)
	panicIf(err)
	return "users_became_editors", nil
}

func (r *mutationResolver) DeleteRoom(ctx context.Context, id string) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	group, err := db.GetModRoom(DB, id, user)
	panicIf(err)
	err = group.DeleteRoom(DB, user)
	panicIf(err)
	return "room_deleted", nil
}

func (r *queryResolver) GetAuthUser(ctx context.Context) (*model.User, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	err := user.Get(DB)
	panicIf(err)
	return &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.Nickname,
		Firstname: user.Firstname, Lastname: user.Lastname,
		Email: user.Email, Role: strconv.Itoa(int(user.RoleID))}, nil
}

func (r *queryResolver) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	if us := auth.ForContext(ctx); us == nil {
		panic("1")
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

func (r *queryResolver) GetUsers(ctx context.Context, nicknames []string) ([]*model.User, error) {
	if us := auth.ForContext(ctx); us == nil {
		panic("1")
	}
	var (
		users        db.Users
		gettingUsers []*model.User
	)
	usersMap := make(map[string]bool)
	for _, user := range nicknames {
		usersMap[user] = false
	}
	for _, user := range users.GetUsers(DB, nicknames) {
		usersMap[user.Nickname] = false
		gettingUsers = append(gettingUsers, &model.User{ID: strconv.Itoa(int(user.ID)),
			Nickname:  user.Nickname,
			Firstname: user.Firstname, Lastname: user.Lastname,
			Email: user.Email, Role: strconv.Itoa(int(user.RoleID))})
	}
	for _, user := range nicknames {
		if exist := usersMap[user]; !exist {
			gettingUsers = append(gettingUsers, &model.User{ID: "0", Nickname: user,
				Firstname: "", Lastname: "",
				Email: "", Role: "0", Error: "user_not_found"})
		}
	}
	return gettingUsers, nil
}

func (r *queryResolver) GetGroups(ctx context.Context, input model.InfoGroups) ([]*model.Group, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	if input.IsPrivate {
		return getGroups(db.GetPrivateGroups(DB, input.GroupsID, user)), nil
	}
	return getGroups(db.GetPublicGroups(DB, input.GroupsID)), nil
}

func (r *queryResolver) GetMyGroups(ctx context.Context) ([]*model.Group, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	return getGroups(db.GetMyGroups(DB, user)), nil
}

func (r *queryResolver) GetMyRooms(ctx context.Context) ([]*model.Room, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	return getRooms(db.GetMyRooms(DB, user)), nil
}

func (r *queryResolver) GetRooms(ctx context.Context, id []string) ([]*model.Room, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		panic("1")
	}
	return getRooms(db.GetRooms(DB, id, user)), nil
}

func (r *queryResolver) GetRoomsMonth(ctx context.Context, month string) ([]*model.Room, error) {
	rooms := []*db.Rooms{}
	if user := auth.ForContext(ctx); user == nil {
		panic("1")
	}
	events, err := db.GetEventsMonth(DB, month)
	panicIf(err)
	for _, room := range events {
		rooms = append(rooms, &room.Room)
	}
	return getRooms(rooms), nil
}

func (r *subscriptionResolver) Notification(ctx context.Context) (<-chan *model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
var DB *gorm.DB = db.Connector()

func getGroups(groups []*db.Groups) []*model.Group {
	var gettingGroups []*model.Group
	for _, group := range groups {
		pack := packagingOfUsers(group)
		gettingGroups = append(gettingGroups, &model.Group{ID: strconv.Itoa(int(group.ID)),
			Name: group.Name, CountUsers: int(group.CountUsers),
			Owner: &model.User{ID: strconv.Itoa(int(pack.owner.ID)),
				Nickname: pack.owner.Nickname, Email: pack.owner.Email},
			Users: pack.modelUsers, Editors: pack.modelUsers})
	}
	return gettingGroups
}
func getRooms(rooms []*db.Rooms) []*model.Room {
	var gettingRooms []*model.Room
	for _, room := range rooms {
		pack := packagingOfUsers(room)
		gettingRooms = append(gettingRooms, &model.Room{
			ID: strconv.Itoa(int(room.ID)), Name: room.Name, Link: room.Link,
			Owner: &model.User{ID: strconv.Itoa(int(pack.owner.ID)),
				Nickname: pack.owner.Nickname, Email: pack.owner.Email},
			Users: pack.modelUsers, Editors: pack.modelUsers})
	}
	return gettingRooms
}

type pocket struct {
	users        []*db.Users
	editors      []*db.Users
	modelUsers   []*model.User
	modelEditors []*model.User
	owner        *db.Users
}

func packagingOfUsers(src interface{}) pocket {
	var pack pocket
	switch v := src.(type) {
	case db.Rooms:
		room := v
		pack.owner = room.GetOwner(DB)
		pack.editors = room.GetEditors(DB)
		pack.users = room.GetUsers(DB)
	case db.Groups:
		group := v
		pack.owner = group.GetOwner(DB)
		pack.editors = group.GetEditors(DB)
		pack.users = group.GetUsers(DB)
	}
	for _, user := range pack.users {
		pack.modelUsers = append(pack.modelUsers, &model.User{ID: strconv.Itoa(int(user.ID)),
			Nickname:  user.Nickname,
			Firstname: user.Firstname, Lastname: user.Lastname,
			Email: user.Email, Role: strconv.Itoa(int(user.RoleID))})
	}
	pack.modelEditors = append(pack.modelEditors, &model.User{ID: strconv.Itoa(int(pack.owner.ID)),
		Nickname: pack.owner.Nickname, Email: pack.owner.Email})
	for _, editor := range pack.editors {
		pack.modelEditors = append(pack.modelEditors, &model.User{ID: strconv.Itoa(int(editor.ID)),
			Nickname:  editor.Nickname,
			Firstname: editor.Firstname, Lastname: editor.Lastname,
			Email: editor.Email, Role: strconv.Itoa(int(editor.RoleID))})
	}
	return pack
}
func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
