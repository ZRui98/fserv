package models

import (
	"context"
	"time"
)

type User struct {
	Username string
	Password string
	LastLoginTime time.Time
}

type UserRepository interface {
	GetUserById(context.Context, string) (*User, error)
	AddUserById(context.Context, *User) (error)
	UpdateUserLastLoginTime(context.Context, string, time.Time) (error)
}

func (pool *DbPool) GetUserById(ctx context.Context, username string) (user *User, err error) {
	var foundUser User
	queryString := `SELECT "username", "password", "last_login" FROM "users" WHERE "username"=$1;`
	err = pool.db.QueryRow(ctx, queryString, username).Scan(&foundUser.Username, &foundUser.Password, &foundUser.LastLoginTime)
	return &foundUser, err
}

func (pool *DbPool) AddUserById(ctx context.Context, user *User) (err error) {
	queryString := `INSERT INTO "users" (username, password, last_login) VALUES ($1, $2, $3);`
	_, err = pool.db.Exec(ctx, queryString, user.Username, user.Password, user.LastLoginTime)
	return err
}

func (pool *DbPool) UpdateUserLastLoginTime(ctx context.Context, username string, newTime time.Time) (err error) {
	queryString := `UPDATE "users" SET last_login=$1 WHERE username=$2`
	_, err = pool.db.Exec(context.Background(), queryString, newTime, username)
	return err
}
