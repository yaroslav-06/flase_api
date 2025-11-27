package socket

import (
	"flase_api/internal/auth"
	"flase_api/internal/db"
	"flase_api/internal/server/handlr"
	"flase_api/internal/session"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 2,
	WriteBufferSize: 1024 * 2,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request, ctx *db.DbCtx, authSys *auth.AuthSystem, controller *session.SessionController) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if ws != nil {
		fmt.Println("connected")
		wsHandler(controller.NewSession(ws, authSys), ctx)
	}
}

func wsHandler(sess *session.Session, ctx *db.DbCtx) {
	hdlrs := handlr.Get()
	for {
		req, data, err := sess.GetRW().Reader()
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("request \"%s\" data: \"%s\"\n", req, data)
		hdlr, exists := (*hdlrs)[req]
		if !exists {
			log.Printf("request: %s doesn't exist\n", req)
			continue
		}
		fmt.Println("handling: ", req)
		hdlr.HandleRequest(ctx, sess, data)
	}
}

func Init(port string, ctx *db.DbCtx, authSys *auth.AuthSystem, controller *session.SessionController) {
	log.Printf("starting server at port %s\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, ctx, authSys, controller)
	})

	err := http.ListenAndServe(":"+port, nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
	} else if err != nil {
		log.Printf("unknown error while starting server %s\n", err)
		os.Exit(1)
	} else {
		log.Printf("server started at the port %s\n", port)
	}
}
