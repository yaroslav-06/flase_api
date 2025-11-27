package smartscheduler

import (
	"flase_api/internal/db"
	"fmt"
	"log"
	"time"
)

/*
Implementation of Scheduling structure.
Which loads from memory, and saves all events in memory.
It will run as a go rootine, and execute the event if current
datetime is larger then the defined execution datetime.
Each check is O(1), using min-heap.

Supports:
 - addition of events: O(log N)
 - removal of events: O(log N)

Constructor: InitSmartScheduler
*/
type SmartScheduler struct {
	ctx 			*db.DbCtx
	len				int64
	executor 	SchedulerEventType
}

func (sch *SmartScheduler) createSmartScheduler(ctx *db.DbCtx) error {
	err := ctx.SaveString("created", "yes")
	if err != nil{
		return err
	}
	err = ctx.SaveString("len", "0")
	if err != nil{
		return err
	}
	ctx.SaveMap("values", map[string] string {})
	return nil
}

func InitSmartScheduler(ctx *db.DbCtx, executor SchedulerEventType) (*SmartScheduler, error) { 
	log.Println("init scheduler")
	rs := &SmartScheduler{}
	rs.ctx = ctx
	rs.executor = executor
	iscr, err := ctx.GetString("created")
	if err != nil{
		err = rs.createSmartScheduler(ctx)
		if err != nil{
			return nil, err
		}
		iscr = "yes"
	}
	if iscr != "yes"{
		return nil, fmt.Errorf("the created not yes %s", err)
	}

	len, err := ctx.GetInt("len")
	if err != nil{
		return nil, fmt.Errorf("smartscheduler length err: %s", err)
	}
	rs.len = len
		
	log.Println(iscr)
	go rs.run()
	return rs, nil
}

func (sch *SmartScheduler) run() {
	for {
		time.Sleep(time.Second * 2)
		if (sch.len != 0){
			vl, tm, err := sch.geti(1)
			// log.Printf("othr: %s", tm)
			// log.Printf("cur: %s", time.Now())
			if err != nil{
				log.Println(err)
				continue
			}
			if(tm.Before(time.Now())){
				log.Println("execute")
				err = sch.executor.ExecuteTask(vl)
				sch.RemoveEvent(vl)
				if err != nil{
					log.Println(err, vl)
					continue
				}
			}
		}
	}	
}
