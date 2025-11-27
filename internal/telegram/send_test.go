package telegram

import (
	"testing"
)


func TestSendTelegram(t *testing.T){
	err := SendMsg("@parolk06", "successful test of flase_api")
	if err != nil{
		t.Errorf("error: %s", err)
	}
}
