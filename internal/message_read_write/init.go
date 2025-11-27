package messagereadwrite

import (
	"log"

	"github.com/gorilla/websocket"
)

const STRING_WRITE_ID = 1

type ReadWriter struct {
	received_counter int64
	sent_counter     int64
	writer_jobs      chan string

	closed bool
	done   chan bool
	conn   *websocket.Conn
}

type WriteRequest struct {
	req_type string
	data     string
}

func (rw *ReadWriter) IsClosed() bool {
	return rw.closed
}

func (rw *ReadWriter) Destruct() {
	rw.done <- true
	rw.closed = true
	rw.conn.Close()
}

func (rw *ReadWriter) closeHandler(code int, text string) error {
	rw.Destruct()
	log.Printf("websocket closed with code %d, and message: %s\n", code, text)
	return nil
}

func NewMessageReadWriter(conn *websocket.Conn) *ReadWriter {
	rw := ReadWriter{
		sent_counter:     0,
		received_counter: 0,
		conn:             conn,
		closed:           false,
		writer_jobs:      make(chan string),
		done:             make(chan bool),
	}
	conn.SetCloseHandler(rw.closeHandler)
	go func() {
		for {
			select {
			case wr := <-rw.writer_jobs:
				err := rw.conn.WriteMessage(STRING_WRITE_ID, []byte(wr))
				if err != nil {
					log.Printf("err")
					rw.Destruct()
					return
				}
			case <-rw.done:
				return
			}
		}
	}()
	return &rw
}
