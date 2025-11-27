package packets

import (
	"flase_api/internal/db"
	"flase_api/internal/encoder"
	messagereadwrite "flase_api/internal/message_read_write"
	"flase_api/internal/packets/actions"
	"flase_api/internal/packets/recievers"
	smartscheduler "flase_api/internal/smart_scheduler"
	"encoding/json"
	"fmt"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"

type Packet struct {
	name         string
	deliveryTime *time.Time
	sched_id     string
	hPass        string

	destructed  bool
	sched 			*smartscheduler.SmartScheduler
	subscribers map[string]*messagereadwrite.ReadWriter
	actionSys   *actions.ActionSystem
	recieverSys *recievers.RecieverSystem
}

type newPacketInfo struct {
	Name         string           `json:"name"`
	DeliveryTime string           `json:"deliveryTime"`
	Pass         string           `json:"pass"`
	Receivers    *json.RawMessage `json:"recievers"`
	Actions      *json.RawMessage `json:"actions"`
}

func (p *Packet) UpdateDbCtx(ctx *db.DbCtx) *db.DbCtx {
	return ctx.GetChild(p.hPass).GetChild("actions")
}

func (p *Packet) SetDeliveryTime(ctx *db.DbCtx, deliveryTime *time.Time) error {
	p.deliveryTime = deliveryTime
	p.sched.RemoveEvent(p.sched_id)
	p.sched.AddEvent(p.sched_id, (*deliveryTime).UTC())
	p.DbSave(ctx)

	return nil
}

func (p *Packet) GetDeliveryTime() *time.Time {
	return p.deliveryTime
}

func NewPacket() *Packet {
	rs := &Packet{}
	rs.destructed = false
	rs.actionSys = actions.NewActionSystem(rs)
	rs.recieverSys = recievers.NewRecieverSystem()
	rs.subscribers = make(map[string]*messagereadwrite.ReadWriter)
	return rs
}

func (rs *Packet) FromJsonPacket(ctx *db.DbCtx, data *json.RawMessage, uid string, scheduler *smartscheduler.SmartScheduler) (string, error) {
	packetInfo := newPacketInfo{}
	err := json.Unmarshal(*data, &packetInfo)
	if err != nil {
		return "", err
	}
	rs.name = packetInfo.Name
	rs.hPass = encoder.Enc(packetInfo.Pass)
	_, err = DbLoad(ctx.GetChild(rs.hPass), scheduler)
	if err == nil {
		return "", fmt.Errorf("the packet with this pass already exists")
	}
	tm, err := time.Parse(TimeLayout, packetInfo.DeliveryTime)
	if err != nil {
		tm, err = time.Parse(time.RFC3339, packetInfo.DeliveryTime)
		if err != nil{
			return "", err
		}
	}
	fmt.Printf("\n\n\n")
	rs.deliveryTime = &tm
	err = rs.actionSys.LoadNewAction(ctx.GetChild(rs.hPass).GetChild("actions"), packetInfo.Actions)
	if err != nil {
		return "", err
	}

	fmt.Println(packetInfo.Receivers)
	fmt.Printf("\n\n\n")
	rs.sched_id = uid + ":" + string(*packetInfo.Receivers)
	rs.sched = scheduler
	fmt.Printf("sched_id: %s\n", rs.sched_id)
	fmt.Printf("sched: %s\n", scheduler)
	fmt.Printf("time: %s\n", *rs.deliveryTime)
	
	scheduler.AddEvent(rs.sched_id, (*rs.deliveryTime).UTC())
	// err = rs.recieverSys.Load(ctx.GetChild(rs.hPass).GetChild("recievers"), packetInfo.Receivers)
	// if err != nil {
	// 	return "", err
	// }
	return packetInfo.Pass, nil
}

func (pac *Packet) IsDestructed() bool {
	return pac.destructed
}

func DbLoad(ctx *db.DbCtx, sched *smartscheduler.SmartScheduler) (*Packet, error) {
	rs := NewPacket()
	rs.sched = sched
	fmt.Printf("loading packet at: %s\n", ctx.GetPath())
	timestr, err := ctx.GetString("deliveryTime")
	if err != nil {
		return nil, err
	}
	tm, err := time.Parse(TimeLayout, timestr)
	if err != nil {
		return nil, err
	}
	rs.name, err = ctx.GetString("name")
	if err != nil {
		return nil, err
	}
	rs.sched_id, err = ctx.GetString("sched_id")
	if err != nil {
		return nil, err
	}
	rs.hPass, err = ctx.GetString("pass")
	if err != nil {
		return nil, err
	}
	rs.deliveryTime = &tm
	rs.actionSys.DbLoad(ctx.GetChild("actions"))
	return rs, nil
}

func (pc *Packet) DbSave(ctx *db.DbCtx) error {
	fmt.Printf("saving action to[%s]", ctx.GetPath())
	err := ctx.SaveString("deliveryTime", pc.deliveryTime.Format(TimeLayout))
	if err != nil {
		return err
	}
	err = ctx.SaveString("name", pc.name)
	if err != nil {
		return err
	}
	err = ctx.SaveString("sched_id", pc.sched_id)
	if err != nil {
		return err
	}
	err = ctx.SaveString("pass", pc.hPass)
	if err != nil {
		return err
	}
	return nil
}

func (pc *Packet) Destruct(ctx *db.DbCtx) error {
	fmt.Printf("destructing packet at: %s\n", ctx.GetPath())
	err := ctx.ClearField("deliveryTime")
	if err != nil {
		return err
	}
	err = ctx.ClearField("name")
	if err != nil {
		return err
	}
	err = ctx.ClearField("receiverIds")
	if err != nil {
		return err
	}
	err = pc.actionSys.Destruct(ctx.GetChild("actions"))
	if err != nil{
		return err
	}
	err = pc.recieverSys.Destruct(ctx.GetChild("recievers"))
	if err != nil{
		return err
	}
	pc.sched.RemoveEvent(pc.sched_id)
	pc.destructed = true
	*pc.deliveryTime = time.Now().UTC()
	return nil
}

func (pac *Packet) PerformAction(ctx *db.DbCtx, pass string) error {
	fmt.Printf("loading action from pass: %s, [%s]\n", pass, ctx.GetPath())
	act, err := pac.actionSys.GetAction(ctx, encoder.Enc(pass))
	if err != nil {
		return err
	}
	act.Perform(ctx.GetChild(pac.hPass), pac)
	pac.onUpdated()
	return nil
}

func (pac *Packet) GetName() string {
	return pac.name
}
