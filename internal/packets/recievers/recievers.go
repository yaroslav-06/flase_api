package recievers

import (
	"flase_api/internal/db"
	"encoding/json"
	"fmt"
)

type RecieverSystem struct {
	handlr 		map[string] RecieverHandler
	data 		*json.RawMessage
}

type RecieverHandler interface {
	LoadReciever(data *json.RawMessage) Reciever
	DbLoad(ctx *db.DbCtx) (Reciever, error)
}

type Reciever interface {
	Send() error
	DbSave(ctx *db.DbCtx) error
	Destruct(ctx *db.DbCtx) error
}

func NewRecieverSystem() *RecieverSystem {
	rs := &RecieverSystem{
		handlr: make(map[string]RecieverHandler),
	}
	return rs
}

type commonGetter struct {
	RecieverType string `json:"type"`
	Message    string `json:"message"`
}

func (sys *RecieverSystem) LoadNewReciever(ctx *db.DbCtx, data *json.RawMessage) error {
	var listOfActionsRaw []*json.RawMessage
	json.Unmarshal(*data, &listOfActionsRaw)
	for _, v := range listOfActionsRaw {
		tg := commonGetter{}
		err := json.Unmarshal(*v, &tg)
		if err != nil {
			return err
		}
		hdlr, exists := sys.handlr[tg.RecieverType]
		if !exists {
			return fmt.Errorf("can't get action handler '%s'", tg.RecieverType)
		}
		hdlr.LoadReciever(v)
	}
	sys.DbSave(ctx)
	return nil
}

func (sys *RecieverSystem) DbSave(ctx *db.DbCtx) error {
	ctx.SaveString("recievers", string(*sys.data))
	return nil
}

func (sys *RecieverSystem) DbLoad(ctx *db.DbCtx) error {
	data, err := ctx.GetString("recievers")
	if err != nil{
		return err
	}
	data_byte := json.RawMessage(data)
	sys.data = &data_byte
	return nil
}

func (sys *RecieverSystem) Destruct(ctx *db.DbCtx) error {
	ctx.ClearField("recievers")
	return nil
}
