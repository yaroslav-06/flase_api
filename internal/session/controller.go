package session

import (
	"flase_api/internal/db"
	smartscheduler "flase_api/internal/smart_scheduler"
	uniqueid "flase_api/internal/unique_id"
	"container/list"
)

type SessionController struct {
	fromId  map[string]Session
	fromUid map[string]list.List

	msgScheduler	*smartscheduler.SmartScheduler

	sessIdGen 		*uniqueid.Generator
}

func (controller *SessionController) Exists(sessId string) bool {
	_, exists := controller.fromId[sessId]
	return exists
}

func InitController(ctx *db.DbCtx) *SessionController {
	controller := &SessionController{}
	controller.fromId = make(map[string]Session)
	controller.fromUid = make(map[string]list.List)

	msgScheduler, err := smartscheduler.InitSmartScheduler(ctx, &MsgScheduleEvent{})
	if err != nil{
		panic(err)
	}
	controller.msgScheduler = msgScheduler

	controller.sessIdGen = uniqueid.NewGenerator(controller)

	return controller
}

func (controller *SessionController) FromSessionId(sessId string) *Session {
	sess, exists := controller.fromId[sessId]
	if !exists {
		return nil
	}
	return &sess
}

func (controller *SessionController) GetNewId() string {
	return controller.sessIdGen.GetNewId()
}

func (controller *SessionController) GetScheduler() (*smartscheduler.SmartScheduler) {
	return controller.msgScheduler
}
