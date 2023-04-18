package main

import (
	Controller "SimpleGo_xpns/Controllers"
	service "SimpleGo_xpns/Services"
	utilities "SimpleGo_xpns/Utilities"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	port := flag.String("port", "9005", "Service active Port number, default: 9999")
	//configFile := flag.String("config", "config", "Config JSON file name,default: config")
	fmt.Println("Hello world !")
	// utilities.ConnectDB()

	utilities.Load()

	fmt.Println("Application running on port", *port)

	r := chi.NewRouter()

	//Registering all the routes for the applications
	r.Route("/expense", func(r chi.Router) {
		Controller.RegisterUserAPI(r)
		Controller.RegisterDataAPI(r)
	})
	r.Route("/", func(r chi.Router) {
		service.RegisterGoogleAPIs(r)
	})
	log.Println(http.ListenAndServe("0.0.0.0:"+*port, r))
}
