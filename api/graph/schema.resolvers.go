package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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
	if err := newUser.Create(DB); err != nil {
		return &model.Tokens{}, err
	}
	accessToken, err := jwt.GenerateAccessToken(newUser.Nickname, strconv.Itoa(int(newUser.ID)))
	if err != nil {
		return &model.Tokens{}, err
	}
	refreshToken, err := jwt.GenerateRefreshToken(newUser.Nickname, newUser.Email, strconv.Itoa(int(newUser.ID)))
	if err != nil {
		return &model.Tokens{}, err
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
		return &model.Tokens{}, &db.WrongUsernameOrPasswordError{}
	}
	userID := strconv.Itoa(int(user.ID))
	accessToken, err := jwt.GenerateAccessToken(user.Nickname, userID)
	if err != nil {
		return &model.Tokens{}, err
	}
	refreshToken, err := jwt.GenerateRefreshToken(user.Nickname, user.Email, userID)
	if err != nil {
		return &model.Tokens{}, err
	}
	user.RefreshToken = refreshToken
	DB.Save(&user)
	return &model.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	var user db.Users
	userID, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	refreshToken, err := user.GetRefreshToken(DB, userID)
	if err != nil {
		return "", nil
	}
	if refreshToken != input.Token {
		return "", errors.New("invalid refresh token")
	}
	token, err := jwt.GenerateAccessToken(user.Nickname, userID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) CreateGroup(ctx context.Context, input model.NewGroup) (*model.Group, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		return &model.Group{}, fmt.Errorf("access denied")
	}
	group := &db.Groups{Name: input.Name, Private: input.Private}
	if err := group.Create(DB, user); err != nil {
		return &model.Group{}, err
	}
	owner := &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.Nickname,
		Firstname: user.Firstname, Lastname: user.Lastname,
		Email: user.Email, Role: strconv.Itoa(int(user.RoleID))}
	return &model.Group{ID: strconv.Itoa(int(group.ID)), Name: group.Name, CountUsers: int(group.CountUsers),
		Owner: owner, Users: []*model.User{owner}}, nil
}

func (r *mutationResolver) AddUsersToGroup(ctx context.Context, input model.UsersToGroup) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		return "", fmt.Errorf("access denied")
	}
	group, err := db.GetModGroup(DB, input.GroupID, user)
	if err != nil {
		return "", err
	} else if err := group.AddUsers(DB, input.UsersID, user); err != nil {
		return "", err
	}
	return "users_added", nil
}

func (r *mutationResolver) AddEditorsToGroup(ctx context.Context, input model.UsersToGroup) (string, error) {
	var user *db.Users
	var group *db.Groups
	if user = auth.ForContext(ctx); user == nil {
		return "", fmt.Errorf("access denied")
	}
	group, err := db.GetModGroup(DB, input.GroupID, user)
	if err != nil {
		return "", err
	} else if err := group.AddEditors(DB, input.UsersID, user); err != nil {
		return "", err
	}
	return "users_became_editors", nil
}

func (r *mutationResolver) KickUsersFromGroup(ctx context.Context, input model.UsersToGroup) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		return "", fmt.Errorf("access denied")
	}
	group, err := db.GetModGroup(DB, input.GroupID, user)
	if err != nil {
		return "", err
	} else if err := group.KickUsers(DB, input.UsersID, user); err != nil {
		return "", err
	}
	return "users_kicked", nil
}

func (r *mutationResolver) DeleteGroup(ctx context.Context, input string) (string, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		return "", fmt.Errorf("access denied")
	}
	group, err := db.GetModGroup(DB, input, user)
	if err != nil {
		return "", err
	} else if err := group.DeleteGroup(DB, user); err != nil {
		return "", err
	}
	return "group_deleted", nil
}

func (r *mutationResolver) CreateRoom(ctx context.Context, input model.NewRoom) (*model.Room, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddUsersToRoom(ctx context.Context, input *model.UsersToRoom) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddGroupToRoom(ctx context.Context, input *model.GroupsToRoom) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) KickUsersFromRoom(ctx context.Context, input *model.UsersToRoom) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) KickGroupsFromRoom(ctx context.Context, input *model.GroupsToRoom) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteRoom(ctx context.Context, input string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetAuthUser(ctx context.Context) (*model.User, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		return &model.User{}, fmt.Errorf("access denied")
	}
	if err := user.Get(DB); err != nil {
		return &model.User{}, err
	}
	return &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.Nickname,
		Firstname: user.Firstname, Lastname: user.Lastname,
		Email: user.Email, Role: strconv.Itoa(int(user.RoleID))}, nil
}

func (r *queryResolver) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	if us := auth.ForContext(ctx); us == nil {
		return []*model.User{}, fmt.Errorf("access denied")
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

func (r *queryResolver) GetUsers(ctx context.Context, input []string) ([]*model.User, error) {
	if us := auth.ForContext(ctx); us == nil {
		return []*model.User{}, fmt.Errorf("access denied")
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

func (r *queryResolver) GetGroups(ctx context.Context, input model.InfoGroups) ([]*model.Group, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		return []*model.Group{}, fmt.Errorf("access denied")
	}
	if input.IsPrivate {
		return getGroups(db.GetPrivateGroups(DB, input.GroupsID, user), input.GroupsID, input.IsPrivate), nil
	}
	return getGroups(db.GetPublicGroups(DB, input.GroupsID), input.GroupsID, input.IsPrivate), nil
}

func (r *queryResolver) GetMyGroups(ctx context.Context) ([]*model.Group, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		return []*model.Group{}, fmt.Errorf("access denied")
	}
	fmt.Println(user)
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
					Owner: &model.User{}, Users: []*model.User{}, Error: "group_not_found"})
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
