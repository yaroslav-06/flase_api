package handlr

import (
	"flase_api/internal/db"
	"flase_api/internal/session"
	"encoding/json"
	"fmt"
	"log"
)

type getPacketHandler struct {
	APIRequestHandler
}

type getPacketInfo struct {
	Pass string `json:"pass"`
}

func NewGetPacketHandlr() APIRequestHandler {
	return &getPacketHandler{}
}

func (lh getPacketHandler) GetName() string {
	return "get packet"
}

func (lh getPacketHandler) HandleRequest(ctx *db.DbCtx, sess *session.Session, data *json.RawMessage) {
	usr := sess.GetUser()
	if usr == nil {
		sess.GetRW().WriteError(lh.GetName(), "not logged in")
		return
	}
	ctx = ctx.GetChild(usr.GetUid())
	info := getPacketInfo{}
	err := json.Unmarshal(*data, &info)
	if err != nil {
		sess.GetRW().WriteError(lh.GetName(), "wrong reqest format")
		return
	}
	pac, err := usr.LoadPacket(ctx, info.Pass, sess.GetSessController().GetScheduler())
	if err != nil {
		sess.GetRW().WriteError(lh.GetName(), fmt.Sprintf("couldn't load the packet: %s", err.Error()))
		return
	}
	sess.GetUser().GetPacket().Unsubscribe(*sess.GetId())
	pac.Subscribe(*sess.GetId(), sess.GetRW())
	sess.GetRW().Write(lh.GetName(), "200")
	log.Printf("user loaded packet: %s\n", pac.GetName())
}
