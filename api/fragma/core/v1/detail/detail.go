package detail

import (
	"strings"

	core_v1 "github.com/mmbednarek/fragma/api/fragma/core/v1"
	"github.com/mmbednarek/fragma/model"
)

var objects = map[string]model.ObjectDetail{
	"application": {
		Version:           "v1",
		SingularName:      "application",
		PluralName:        "applications",
		FullName:          "fragma.core.v1.Application",
		ProtoType:         (&core_v1.Application{}).ProtoReflect().Type(),
		HighlightedFields: []string{"path", "name"},
	},
}

var aliases = map[string]string{
	"app":          "application",
	"apps":         "application",
	"applications": "application",
}

type ApiDetail struct {
}

func (ApiDetail) Name() string {
	return "fragma.core.v1"
}

func (d ApiDetail) GetObjectDetail(name string) (model.ObjectDetail, error) {
	nameLC := strings.ToLower(name)

	obj, ok := objects[nameLC]
	if !ok {
		alias, ok := aliases[nameLC]
		if ok {
			return d.GetObjectDetail(alias)
		}

		return model.ObjectDetail{}, model.ErrObjectNotFound
	}
	return obj, nil
}
