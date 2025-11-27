package smartscheduler

import (
	"flase_api/internal/db"
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

type tstExec struct {
	SchedulerEventType
}

func (ex tstExec) ExecuteTask(val string) error {
	log.Println("executing")
	return nil
}

func TestInitSmartScheduler(t *testing.T){
	//get parent ctx
	ctx, err := db.GetParentContext("6379", context.Background())
	if err != nil{
		t.Errorf("Couldn't get the inital ctx: %e", err)
		return
	}

	tst := ctx.GetChild("test")
	ctx = tst.GetChild("smartscheduler")
	executor := tstExec{}
	scheduler, err := InitSmartScheduler(ctx, executor)
	if err != nil{
		t.Errorf("Couldn't do init scheduler: %e", err)
		return
	}
	if scheduler.ctx != ctx {
		t.Errorf("ctx is wrong")
		return
	}
	/// Test addition
	t1 := time.Now().UTC().Add(time.Second * 3)
	t2 := t1.Add(time.Second)
	t3 := t1.AddDate(0, 1, 2)
	t4 := t3.AddDate(0, 0, 2)
	t6 := t3.Add(time.Minute * 3)
	t7 := t3.Add(time.Minute * 4)
	fmt.Printf("t4: %s", t4)
	scheduler.AddEvent("e_1", t2)
	if scheduler.len != int64(1){
		t.Errorf("length incorrect\n")
	}
	scheduler.AddEvent("e_2", t4)
	scheduler.AddEvent("e_3", t1)
	scheduler.AddEvent("e_4", t3)
	scheduler.AddEvent("e_5", t3.AddDate(1, 0, 0))
	scheduler.AddEvent("e_6", t6)
	scheduler.AddEvent("e_7", t7)
	expected := []string{"e_3", "e_4", "e_1", "e_2", "e_5", "e_6", "e_7"}
	if scheduler.len != int64(len(expected)){
		t.Errorf("length incorrect\n")
		return
	}

	for i:=int64(1); i <= int64(len(expected)); i++ {
		v, _, err := scheduler.geti(i)
		if err != nil{
			t.Errorf("had an err: %e", err)
			return
		}
		log.Println(v)
		if expected[i - 1] != v{
			t.Errorf("expected %s, got %s", expected[i], v)
			return
		}
	}

	expected = []string{"e_1", "e_4", "e_6", "e_2", "e_5", "e_7"}
	scheduler.RemoveEvent("e_3")
	for i:=int64(1); i <= int64(len(expected)); i++ {
		v, _, err := scheduler.geti(i)
		if err != nil{
			t.Errorf("had an err: %e", err)
			return
		}
		log.Println(v)
		if expected[i - 1] != v{
			t.Errorf("expected %s, got %s", expected[i], v)
			return
		}
	}
	// time.Sleep(time.Minute * 2)
}

