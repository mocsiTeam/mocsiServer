// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	Firstname string `json:"firstname"`
	LastName  string `json:"lastName"`
}

type RefreshTokenInput struct {
	Token string `json:"token"`
}

type User struct {
	ID        string   `json:"id"`
	Nickname  string   `json:"nickname"`
	Firstname string   `json:"firstname"`
	LastName  string   `json:"lastName"`
	Email     string   `json:"email"`
	Role      string   `json:"role"`
	Groups    []string `json:"groups"`
}
