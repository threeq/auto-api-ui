package main

import (
	"auto-api-ui/store"
	"auto-api-ui/util"
	"context"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"time"
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
	initApi(router)
	printRouters(router)

	serve := &http.Server{Addr: Args.Addr, Handler: router}

	go func() {
		log.Printf("start server at %s", Args.Addr)

		if err := serve.ListenAndServe(); err != nil {
			log.Fatalf("%v", err)
		}
	}()

	// 处理系统退出信号
	util.ObserveExitSignal(func(signal os.Signal) {
		log.Printf("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := serve.Shutdown(ctx); err != nil {
			log.Printf("Server Shutdown error: %v", err)
		}
		_ = store.Close(ctx)

		log.Printf("Server exited")
	})
}

func printRouters(router *httprouter.Router) {
}
