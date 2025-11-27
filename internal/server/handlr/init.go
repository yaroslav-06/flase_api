package handlr

import (
	"flase_api/internal/db"
	"flase_api/internal/session"
	"encoding/json"
)

type APIRequestHandler interface {
	GetName() string
	HandleRequest(ctx *db.DbCtx, sess *session.Session, data *json.RawMessage)
}

func Get() *map[string]APIRequestHandler {
	rs := make(map[string]APIRequestHandler)

	var hdlr APIRequestHandler
	hdlr = NewLoginHandlr()
	rs[hdlr.GetName()] = hdlr
	hdlr = NewLogoutHandlr()
	rs[hdlr.GetName()] = hdlr

	hdlr = NewPacketCreatorHandlr()
	rs[hdlr.GetName()] = hdlr

	hdlr = NewGetPacketHandlr()
	rs[hdlr.GetName()] = hdlr

	hdlr = NewPerformActionHandler()
	rs[hdlr.GetName()] = hdlr

	return &rs
}
