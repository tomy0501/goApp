package db

import "database/sql"

func NewUserDB(db *sql.DB) UserDB {
	return UserDB{
		UserDB: db,
	}
}

type UserDB struct {
	UserDB *sql.DB
}
