package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mocsiTeam/mocsiServer/api/graph/generated"
	"github.com/mocsiTeam/mocsiServer/api/graph/model"
	"github.com/mocsiTeam/mocsiServer/auth"
	"github.com/mocsiTeam/mocsiServer/auth/jwt"
	"github.com/mocsiTeam/mocsiServer/db"
	"gorm.io/gorm"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var newUser = db.Users{
		Email:     input.Email,
		NickName:  input.Nickname,
		FirstName: input.Firstname,
		LastName:  input.Lastname,
		Pass:      input.Password,
		RoleID:    3,
	}
	if err := newUser.Create(DB); err != nil {
		return "", err
	}
	token, err := jwt.GenerateToken(newUser.NickName)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user = db.Users{
		NickName: input.Nickname,
		Pass:     input.Password,
	}
	if correct := user.Authenticate(DB); !correct {
		return "", &db.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.NickName)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
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
		allUsers = append(allUsers, &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.NickName,
			Firstname: user.FirstName, LastName: user.LastName,
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
		for _, v := range input {
			if v == user.NickName {
				gettingUsers = append(gettingUsers, &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.NickName,
					Firstname: user.FirstName, LastName: user.LastName,
					Email: user.Email, Role: strconv.Itoa(int(user.RoleID))})
			} else {
				gettingUsers = append(gettingUsers, &model.User{ID: "0", Nickname: "",
					Firstname: "", LastName: "",
					Email: "", Role: "0", Error: "user_not_found"})
			}
		}

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
