package handlr

import (
	"flase_api/internal/db"
	"flase_api/internal/encoder"
	"flase_api/internal/packets"
	"flase_api/internal/session"
	"encoding/json"
	"fmt"
	"log"
)

type packetCreatorHandler struct {
	APIRequestHandler
}

func NewPacketCreatorHandlr() APIRequestHandler {
	return &packetCreatorHandler{}
}

func (lh packetCreatorHandler) GetName() string {
	return "packet creator"
}

func (lh packetCreatorHandler) HandleRequest(ctx *db.DbCtx, sess *session.Session, data *json.RawMessage) {
	usr := sess.GetUser()
	if usr == nil {
		sess.GetRW().WriteError(lh.GetName(), "not logged in")
		return
	}
	ctx = ctx.GetChild(usr.GetUid())
	pac := packets.NewPacket()
	log.Println(sess.GetSessController())
	log.Println(sess.GetSessController().GetScheduler())
	log.Println(*sess.GetUid())
	pass, err := pac.FromJsonPacket(ctx, data, *sess.GetUid(), sess.GetSessController().GetScheduler())
	if err != nil {
		sess.GetRW().WriteError(lh.GetName(), fmt.Sprintf("can't create packet: %s", err.Error()))
		return
	}
	hPass := encoder.Enc(pass)
	log.Println("saving the packet: ", ctx.GetChild(hPass).GetPath())
	pac.DbSave(ctx.GetChild(hPass))
	sess.GetRW().Write(lh.GetName(), "200")
	// sess.GetRW().Write(lh.GetName(), fmt.Sprint(pac))
}
