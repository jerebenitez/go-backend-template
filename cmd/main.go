package main

import (
	"fmt"
	"log"

	"github.com/jerebenitez/go-backend-template/cmd/api"
	"github.com/jerebenitez/go-backend-template/utils"
)

func main() {
	server := api.NewAPIServer(
		fmt.Sprintf("%s:%s", utils.Envs.PublicHost, utils.Envs.Port), 
		nil,
	)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
