package service

import (
	"fmt"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/hudl/fargo"
	"github.com/unrolled/render"
)

func NewServer(connection *fargo.EurekaConnection) *negroni.Negroni {
	webClient := fulfillmentWebClient{
		rootURL: "http://localhost:3001/skus",
	}

	app, err := connection.GetApp("BACKING_FULFILLMENT")
	if err == nil {
		instance := app.Instances[0]
		webClient.rootURL = instance.HostName + "/skus"
	} else {
		fmt.Printf("Failed to get registered URL from Eureka: %v\n", err)
	}
	fmt.Printf("Using the following URL for fulfillment backing service: %s\n", webClient.rootURL)

	return NewServerFromClient(webClient)
}

// NewServerFromClient configures and returns a Server.
func NewServerFromClient(webClient fulfillmentClient) *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter, webClient)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render, webClient fulfillmentClient) {
	mx.HandleFunc("/", rootHandler(formatter)).Methods("GET")
	mx.HandleFunc("/catalog", getAllCatalogItemsHandler(formatter)).Methods("GET")
	mx.HandleFunc("/catalog/{sku}", getCatalogItemDetailsHandler(formatter, webClient)).Methods("GET")
}
