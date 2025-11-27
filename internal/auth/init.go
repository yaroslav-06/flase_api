package auth

import (
	"flase_api/internal/db"
	uniqueid "flase_api/internal/unique_id"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type AuthSystem struct {
	uidGen      *uniqueid.Generator
	loadedUsers map[string]*User
	ctx         *db.DbCtx
}

func InitAuthSystem(ctx *db.DbCtx) (*AuthSystem, error) {
	auth := &AuthSystem{}
	auth.ctx = ctx
	auth.uidGen = uniqueid.NewGenerator(auth)
	auth.loadedUsers = make(map[string]*User)

	config_ctx := ctx.GetChild(":config")

	usrn, err := config_ctx.GetString("admin_username")
	if err != nil{
		var pass string

		fmt.Printf("no admin user found.\n Enter admin name: ")
		fmt.Scanf("%s", &usrn)
		fmt.Printf("enter admin password: ")
		fmt.Scanf("%s", &pass)
		fmt.Printf("creating user %s\n", usrn)
		fmt.Printf("creating user with pass %s\n", pass)
		_, err = auth.CreateUser("e@e", usrn, pass)
		if err != nil {
			return nil, fmt.Errorf("couldn't create such user: %s", err.Error())
		}
		fmt.Println("creation successful")
		err = config_ctx.SaveString("admin_username", usrn)
		if err != nil {
			return nil, fmt.Errorf("couldn't save admin user: %s", err.Error())
		}
	}

	uid, err := auth.FromUsername(&usrn)
	if err != nil {
		return nil, fmt.Errorf("couldn't get admin user (created improperly: %s)", err.Error())
	}
	usr, err := auth.LoadUser(ctx, uid)
	log.Printf("uid:'%s', uname:'%s'\n", usr.uid, usr.username)
	log.Println("exists: ", auth.Exists(usr.uid))

	return auth, nil
}

func (auth *AuthSystem) FromUsername(usrnm *string) (string, error) {
	uid, err := auth.ctx.GetMapField("fromUsername", *usrnm)
	if uid == "" {
		return uid, fmt.Errorf("can't find user with this username")
	}
	if err == redis.Nil {
		return uid, nil
	}
	return uid, err
}

func (auth *AuthSystem) SetUsername(usrnm *string, uid *string) {
	log.Printf("uid: '%s', err: '%s'", *uid, *usrnm)
	auth.ctx.SaveMapField("fromUsername", *usrnm, *uid)
	log.Println("saved from username: ")
	ud, err := auth.ctx.GetMapField("fromUsername", *usrnm)
	log.Printf("uid: '%s', err: '%s'", ud, err)
}

func (auth *AuthSystem) SaveUser(usr *User) {
	auth.SetUsername(&usr.username, &usr.uid)
	usr.dbSave(auth.ctx.GetChild(usr.uid))
}

func (auth *AuthSystem) Exists(uid string) bool {
	_, err := auth.LoadUser(auth.ctx, uid)
	return err == nil
}

func (auth *AuthSystem) LoadUser(ctx *db.DbCtx, uid string) (*User, error) {
	if usr, exists := auth.loadedUsers[uid]; exists {
		log.Println("user already loaded")
		return usr, nil
	}
	rs := &User{}
	err := rs.DbLoad(ctx.GetChild(uid))
	if err != nil {
		return nil, err
	}
	auth.loadedUsers[uid] = rs
	return rs, nil
}
