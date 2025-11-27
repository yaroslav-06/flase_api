package recievers

import (
	"flase_api/internal/db"
	"encoding/json"
	"fmt"
)

type emailHandler struct {
	RecieverHandler
}

type emailReciever struct {
	Reciever
	Email   string `json:"email"`
	Message string `json:"message"`
}

func (hd *emailHandler) LoadReciever(data *json.RawMessage) Reciever {
	rc := emailReciever{}
	json.Unmarshal(*data, &rc)
	return &rc
}
func (hd *emailHandler) DbLoad(ctx *db.DbCtx) (Reciever, error) {
	rc := emailReciever{}
	return &rc, nil
}

func (rc *emailReciever) Send() error {
	fmt.Printf("email send: %s\n", rc.Message)
	return nil
}

func (rc *emailReciever) DbSave(ctx *db.DbCtx) error {
	return nil
}
