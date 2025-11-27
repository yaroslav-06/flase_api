package session

import (
	"flase_api/internal/auth"
	messagereadwrite "flase_api/internal/message_read_write"

	"github.com/gorilla/websocket"
)

type Session struct {
	id       string
	uid      string
	packetId string
	rw       *messagereadwrite.ReadWriter

	usr        *auth.User
	controller *SessionController
	authSystem *auth.AuthSystem
}

func (controller *SessionController) NewSession(conn *websocket.Conn, authSys *auth.AuthSystem) *Session {
	sess := &Session{
		id:         controller.GetNewId(),
		uid:        "",
		packetId:   "",
		rw:         messagereadwrite.NewMessageReadWriter(conn),
		controller: controller,
		authSystem: authSys,
	}
	return sess
}

func (sess *Session) GetAuthSystem() *auth.AuthSystem {
	return sess.authSystem
}

func (sess *Session) GetSessController() *SessionController {
	return sess.controller
}

func (sess *Session) GetUid() *string {
	if sess.uid != "" {
		return &sess.uid
	}
	if sess.usr == nil {
		return nil
	}
	sess.uid = sess.usr.GetUid()
	return &sess.uid
}

func (sess *Session) GetId() *string {
	if sess.id == "" {
		return nil
	}
	return &sess.id
}

func (sess *Session) GetRW() *messagereadwrite.ReadWriter {
	return sess.rw
}

func (sess *Session) SetUser(usr *auth.User) {
	sess.usr = usr
}

func (sess *Session) GetUser() *auth.User {
	return sess.usr
}
