package handlr

import (
	"flase_api/internal/auth"
	"flase_api/internal/db"
	"flase_api/internal/session"
	"encoding/json"
	"fmt"
	"log"
)

type loginHandler struct {
	APIRequestHandler
}

func NewLoginHandlr() APIRequestHandler {
	return &loginHandler{}
}

func (lh loginHandler) GetName() string {
	return "login"
}

func (lh loginHandler) HandleRequest(ctx *db.DbCtx, sess *session.Session, data *json.RawMessage) {
	logInfo := auth.LoginInfo{}
	err := json.Unmarshal(*data, &logInfo)
	if err != nil {
		msg := fmt.Sprintf("login info unmarshal: %s\n", err.Error())
		sess.GetRW().WriteError("auth", msg)
		log.Println(msg)
		return
	}
	authSystem := sess.GetAuthSystem()
	usr, err := authSystem.Login(ctx, &logInfo)
	if err != nil {
		msg := fmt.Sprintf("can't get user from credentials: %s", err.Error())
		sess.GetRW().WriteError("auth", msg)
		log.Println(msg)
		return
	}
	sess.SetUser(usr)
	sess.GetRW().Write("auth", "200")
}
