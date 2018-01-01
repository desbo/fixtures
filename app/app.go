package app

import (
	"log"
	"net/http"

	loads "github.com/go-openapi/loads"

	"github.com/desbo/fixtures/restapi"
	"github.com/desbo/fixtures/restapi/operations"
)

func init() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	var server *restapi.Server

	api := operations.NewTabletennis365ComFixturesAPI(swaggerSpec)
	server = restapi.NewServer(api)

	defer server.Shutdown()

	server.ConfigureAPI()
	http.Handle("/", server.GetHandler())
}
