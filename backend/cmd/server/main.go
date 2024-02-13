package main

import (
	"log"

	"github.com/babadro/tutor/internal/infra/restapi"
	"github.com/babadro/tutor/internal/infra/restapi/operations"
	"github.com/caarlos0/env"
	"github.com/go-openapi/loads"
)

type envVars struct {
	AppPort int `env:"APP_PORT,required"`
}

func main() {
	var envs envVars
	if err := env.Parse(&envs); err != nil {
		log.Fatalf("Unable to parse env vars: %v\n", err)
	}

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewTutorAPI(swaggerSpec)

	server := restapi.NewServer(api)
	defer func(server *restapi.Server) {
		if err = server.Shutdown(); err != nil {
			log.Printf("error while shutting down server: %v", err)
		}
	}(server)

	server.Port = envs.AppPort

	server.ConfigureAPI()

	if err = server.Serve(); err != nil {
		log.Printf("error while serving: %v", err)
	}
}
