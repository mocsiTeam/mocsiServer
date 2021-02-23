// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Group struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	CountUsers int     `json:"countUsers"`
	Owner      *User   `json:"owner"`
	Editors    []*User `json:"editors"`
	Users      []*User `json:"users"`
	Error      string  `json:"error"`
}

type GroupsToRoom struct {
	RoomID   string   `json:"roomID"`
	GroupsID []string `json:"groupsID"`
}

type InfoGroups struct {
	GroupsID  []string `json:"groupsID"`
	IsPrivate bool     `json:"isPrivate"`
}

type Login struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type NewGroup struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

type NewRoom struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type NewUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type RefreshTokenInput struct {
	Token string `json:"token"`
}

type Room struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Link    string  `json:"link"`
	Owner   *User   `json:"owner"`
	Editors []*User `json:"editors"`
	Users   []*User `json:"users"`
	Error   string  `json:"error"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type User struct {
	ID        string   `json:"id"`
	Nickname  string   `json:"nickname"`
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Email     string   `json:"email"`
	Role      string   `json:"role"`
	Groups    []string `json:"groups"`
	Error     string   `json:"error"`
}

type UsersToGroup struct {
	GroupID string   `json:"groupID"`
	UsersID []string `json:"usersID"`
}

type UsersToRoom struct {
	RoomID  string   `json:"roomID"`
	UsersID []string `json:"usersID"`
}
