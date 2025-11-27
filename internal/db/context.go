package db

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
)

/*
Object that represents a reference to id in database.

Constructors:
 - GetParentContext -- for initial redis connection, and base reference
 - GetChild(exten) -- to get/create a child of this DbCtx with given exten
 - GetParent -- to get the parent of this DbCtx
*/
type DbCtx struct {
	rp   Repo
	path string
}

func GetParentContext(port string, cntx context.Context) (*DbCtx, error) {
	ctx := &DbCtx{}
	rp, err := Connect(port, cntx)
	if err != nil {
		return nil, err
	}
	ctx.rp = rp
	ctx.path = ""
	return ctx.GetChild("fapi"), nil
}

func (parent DbCtx) GetChild(pathExtention string) *DbCtx {
	ctx := &DbCtx{}
	ctx.rp = parent.rp
	ctx.path = parent.path + ":" + pathExtention
	return ctx
}

func (child DbCtx) GetParent() *DbCtx {
	ctx := &DbCtx{}
	ctx.rp = child.rp
	lst_ind := strings.LastIndex(child.path, ":")
	ctx.path = child.path[0:lst_ind]
	return ctx
}

func (nd DbCtx) SaveString(name string, data string) error {
	fmt.Printf("[%s] = %s\n", nd.path+":"+name, data)
	vl := nd.rp.Client.Set(nd.rp.ctx, nd.path+":"+name, data, 0)
	return vl.Err()
}

func (ctx DbCtx) GetInt(name string) (int64, error) {
	rq := ctx.rp.Client.Get(ctx.rp.ctx, ctx.path+":"+name)
	if rq.Err() != nil {
		return 0, rq.Err()
	}
	rs, err := strconv.ParseInt(rq.Val(), 10, 64)
	if err != nil{
		return 0, err
	}
	return rs, nil
}

func (ctx DbCtx) GetString(name string) (string, error) {
	rq := ctx.rp.Client.Get(ctx.rp.ctx, ctx.path+":"+name)
	if rq.Err() != nil {
		return "", rq.Err()
	}
	return rq.Val(), nil
}

func (nd DbCtx) SaveMap(name string, data map[string]string) {
	nd.rp.Client.HSet(nd.rp.ctx, nd.path+":"+name, data, 0)
}

func (nd DbCtx) SaveMapField(name string, field string, val string) error {
	log.Printf("[%s] = %s\n", field, val)
	vl := nd.rp.Client.HSet(nd.rp.ctx, nd.path+":"+name, field, val)
	return vl.Err()
}

func (ctx DbCtx) GetMapField(name string, field string) (string, error) {
	rq := ctx.rp.Client.HGet(ctx.rp.ctx, ctx.path+":"+name, field)
	if rq.Err() != nil {
		return "", rq.Err()
	}
	return rq.Val(), nil
}

func (ctx DbCtx) GetMapAll(name string) (map[string]string, error) {
	rq := ctx.rp.Client.HGetAll(ctx.rp.ctx, ctx.path+":"+name)
	if rq.Err() != nil {
		return nil, rq.Err()
	}
	return rq.Val(), nil
}

func (nd DbCtx) ClearField(name string) error {
	log.Printf("delete [%s]\n", name)
	vl := nd.rp.Client.Del(nd.rp.ctx, nd.path+":"+name)
	return vl.Err()
}

func (nd *DbCtx) GetPath() string {
	return nd.path
}
