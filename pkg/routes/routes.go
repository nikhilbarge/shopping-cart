package routes

import (
	"shopping-cart/pkg/controllers/cart"
	"shopping-cart/pkg/controllers/inventory"
	"shopping-cart/pkg/controllers/user"

	"github.com/gorilla/mux"
)

func SetRoutes(router *mux.Router) *mux.Router {

	router.HandleFunc("/v1/info", user.AppDetails).Methods("GET")
	
	router.HandleFunc("/v1/register", user.RegisterUser).Methods("POST")
	router.HandleFunc("/v1/login", user.Login).Methods("POST")
	router.HandleFunc("/v1/logout", user.LogOut).Methods("POST")

	router.HandleFunc("/v1/inventory", inventory.ViewInventory).Methods("GET")
	router.HandleFunc("/v1/inventory", inventory.AddItemToInventory).Methods("POST")
	router.HandleFunc("/v1/inventory", inventory.RemoveItem).Methods("DELETE")
	router.HandleFunc("/v1/inventory{itemid}", inventory.RemoveItem).Methods("DELETE")

	router.HandleFunc("/v1/cart", cart.ViewCart).Methods("GET")
	router.HandleFunc("/v1/cart", cart.AddItem).Methods("PATCH")
	router.HandleFunc("/v1/cart", cart.ClearCart).Methods("DELETE")
	router.HandleFunc("/v1/cart/{itemid}", cart.RemoveItem).Methods("DELETE") 
	return router
}