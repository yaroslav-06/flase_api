package smartscheduler

import (
	"fmt"
	"strconv"
)

func (sch *SmartScheduler) RemoveEvent(val string) error {
	i_str, err := sch.ctx.GetMapField("vals", val)
	if err != nil{ return nil }

	i, err := strconv.ParseInt(i_str, 10, 64)
	if err != nil{ return nil }

	val, tm, err := sch.geti(sch.len)
	if err != nil{ return nil }
	fmt.Printf("inserting %s\n", val)

	for 2*i <= sch.len {
		fmt.Printf("new i: %d\n", i)
		v := 2*i; u := 2*i + 1;

		v_vl, v_tm, err := sch.geti(v)
		if err != nil{ return nil }
		u_vl, u_tm, err := sch.geti(u)
		if(u <= sch.len){
			if err != nil{ return nil }
		}

		if(tm.Before(v_tm) && tm.Before(u_tm)){
			break
		}

		if(u > sch.len || v_tm.Before(u_tm)){
			sch.seti(i, v_vl, v_tm)
			i = v
		}else{
			sch.seti(i, u_vl, u_tm)
			i = u
		}
	}
	fmt.Printf("%d)inserting %s\n", i, val)
	sch.seti(i, val, tm)
	sch.len -= 1
	sch.ctx.SaveString("len", fmt.Sprint(sch.len))
	return nil
}
