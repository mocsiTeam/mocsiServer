package db

const (
	createTableRoles = `
	CREATE TABLE Roles (
		id integer PRIMARY KEY,
		Name text UNIQUE NOT NULL
	);
	`
	createTableGroups = `
	CREATE TABLE Groups (
		id integer PRIMARY KEY,
		Name text UNIQUE NOT NULL
	);
	`
	createTableRoom = `
	CREATE TABLE Rooms (
		id integer PRIMARY KEY,
		Name text UNIQUE NOT NULL,
		Link text UNIQUE NOT NULL,
		Pass text,
		Timer TIME
	);
	`
	createTableAccessLevel = `
	CREATE TABLE AccessLevel (
		id integer PRIMARY KEY,
		Name text UNIQUE NOT NULL
	);
	`
	creataTableUsers = `
	CREATE TABLE Users (
		id integer PRIMARY KEY,
		NickName text UNIQUE NOT NULL,
		Email text UNIQUE NOT NULL,
		Pass text NOT NULL,
		Name text NOT NULL,
		SurName text NOT NULL,
		SecondName text NOT NULL,
		UserRole integer REFERENCES Roles(id) NOT NULL,
		UserGroup integer REFERENCES Groups(id) NOT NULL
	);
	`
	createTableRoomAccess = `
	CREATE TABLE RoomAccess (
		id integer PRIMARY KEY,
		User_ID integer REFERENCES Users(id) NOT NULL,
		Room_ID integer REFERENCES Rooms(id) NOT NULL,
		Level_ID integer REFERENCES AccessLevel(id) NOT NULL
	);
	`
	creataTableStatsRoom = `
	CREATE TABLE StatsRoom (
		id integer PRIMARY KEY,
		Room_ID integer REFERENCES Rooms(id) NOT NULL,
		DateTime TIME,
		CountUsers integer
	);
	`
	createTableStatsUser = `
	CREATE TABLE StatsUser (
		User_ID integer REFERENCES Users(id) NOT NULL,
		DateTime TIME,
		Room_ID integer REFERENCES Rooms(id) NOT NULL,
		TimeSpend TIME
	);
	`
)
