package handlr

import (
	"flase_api/internal/db"
	"flase_api/internal/session"
	"encoding/json"
	"fmt"
)

type performActionHandler struct {
	APIRequestHandler
}

type performActionInfo struct {
	Pass string `json:"pass"`
}

func NewPerformActionHandler() APIRequestHandler {
	return &performActionHandler{}
}

func (lh performActionHandler) GetName() string {
	return "perform action"
}

func (lh performActionHandler) HandleRequest(ctx *db.DbCtx, sess *session.Session, data *json.RawMessage) {
	usr := sess.GetUser()
	if usr == nil {
		sess.GetRW().WriteError(lh.GetName(), "not logged in")
		return
	}
	ctx = ctx.GetChild(usr.GetUid())
	pac := usr.GetPacket()
	if pac == nil {
		sess.GetRW().WriteError(lh.GetName(), "no packet openned")
		return
	}
	ctx = pac.UpdateDbCtx(ctx)
	info := performActionInfo{}
	err := json.Unmarshal(*data, &info)
	if err != nil {
		sess.GetRW().WriteError(lh.GetName(), "wrong reqest format")
		return
	}
	err = pac.PerformAction(ctx, info.Pass)
	if err != nil {
		sess.GetRW().WriteError(lh.GetName(), fmt.Sprintf("couldn't perform the action: %s", err.Error()))
		return
	}
	sess.GetRW().Write(lh.GetName(), "200")
}
