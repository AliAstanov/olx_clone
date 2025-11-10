package main

import (
	"fmt"
	"log"
	"time"

	"github.com/AliAstanov/olx_clone/api"
	"github.com/AliAstanov/olx_clone/config"
	"github.com/AliAstanov/olx_clone/pkg/db"
	"github.com/AliAstanov/olx_clone/storage"
)

func main() {

	cfg := config.Load()

	db, err := db.ConnToDb(cfg.PgConfig)
	if err != nil {
		log.Println("error on connect to ConToDb:", err)
		return
	}

	storage := storage.NewStorage(db)
	var time time.Time
	var bol bool
	fmt.Println("time=",time)
	fmt.Println("bool=",bol)

	api.Api(storage)
}
