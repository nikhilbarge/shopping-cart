package routes

import (
	"shopping-cart/pkg/controllers/cart"
	"shopping-cart/pkg/controllers/categories"
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
 
	router.HandleFunc("/v1/categories", categories.ViewCategories).Methods("GET")
	router.HandleFunc("/v1/categories", categories.AddCategories).Methods("POST")
	router.HandleFunc("/v1/categories", categories.RemoveCategory).Methods("DELETE")

	router.HandleFunc("/v1/cart", cart.CreateNewCart).Methods("POST") 
	router.HandleFunc("/v1/cart", cart.GetAllUserCarts).Methods("GET")
	router.HandleFunc("/v1/cart/{listid}", cart.ViewCart).Methods("GET")
	router.HandleFunc("/v1/cart/{listid}", cart.AddItem).Methods("PATCH")
	router.HandleFunc("/v1/cart/{listid}", cart.DeleteCart).Methods("DELETE")
	router.HandleFunc("/v1/cart/{listid}/{itemid}", cart.RemoveItem).Methods("DELETE") 
	return router
}