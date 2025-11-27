package main

import (
	"flase_api/internal/auth"
	"flase_api/internal/db"
	socket "flase_api/internal/server"
	"flase_api/internal/session"
	"context"
	"fmt"
	"log"
)

func main() {
	fmt.Println("flase_api server")
	ctx, err := db.GetParentContext("6379", context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	controller := session.InitController(ctx)

	auth, err := auth.InitAuthSystem(ctx)
	if err != nil{
		log.Printf("auth error: %s\n", err.Error())
		return
	}


	socket.Init("14539", ctx, auth, controller)
}
