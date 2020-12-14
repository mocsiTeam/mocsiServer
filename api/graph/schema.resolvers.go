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

func (r *queryResolver) GetAuthUser(ctx context.Context) (*model.User, error) {
	var user *db.Users
	if user = auth.ForContext(ctx); user == nil {
		return &model.User{}, fmt.Errorf("access denied")
	}
	if err := user.Get(DB); err != nil {
		return &model.User{}, err
	}
	return &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.NickName,
		Firstname: user.FirstName, LastName: user.LastName,
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
	for _, user := range users.GetUsers(DB, input) {
		gettingUsers = append(gettingUsers, &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.NickName,
			Firstname: user.FirstName, LastName: user.LastName,
			Email: user.Email, Role: strconv.Itoa(int(user.RoleID))})
	}
	return gettingUsers, nil
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
