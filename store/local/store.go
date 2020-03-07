package local

import (
	"auto-api-ui/store"
)

type localStore struct {
	Meta []store.ApiDef `json:"meta"`
}

func (l localStore) APIs(uiModule string) []store.ApiDef {
	return data.Meta
}
