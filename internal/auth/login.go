package auth

import (
	"flase_api/internal/db"
	"fmt"
)

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (auth *AuthSystem) Login(ctx *db.DbCtx, info *LoginInfo) (*User, error) {
	uid, err := auth.FromUsername(&info.Username)
	if err != nil {
		return nil, fmt.Errorf("can't find user with this username: %w", err)
	}
	usr, err := auth.LoadUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("unexpected error: can't find user from uid: %w", err)
	}
	if err := usr.checkPassword(info.Password); err != nil {
		return nil, fmt.Errorf("wrong password")
	}
	return usr, nil
}
