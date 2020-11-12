package db

const (
	createTableRoles = `
	CREATE TABLE Roles (
		id integer PRIMARY KEY,
		Name text
	);
	`
	createTableGroups = `
	CREATE TABLE Groups (
		id integer PRIMARY KEY,
		Name text
	);
	`
	createTableRoom = `
	CREATE TABLE Rooms (
		id integer PRIMARY KEY,
		Name text,
		Link text,
		Pass text,
		Timer TIME
	);
	`
	createTableAccessLevel = `
	CREATE TABLE AccessLevel (
		id integer PRIMARY KEY,
		Name text
	);
	`
	creataTableUsers = `
	CREATE TABLE Users (
		id integer PRIMARY KEY,
		NickName text UNIQUE,
		Email text UNIQUE,
		Pass text,
		Name text,
		SurName text,
		SecondName text,
		UserRole integer REFERENCES Roles(id),
		UserGroup integer REFERENCES Groups(id)
	);
	`
	createTableRoomAccess = `
	CREATE TABLE RoomAccess (
		id integer PRIMARY KEY,
		User_ID integer REFERENCES Users(id),
		Room_ID integer REFERENCES Rooms(id),
		Level_ID integer REFERENCES AccessLevel(id)
	);
	`
	creataTableStatsRoom = `
	CREATE TABLE StatsRoom (
		id integer PRIMARY KEY,
		Room_ID integer REFERENCES Rooms(id),
		DateTime TIME,
		CountUsers integer
	);
	`
	createTableStatsUser = `
	CREATE TABLE StatsUser (
		User_ID integer REFERENCES Users(id),
		DateTime TIME,
		Room_ID integer REFERENCES Rooms(id),
		TimeSpend TIME
	);
	`
)
