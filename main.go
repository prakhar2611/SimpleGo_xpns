package main

import (
	Controller "SimpleGo_xpns/Controllers"
	utilities "SimpleGo_xpns/Utilities"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	port := flag.String("port", "9005", "Service active Port number, default: 9999")
	//configFile := flag.String("config", "config", "Config JSON file name,default: config")
	fmt.Println("Hello world !")
	// utilities.ConnectDB()

	utilities.Load()

	fmt.Println("Application running on port", *port)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*", "http://localhost:3001"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	//Registering all the routes for the applications
	r.Route("/expense", func(r chi.Router) {

		Controller.RegisterDataAPI(r)
	})
	r.Route("/docs", func(r chi.Router) {

		Controller.RegisterDocsAPI(r)
	})
	r.Route("/", func(r chi.Router) {
		Controller.RegisterUserAPI(r)
		Controller.RegisterGoogleAPIs(r)
	})
	log.Println(http.ListenAndServe("0.0.0.0:"+*port, r))
}
