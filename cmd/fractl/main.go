package main

import (
	"github.com/mmbednarek/fragma/model/client"
	"github.com/mmbednarek/fragma/model/repo"
)

func main() {
	frontend := NewFrontend(client.NewClient("127.0.0.1:8000"), repo.GetStandardRepository())
	rootCmd := frontend.Mount()
	if err := rootCmd.Execute(); err != nil {
		die(err.Error())
	}
}
