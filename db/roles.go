package db

//global roles
var (
	SuperAdmin = Roles{RoleName: "SuperAdmin"}
	Admin      = Roles{RoleName: "Admin"}
	Moder      = Roles{RoleName: "Moder"}
	User       = Roles{RoleName: "User"}
	roles      = []Roles{SuperAdmin, Admin, Moder, User}
)

//roles in rooms(Access Level)
var (
	Owner    = AccessLevel{Name: "Owner"}
	Editor   = AccessLevel{Name: "Editor"}
	Listener = AccessLevel{Name: "Listener"}
	alevels  = []AccessLevel{Owner, Editor, Listener}
)
