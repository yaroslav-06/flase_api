package actions

import (
	"flase_api/internal/db"
	"encoding/json"
	"log"
)

type deleteHandler struct {
	ActionHandler
}

type deleteAction struct {
	Action
}

func (hdlr *deleteHandler) LoadDeleteAction(data *json.RawMessage) Action {
	return &deleteAction{}
}

func (hdlr *deleteHandler) DbLoad(ctx *db.DbCtx) (Action, error) {
	return &deleteAction{}, nil
}

func (hdlr *deleteHandler) Name() string {
	return "delete"
}

func (ac *deleteAction) Perform(ctx *db.DbCtx, packet PacketInterface) {
	log.Println("performing delete")
	packet.Destruct(ctx.GetParent().GetParent())
}

func (ac *deleteAction) DbSave(ctx *db.DbCtx) error {
	return nil
}

func (ac *deleteAction) Destruct(ctx *db.DbCtx) error {
	return nil
}
