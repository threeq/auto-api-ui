package local

import (
	"auto-api-ui/store"
	"auto-api-ui/util"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"sync"
)

const (
	fileName = "data.json"
)

var (
	dataPath string
	data     = &localStore{
		Meta: []*store.DocDesc{},
	}
	rwmux = sync.RWMutex{}
)

func init() {
	store.Register(initStore, data)
	flag.StringVar(&dataPath, "data", "data", "数据存储路径")
}

func initStore() (err error) {
	if !util.Exists(dataPath) {
		err = os.MkdirAll(dataPath, 0644)
		if err != nil {
			return err
		}
	} else if util.IsFile(dataPath) {
		return errors.New(fmt.Sprintf("%s is file not directory", dataPath))
	}
	metaPath := path.Join(dataPath, fileName)

	if !util.Exists(metaPath) {
		return writeData(data)
	} else if util.IsDir(metaPath) {
		return errors.New(fmt.Sprintf("%s is directory not file", metaPath))
	}
	return safeReadData()

}
