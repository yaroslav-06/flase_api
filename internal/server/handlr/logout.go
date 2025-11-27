package handlr

import (
	"flase_api/internal/db"
	"flase_api/internal/session"
	"encoding/json"
)

type logoutHandler struct {
	APIRequestHandler
}

func NewLogoutHandlr() APIRequestHandler {
	return &logoutHandler{}
}

func (lh logoutHandler) GetName() string {
	return "logout"
}

func (lh logoutHandler) HandleRequest(ctx *db.DbCtx, sess *session.Session, data *json.RawMessage) {
	usr := sess.GetUser()
	if usr == nil {
		sess.GetRW().WriteError(lh.GetName(), "not logged in")
		return
	}
	sess.SetUser(nil)
	sess.GetRW().Write(lh.GetName(), "200")
}
