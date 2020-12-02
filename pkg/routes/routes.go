package routes

import (
	"shopping-cart/pkg/controllers/cart"
	"shopping-cart/pkg/controllers/user"

	"github.com/gorilla/mux"
)

func SetRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/v1/login", user.Login).Methods("GET")
	router.HandleFunc("/v1/register", user.RegisterUser).Methods("POST")
	router.HandleFunc("/v1/logout", user.LogOut).Methods("GET")

	router.HandleFunc("/v1/cart", cart.ViewCart).Methods("GET")
	router.HandleFunc("/v1/cart", cart.AddItem).Methods("PATCH")
	router.HandleFunc("/v1/cart", cart.ClearCart).Methods("DELETE")
	router.HandleFunc("/v1/cart/{itemid}", cart.RemoveItem).Methods("DELETE") 
	return router
}