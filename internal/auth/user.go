package auth

import (
	"flase_api/internal/db"
	"flase_api/internal/packets"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	email          string
	username       string
	uid            string
	hashedPassword string

	packet       *packets.Packet
	packetSystem *PacketSystem
}

func (usr *User) GetUid() string {
	return usr.uid
}

func (auth *AuthSystem) CreateUser(email string, username string, pass string) (*User, error) {
	ps, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := User{
		uid:   auth.uidGen.GetNewId(),
		email: email, username: username,
		packetSystem:   NewPacketSystem(),
		hashedPassword: string(ps)}
	auth.SaveUser(&u)
	return &u, nil
}

func (usr *User) checkPassword(pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(usr.hashedPassword), []byte(pass))
}

func (usr *User) dbSave(ctx *db.DbCtx) {
	ctx.SaveString("email", usr.email)
	ctx.SaveString("username", usr.username)
	ctx.SaveString("uid", usr.uid)
	ctx.SaveString("hashedPassword", usr.hashedPassword)
}

func (usr *User) DbLoad(ctx *db.DbCtx) error {
	var err error
	usr.email, err = ctx.GetString("email")
	if err != nil {
		return err
	}
	usr.uid, err = ctx.GetString("uid")
	if err != nil {
		return err
	}
	usr.username, err = ctx.GetString("username")
	if err != nil {
		return err
	}
	usr.hashedPassword, err = ctx.GetString("hashedPassword")
	if err != nil {
		return err
	}
	usr.packetSystem = NewPacketSystem()
	if usr.email == "" || usr.uid == "" || usr.username == "" {
		return fmt.Errorf("can't load the user info")
	}
	return nil
}
