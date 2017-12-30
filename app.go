package main

import (
	"log"
	"net/http"
	"os"

	loads "github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"

	"github.com/desbo/fixtures/restapi"
	"github.com/desbo/fixtures/restapi/operations"
)

func init() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewTabletennis365ComFixturesAPI(swaggerSpec)
	server := restapi.NewServer(api)

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "tabletennis365.com fixtures"
	parser.LongDescription = "tabletennis365.com fixtures"

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI()

	http.Handle("/", server.GetHandler())
}
