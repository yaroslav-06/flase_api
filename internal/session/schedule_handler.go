package session

import (
	smartscheduler "flase_api/internal/smart_scheduler"
	"flase_api/internal/telegram"
	"encoding/json"
	"fmt"
	"strings"
)

type MsgScheduleEvent struct {
	smartscheduler.SchedulerEventType
}

func (ev *MsgScheduleEvent) execEach(dt *json.RawMessage) error {
	var vl map[string] string
	err := json.Unmarshal(*dt, &vl)
	if err != nil{
		return fmt.Errorf("message_decoding: %s", err)
	}
	tp, exists := vl["type"]
	if !exists {
		return fmt.Errorf("no type specified")
	}
	msg, exists := vl["message"]
	if !exists {
		return fmt.Errorf("no message specified")
	}
	if(tp == "telegram"){
		fmt.Println("execute telegram")
		usrn, exists := vl["username"]
		if !exists {
			return fmt.Errorf("no tg username specified")
		}
		return telegram.SendMsg(usrn, msg)
	}
	return nil
}

func (ev *MsgScheduleEvent) ExecuteTask(val string) error {
	fmt.Println("execute task")
	fmt.Println(val)
	var listOfRaw []*json.RawMessage
	uid := strings.Split(val, ":")[0]
	val = val[len(uid) + 1:]
	fmt.Println(val)
	err := json.Unmarshal([]byte(val), &listOfRaw)
	if err != nil{
		return err
	}
	for _, v := range listOfRaw {
		fmt.Println("h:")
		fmt.Println(v)
		er := ev.execEach(v)
		if er != nil{
			fmt.Println(err)
			err = er
		}
	}
	return err
}
