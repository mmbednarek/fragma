package main

import (
	"context"
	"os"

	core "github.com/mmbednarek/fragma/api/fragma/core/v1"
	"github.com/mmbednarek/fragma/daemon/service/v1"
	"github.com/mmbednarek/fragma/pkg/log"
	_ "github.com/mmbednarek/fragma/pkg/log/formatter"
)

func main() {
	ctx := context.Background()
	srv := service.NewService()

	img := os.Getenv("FRAGMA_IMAGE")
	if len(img) == 0 {
		img = "./img"
	}

	bin := "/usr/bin/bash"
	if len(os.Args) > 1 {
		bin = os.Args[1]
	}

	app := &core.Application{
		Name: "Bash SHELL",
		Path: bin,
	}

	volume := &core.Volume{
		Path: img,
	}

	options := &core.RunOptions{
		Arguments: os.Args[1:],
	}

	log.With(ctx).Info("starting application")
	if err := srv.RunApplication(ctx, volume, app, options); err != nil {
		log.With(ctx, "msg", err).Error("could not run application")
	}
}
