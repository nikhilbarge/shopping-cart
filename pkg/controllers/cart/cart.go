package cart

import (
	"encoding/json"
	"net/http"
	"shopping-cart/pkg/service"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"github.com/gorilla/mux"
)

// AddItem : handler function for PATCH /v1/cart call
func AddItem(w http.ResponseWriter, r *http.Request) {
	// authenticating user
	accessToken := &types.AccessToken{}
	accessToken.GetTokenFromRequest(r)
	authService := &service.AuthService{}
	if !authService.AuthorizeByToken(w, accessToken) {
		return
	}
	// proccess add to cart request
	cartid := authService.GetUser().CartID
	cart := &types.Cart{}
	cart.ID = cartid
	applog.Info("adding item to cart")
	cartService:=service.CartService{}
	if cartService.Validate(w, r, cart) && cartService.AddToCart(w, cart) { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": cart, "status": 1}
		json.NewEncoder(w).Encode(response)
	}
	applog.Info("add to cart request completed")
}
// ViewCart Get All items in cart
func ViewCart(w http.ResponseWriter, r *http.Request) {
	accessToken := &types.AccessToken{}
	accessToken.GetTokenFromRequest(r)
	authService := &service.AuthService{}
	if !authService.AuthorizeByToken(w, accessToken) {
		return
	}
	
	cartid := authService.GetUser().CartID
	cart := &types.Cart{}
	cart.ID = cartid 
	applog.Info("get all items from cart")
	cartService:=service.CartService{}
	if cartService.ViewCart(w, cart) { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": cart.Items, "status": 1}
		json.NewEncoder(w).Encode(response)
	}

}

// RemoveItem : delete item from cart
func RemoveItem(w http.ResponseWriter, r *http.Request) {
	accessToken := &types.AccessToken{}
	accessToken.GetTokenFromRequest(r)
	authService := &service.AuthService{}
	if !authService.AuthorizeByToken(w, accessToken) {
		return
	}
	params := mux.Vars(r)
	cartid := authService.GetUser().CartID
	cart := &types.Cart{}
    cart.ID = cartid 
	cartService:=service.CartService{}
	if cartService.RemoveItem(w, cart, params["itemid"]) { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": cart, "message": "Item Deleted Successfully", "status": 1}
		json.NewEncoder(w).Encode(response)
	} 
}
// ClearCart : delete item from cart
func ClearCart(w http.ResponseWriter, r *http.Request) {
	accessToken := &types.AccessToken{}
	accessToken.GetTokenFromRequest(r)
	authService := &service.AuthService{}
	if !authService.AuthorizeByToken(w, accessToken) {
		return
	}
	cartid := authService.GetUser().CartID
	cart := &types.Cart{}
    cart.ID = cartid 
	cartService:=service.CartService{}
	if cartService.ClearCart(w,cart) { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": cart, "message": "Deleted Successfully", "status": 1}
		json.NewEncoder(w).Encode(response)
	} 
}