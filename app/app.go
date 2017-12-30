package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	loads "github.com/go-openapi/loads"
	flag "github.com/spf13/pflag"

	"github.com/desbo/fixtures/restapi"
	"github.com/desbo/fixtures/restapi/operations"
)

func init() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	var server *restapi.Server // make sure init is called

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage:\n")
		fmt.Fprint(os.Stderr, "  tabletennis365-com-fixtures-server [OPTIONS]\n\n")

		title := "tabletennis365.com fixtures"
		fmt.Fprint(os.Stderr, title+"\n\n")
		desc := "tabletennis365.com fixtures"
		if desc != "" {
			fmt.Fprintf(os.Stderr, desc+"\n\n")
		}
		fmt.Fprintln(os.Stderr, flag.CommandLine.FlagUsages())
	}
	// parse the CLI flags
	flag.Parse()

	api := operations.NewTabletennis365ComFixturesAPI(swaggerSpec)
	// get server with flag values filled out
	server = restapi.NewServer(api)

	defer server.Shutdown()

	server.ConfigureAPI()
	http.Handle("/", server.GetHandler())
}
