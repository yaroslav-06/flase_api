package actions

import (
	"flase_api/internal/db"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

type timeChngHandler struct {
	ActionHandler
}

type timeChngAction struct {
	Action
	Duration int64 `json:"duration"`
}

func (hdlr *timeChngHandler) LoadDeleteAction(data *json.RawMessage) Action {
	vl := timeChngAction{}
	json.Unmarshal(*data, &vl)
	return &vl
}

func (hdlr *timeChngHandler) DbLoad(ctx *db.DbCtx) (Action, error) {
	vl := timeChngAction{}
	fmt.Printf("get duration: '%s'\n", ctx.GetPath())
	duration, err := ctx.GetString("duration")
	if err != nil {
		return nil, fmt.Errorf("couldn't load duration string of timechanger: %w", err)
	}
	vl.Duration, err = strconv.ParseInt(duration, 10, 64)
	if err != nil {
		return nil, err
	}
	return &vl, nil
}

func (hdlr *timeChngHandler) Name() string {
	return "time changer"
}

func (ac *timeChngAction) Perform(ctx *db.DbCtx, packet PacketInterface) {
	log.Println("fine")
	tm := packet.GetDeliveryTime()
	*tm = tm.Add(-time.Minute * time.Duration(ac.Duration))
	packet.SetDeliveryTime(ctx, tm)
}

func (ac *timeChngAction) DbSave(ctx *db.DbCtx) error {
	fmt.Printf("save duration: '%s'\n", ctx.GetPath())
	ctx.SaveString("duration", fmt.Sprint(ac.Duration))
	return nil
}

func (ac *timeChngAction) Destruct(ctx *db.DbCtx) error {
	err := ctx.ClearField("duration")
	return err
}
