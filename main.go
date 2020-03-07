package main

import (
	"auto-api-ui/store"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)
import _ "auto-api-ui/store/local"

func main() {
	CmdArgsParse()

	// init data
	err := store.Init()
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	router := httprouter.New()
	initUI(router)

	log.Printf("start server at %s", Args.Addr)
	err = http.ListenAndServe(Args.Addr, router)
	if err != nil {
		log.Println(err)
	}
}
