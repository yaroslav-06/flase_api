package actions

import (
	"flase_api/internal/db"
	"flase_api/internal/encoder"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type ActionSystem struct {
	handlr  map[string]ActionHandler
	packet  PacketInterface
	actions map[string]string
}

type PacketInterface interface {
	DbSave(ctx *db.DbCtx) error
	GetDeliveryTime() *time.Time
	SetDeliveryTime(ctx *db.DbCtx, deliveryTime *time.Time) error
	Destruct(ctx *db.DbCtx) error
}

type Action interface {
	Perform(ctx *db.DbCtx, pc PacketInterface)
	DbSave(ctx *db.DbCtx) error
	Destruct(ctx *db.DbCtx) error
}

type ActionHandler interface {
	LoadDeleteAction(data *json.RawMessage) Action
	DbLoad(ctx *db.DbCtx) (Action, error)
	Name() string
}

type commonGetter struct {
	ActionType string `json:"type"`
	Pass       string `json:"pass"`
}

func NewActionSystem(pc PacketInterface) *ActionSystem {
	rs := &ActionSystem{}
	rs.packet = pc

	rs.handlr = map[string]ActionHandler{}

	var hdlr ActionHandler
	hdlr = &deleteHandler{}
	rs.handlr[hdlr.Name()] = hdlr
	hdlr = &timeChngHandler{}
	rs.handlr[hdlr.Name()] = hdlr

	rs.actions = make(map[string]string)
	return rs
}

func (ac *ActionSystem) LoadNewAction(ctx *db.DbCtx, data *json.RawMessage) error {
	var listOfActionsRaw []*json.RawMessage
	json.Unmarshal(*data, &listOfActionsRaw)
	for _, v := range listOfActionsRaw {
		tg := commonGetter{}
		err := json.Unmarshal(*v, &tg)
		if err != nil {
			return err
		}
		hPass := encoder.Enc(tg.Pass)
		if _, exists := ac.actions[hPass]; exists {
			return errors.New("there are multiple actions with the same password")
		}
		hdlr, exists := ac.handlr[tg.ActionType]
		if !exists {
			return fmt.Errorf("can't get action handler '%s'", tg.ActionType)
		}
		act := hdlr.LoadDeleteAction(v)

		err = act.DbSave(ctx.GetChild(hPass))
		if err != nil {
			return err
		}
		ac.actions[hPass] = hdlr.Name()
	}
	ac.DbSave(ctx)
	return nil
}

func (ac *ActionSystem) DbLoad(ctx *db.DbCtx) error {
	var err error
	log.Println("actions load", ctx.GetPath())
	ac.actions, err = ctx.GetMapAll("actions")
	if err != nil {
		return err
	}
	return nil
}

func (ac *ActionSystem) DbSave(ctx *db.DbCtx) error {
	log.Println("actions save: ", ctx.GetPath())
	for key, val := range ac.actions {
		err := ctx.SaveMapField("actions", key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ac *ActionSystem) Destruct(ctx *db.DbCtx) error {
	fmt.Println("action system destruct: ")
	for hPass := range ac.actions {
		act, err := ac.GetAction(ctx, hPass)
		if err != nil {
			return err
		}
		err = act.Destruct(ctx.GetChild(hPass))
		if err != nil{
			return err
		}
	}
	ctx.ClearField("actions")
	return nil
}

func (ac *ActionSystem) GetAction(ctx *db.DbCtx, hPass string) (Action, error) {
	fmt.Println(ac.actions)
	actionType, exists := ac.actions[hPass]
	if !exists {
		return nil, errors.New("couldn't find the action from password")
	}
	hdlr, exists := ac.handlr[actionType]
	if !exists {
		return nil, errors.New("unexpected error couldn't find action of the type")
	}
	act, err := hdlr.DbLoad(ctx.GetChild(hPass))
	return act, err
}
