package smartscheduler

import (
	"fmt"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"

func (sch *SmartScheduler) geti(i int64) (string, time.Time, error) {
	val, err := sch.ctx.GetString(fmt.Sprint(i)+":val")
	if err != nil{
		return "", time.Time{}, err
	}
	time_str, err := sch.ctx.GetString(fmt.Sprint(i)+":time")
	if err != nil{
		return "", time.Time{}, err
	}

	tm, err := time.Parse(TimeLayout, time_str)
	if err != nil{
		return "", time.Time{}, err
	}

	return val, tm, nil
}

func (sch *SmartScheduler) seti(i int64, val string, tm time.Time) {
	sch.ctx.SaveString(fmt.Sprint(i)+":val", val)
	time_str := tm.Format(TimeLayout)
	sch.ctx.SaveString(fmt.Sprint(i)+":time", time_str)
	sch.ctx.SaveMapField("vals", val, fmt.Sprint(i))
}

/*
Add an event (of $val$) that should be executed at $time$.
Addes to the end and swaps forward if needed (min heap).
*/
func (sch *SmartScheduler) AddEvent(val string, tm time.Time) (error) {
	prev_i := sch.len + 1
	sch.seti(prev_i, val, tm)
	for i := prev_i / 2; i > 0; i /= 2 {
		ul, tm_h, err := sch.geti(i)
		if err != nil{return err}
		// fmt.Printf("%d) here: %s, cmp: %s\n", i, tm_h, tm)
		if(tm.Before(tm_h)){
			sch.seti(prev_i, ul, tm_h)
		}else{
			break
		}
		prev_i = i
	}
	// fmt.Printf("set at: %d\n", prev_i)
	sch.seti(prev_i, val, tm)
	sch.len += 1
	sch.ctx.SaveString("len", fmt.Sprint(sch.len))
	return nil
}
