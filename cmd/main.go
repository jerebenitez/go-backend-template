package main

import (
	"fmt"
	"log"

	"github.com/jerebenitez/go-backend-template/cmd/api"
	"github.com/jerebenitez/go-backend-template/utils"
)

func main() {
	pool, ctx, err := utils.NewPool(utils.Envs.DB)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(
		fmt.Sprintf("%s:%s", utils.Envs.PublicHost, utils.Envs.Port),
		pool, ctx,
	)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
