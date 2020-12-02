package main

import (
	"log"
	"net/http"
	"shopping-cart/pkg/database"
	"shopping-cart/pkg/routes"

	"github.com/gorilla/mux"
)
 

func main() {

	database.Connect()
	router := mux.NewRouter() 
	router = routes.SetRoutes(router)
	log.Fatal(http.ListenAndServe(":8008", router))
}
