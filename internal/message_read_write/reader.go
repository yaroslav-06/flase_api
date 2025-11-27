package messagereadwrite

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type APIRequest struct {
	Req  string          `json:"r"`
	Data json.RawMessage `json:"d"`
}

func (rw *ReadWriter) Reader() (string, *json.RawMessage, error) {
	_, p, err := rw.conn.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
			rw.Destruct()
		}
		return "", &json.RawMessage{}, err
	}
	fmt.Println(string(p))
	var dc APIRequest
	if err := json.Unmarshal(p, &dc); err != nil {
		return "", &json.RawMessage{}, fmt.Errorf("can't unmarshal api request Reader: %w", err)
	}

	return dc.Req, &dc.Data, nil
}
