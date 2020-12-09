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

var DB *gorm.DB = db.Connector()

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

func (r *queryResolver) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	if us := auth.ForContext(ctx); us == nil {

		return []*model.User{}, fmt.Errorf("access denied")
	}
	var users db.Users
	var allUsers []*model.User
	for _, user := range users.GetAll(DB) {
		allUsers = append(allUsers, &model.User{ID: strconv.Itoa(int(user.ID)), Nickname: user.NickName,
			Firstname: user.FirstName, LastName: user.LastName,
			Email: user.Email, Role: strconv.Itoa(int(user.RoleID))})
	}
	return allUsers, nil
}

func (r *queryResolver) GetUsers(ctx context.Context, input []string) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
