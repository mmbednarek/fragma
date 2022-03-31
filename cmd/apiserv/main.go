package main

import (
	"log"

	core_v1_det "github.com/mmbednarek/fragma/api/fragma/core/v1/detail"
	"github.com/mmbednarek/fragma/daemon/rest/v1"
	"github.com/mmbednarek/fragma/model"
	"github.com/mmbednarek/fragma/pkg/storage"
	"github.com/valyala/fasthttp"
)

type Crud = *model.CrudService[storage.Storage]

func main() {
	store, err := storage.NewStorage("/tmp/fragmastore")
	if err != nil {
		log.Fatalf("storage.NewStorage: %s", err)
	}

	crud := model.NewCrudService[storage.Storage](store)

	restApi := rest.NewRest[Crud](&crud,
		rest.WithApi[Crud](core_v1_det.ApiDetail{}),
	)
	handler := restApi.RequestHandler()

	if err := fasthttp.ListenAndServe("127.0.0.1:8000", handler); err != nil {
		log.Fatalf("fasthttp.ListenAndServe: %s", err)
	}
}
