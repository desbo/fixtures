package main

import (
	"log"
	"net/http"
	"os"

	loads "github.com/go-openapi/loads"

	"github.com/rs/cors"

	"github.com/desbo/fixtures/restapi"
	"github.com/desbo/fixtures/restapi/operations"
)

func main() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	var server *restapi.Server

	api := operations.NewTabletennis365ComFixturesAPI(swaggerSpec)
	server = restapi.NewServer(api)

	defer server.Shutdown()

	server.ConfigureAPI()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	http.Handle("/", cors.AllowAll().Handler(server.GetHandler()))

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
