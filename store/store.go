package store

import (
	"context"
	"flag"
)

var storeService Store
var initService func() error
var stopWatch Stop

var storeWatch bool

func init() {
	flag.BoolVar(&storeWatch, "store-watch", false, "主动监视存储的改变")
}

func Init() error {
	err := initService()
	if err != nil {
		return err
	}
	stopWatch = storeService.Watch()
	return nil
}

func Register(initFunc func() error, store Store) {
	initService = initFunc
	storeService = store
}

type Stop = func() chan struct{}

type Store interface {
	// 返回所有 api 定义
	APIs(uiModule string) []*DocDesc
	Watch() Stop
	UpSetAPI(name string, desc *DocDesc) error
	DelAPI(name string) error
}

type DocDesc struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Type string `json:"type"`
}

func APIs(uiModule string) []*DocDesc {
	return storeService.APIs(uiModule)
}

func Close(ctx context.Context) error {
	select {
	case <-stopWatch():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func DelAPI(name string) error {
	return storeService.DelAPI(name)
}

func UpdateAPI(name string, apiDesc *DocDesc, ) error {
	return storeService.UpSetAPI(name, apiDesc)
}

func AddAPI(apiDesc *DocDesc) error {
	return storeService.UpSetAPI("", apiDesc)
}
