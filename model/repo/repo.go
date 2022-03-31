package repo

import (
	"github.com/mmbednarek/fragma/model"
)

type ApiDetail interface {
	Name() string
	GetObjectDetail(name string) (model.ObjectDetail, error)
}

type ApiRepository struct {
	Apis map[string]ApiDetail
}

func (a *ApiRepository) RegisterApi(detail ApiDetail) {
	a.Apis[detail.Name()] = detail
}

func (a *ApiRepository) FindApi(name string) (ApiDetail, bool) {
	api, ok := a.Apis[name]
	if !ok {
		return nil, false
	}
	return api, ok
}
