package cart

import (
	"encoding/json"
	"net/http"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"github.com/gorilla/mux"
)

// AddItem : handler function for PATCH /v1/cart call
func AddItem(w http.ResponseWriter, r *http.Request) {
	accessToken := &types.AccessToken{}
	if !accessToken.AuthorizeByToken(w, r) {
		return
	}
	cartid := accessToken.GetUser().CartID
	cart := &types.Cart{}
	cart.ID = cartid
	applog.Info("adding item to cart")
	if cart.Validate(w, r) && cart.AddToCart(w) { 
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
	if !accessToken.AuthorizeByToken(w, r) {
		return
	} 
	
	cartid := accessToken.GetUser().CartID
	cart := &types.Cart{}
	cart.ID = cartid 
	applog.Info("get all items from cart")
	if cart.ViewCart(w) { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": cart.Items, "status": 1}
		json.NewEncoder(w).Encode(response)
	}

}

// RemoveItem : delete item from cart
func RemoveItem(w http.ResponseWriter, r *http.Request) {
	accessToken := &types.AccessToken{}
	if !accessToken.AuthorizeByToken(w, r) {
		return
	}
	params := mux.Vars(r)
	cartid := accessToken.GetUser().CartID
	cart := &types.Cart{}
    cart.ID = cartid 

	if cart.RemoveItem(w, params["itemid"]) { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": cart, "message": "Item Deleted Successfully", "status": 1}
		json.NewEncoder(w).Encode(response)
	} 
}
// ClearCart : delete item from cart
func ClearCart(w http.ResponseWriter, r *http.Request) {
	accessToken := &types.AccessToken{}
	if !accessToken.AuthorizeByToken(w, r) {
		return
	} 
	cartid := accessToken.GetUser().CartID
	cart := &types.Cart{}
    cart.ID = cartid 

	if cart.ClearCart(w) { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": cart, "message": "Deleted Successfully", "status": 1}
		json.NewEncoder(w).Encode(response)
	} 
}