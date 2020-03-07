package store

var implStore Store
var implInit func() error

func Init() error {
	return implInit()
}

func Register(initFunc func() error, store Store) {
	implInit = initFunc
	implStore = store
}

type Store interface {
	APIs(uiModule string) []ApiDef
}

type ApiDef struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Type string `json:"type"`
}

func APIs(uiModule string) []ApiDef {
	return implStore.APIs(uiModule)
}
