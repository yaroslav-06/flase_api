package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// change if you are running the api_telegram_sender on different server
const api_telegram_sender_ip = "0.0.0.0"

// insert string from your_api_key variable in api_telegram_sender/main.py
const your_api_key = "my_key"

func SendMsg(usrname string, msg string) error {
	data := map[string]string{
			"username": usrname,
			"message": 	msg,
			"api_key":	your_api_key,
	}

	jsonData, err := json.Marshal(data)
	if err != nil{
		return err
	}

	resp, err := http.Post(
		"http://0.0.0.0:8000/send",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}else{
		return fmt.Errorf("api status code: %d (configure your_api_key at intrnal/telegram/send.go)", resp.StatusCode)
	}
}
