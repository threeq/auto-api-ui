package local

import (
	"auto-api-ui/store"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type localStore struct {
	Meta []*store.DocDesc `json:"meta"`
}

func (l *localStore) UpSetAPI(name string, desc *store.DocDesc) error {
	rwmux.Lock()
	defer rwmux.Unlock()

	// Add api
	if name == "" {
		existApi := findApiDesc(desc.Name)
		if existApi == nil {
			data.Meta = append(data.Meta, desc)
		}
	} else {
		existApi := findApiDesc(name)
		if existApi == nil {
			data.Meta = append(data.Meta, desc)
		} else {
			existApi.Name = desc.Name
			existApi.Type = desc.Type
			existApi.Url = desc.Url
		}
	}
	_ = writeData(data)
	return nil
}

func (l *localStore) DelAPI(name string) error {
	rwmux.Lock()
	defer rwmux.Unlock()
	panic("implement me")
}

func (l *localStore) Watch() store.Stop {
	stop := make(chan struct{})
	stopWait := make(chan struct{})
	storePath := path.Join(dataPath, fileName)
	log.Printf("start watch data file %s", storePath)

	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			watcher.Close()
			stopWait <- struct{}{}
			log.Print("stop watch data")
		}()
		err = watcher.Add(storePath)
		if err != nil {
			log.Fatal(err)
		}

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					err := safeReadData()

					if err != nil {
						log.Printf("read data fail：%v", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)

			case <-stop:
				return
			}
		}
	}()
	return func() chan struct{} {
		stop <- struct{}{}
		return stopWait
	}
}

func (l *localStore) APIs(uiModule string) []*store.DocDesc {
	rwmux.RLock()
	defer rwmux.RUnlock()
	return data.Meta
}

func safeReadData() error {
	rwmux.Lock()
	defer rwmux.Unlock()

	storePath := path.Join(dataPath, fileName)
	metaFile, err := os.Open(storePath)
	if err != nil {
		return err
	}
	defer metaFile.Close()

	loadData, err := ioutil.ReadAll(metaFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(loadData, &data)
	if err != nil {
		return err
	}

	log.Printf("load document description from %s", storePath)
	return nil
}

func writeData(data *localStore) error {
	metaPath := path.Join(dataPath, fileName)
	metaFile, err := os.Create(metaPath)
	if err != nil {
		return err
	}
	defer metaFile.Close()

	bytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	_, err = metaFile.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func deepCopy(data *localStore) *localStore {
	var buff bytes.Buffer
	cp := &localStore{
		Meta: []*store.DocDesc{},
	}
	gob.NewEncoder(&buff).Encode(*data)
	gob.NewDecoder(bytes.NewBuffer(buff.Bytes())).Decode(cp)
	return cp
}

func findApiDesc(name string) *store.DocDesc {
	for _, desc := range data.Meta {
		if desc.Name == name {
			return desc
		}
	}

	return nil
}
