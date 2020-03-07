package local

import (
	"auto-api-ui/store"
	"auto-api-ui/util"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const (
	fileName = "data.json"
)

var (
	dataPath string
	data     = localStore{
		Meta: []store.ApiDef{},
	}
)

func init() {
	store.Register(initStore, &data)
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
	var metaFile *os.File

	if !util.Exists(metaPath) {
		metaFile, err = os.Create(metaPath)
		if err != nil {
			return err
		}
		defer metaFile.Close()
		data, err := json.Marshal(&data)
		if err != nil {
			return err
		}
		_, err = metaFile.Write(data)
		if err != nil {
			return err
		}
	} else if util.IsDir(metaPath) {
		return errors.New(fmt.Sprintf("%s is directory not file", metaPath))
	} else {
		metaFile, _ = os.Open(metaPath)
		defer metaFile.Close()
		loadData, err := ioutil.ReadAll(metaFile)
		if err != nil {
			return err
		}
		err = json.Unmarshal(loadData, &data)
		if err != nil {
			return err
		}
	}

	return nil
}
