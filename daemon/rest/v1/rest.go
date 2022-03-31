package rest

import (
	"encoding/json"
	"fmt"

	"github.com/fasthttp/router"
	"github.com/mmbednarek/fragma/model"
	"github.com/valyala/fasthttp"
)

type ApiDetail interface {
	Name() string
	Objects() map[string]model.ObjectDetail
}

type CrudService interface {
	Update(obj *model.Object) error
	Read(typeName string, name string) (model.Object, error)
	Delete(typeName string, name string) error
	ReadAll(typeName string) ([]model.Object, error)
}

type Rest[TCrud CrudService] struct {
	apis []ApiDetail
	crud TCrud
}

func WithApi[TCrud CrudService](detail ApiDetail) func(rest *Rest[TCrud]) {
	return func(rest *Rest[TCrud]) {
		rest.apis = append(rest.apis, detail)
	}
}

func NewRest[TCrud CrudService](crud TCrud, opts ...func(rest *Rest[TCrud])) Rest[TCrud] {
	rest := Rest[TCrud]{
		crud: crud,
	}

	for _, opt := range opts {
		opt(&rest)
	}

	return rest
}

func (r *Rest[TCrud]) GetResource(ctx *fasthttp.RequestCtx, objectDetail model.ObjectDetail) {
	name := ctx.UserValue("name").(string)
	obj, err := r.crud.Read(objectDetail.FullName, name)
	if err != nil {
		ctx.Error("could not read object", fasthttp.StatusNotFound)
		return
	}

	result, err := json.Marshal(obj)
	if err != nil {
		ctx.Error("could not marshal object", fasthttp.StatusInternalServerError)
		return
	}

	if _, err := ctx.Write(result); err != nil {
		ctx.Error("could not write message", fasthttp.StatusInternalServerError)
		return
	}
}

func (r *Rest[TCrud]) WriteResource(ctx *fasthttp.RequestCtx, objectDetail model.ObjectDetail) {
	obj := model.Object{
		Spec: model.Spec{Message: objectDetail.ProtoType.New().Interface()},
	}

	if err := json.Unmarshal(ctx.PostBody(), &obj); err != nil {
		ctx.Error("could not unmarshall object", fasthttp.StatusNotFound)
		return
	}

	if err := r.crud.Update(&obj); err != nil {
		ctx.Error("could not read object", fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func (r *Rest[TCrud]) DeleteResource(ctx *fasthttp.RequestCtx, objectDetail model.ObjectDetail) {
	name := ctx.UserValue("name").(string)

	if err := r.crud.Delete(objectDetail.FullName, name); err != nil {
		ctx.Error("could not delete object", fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
}

func (r *Rest[TCrud]) GetAllResources(ctx *fasthttp.RequestCtx, objectDetail model.ObjectDetail) {
	objs, err := r.crud.ReadAll(objectDetail.FullName)
	if err != nil {
		ctx.Error("could not read objects", fasthttp.StatusNotFound)
		return
	}

	for _, obj := range objs {
		result, err := json.Marshal(obj)
		if err != nil {
			ctx.Error("could not marshal object", fasthttp.StatusInternalServerError)
			return
		}

		if _, err := ctx.Write(result); err != nil {
			ctx.Error("could not write message", fasthttp.StatusInternalServerError)
			return
		}
	}
}

func (r *Rest[TCrud]) RequestHandler() fasthttp.RequestHandler {
	rt := router.New()

	for _, api := range r.apis {
		objects := api.Objects()
		for _, object := range objects {
			rt.GET(fmt.Sprintf("/apis/%s/%s/{name}", api.Name(), object.PluralName), func(ctx *fasthttp.RequestCtx) {
				r.GetResource(ctx, object)
			})

			rt.DELETE(fmt.Sprintf("/apis/%s/%s/{name}", api.Name(), object.PluralName), func(ctx *fasthttp.RequestCtx) {
				r.DeleteResource(ctx, object)
			})

			rt.PUT(fmt.Sprintf("/apis/%s/%s", api.Name(), object.SingularName), func(ctx *fasthttp.RequestCtx) {
				r.WriteResource(ctx, object)
			})

			rt.GET(fmt.Sprintf("/apis/%s/%s", api.Name(), object.PluralName), func(ctx *fasthttp.RequestCtx) {
				r.GetAllResources(ctx, object)
			})

			rt.GET(fmt.Sprintf("/apis/%s/%s", api.Name(), object.SingularName), func(ctx *fasthttp.RequestCtx) {
				r.GetAllResources(ctx, object)
			})
		}
	}

	return rt.Handler
}
