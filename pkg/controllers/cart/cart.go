package cart

import (
	"encoding/json"
	"net/http"
	"net/url"
	"shopping-cart/pkg/controllers/common"
	"shopping-cart/pkg/service"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"github.com/gorilla/mux"
)

// AddItem : handler function for PATCH /v1/cart call
func AddItem(w http.ResponseWriter, r *http.Request) {
	// authenticating user
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	as := service.AuthService{}
	authService := as.NewAuthService()  

	// proccess add to cart request
	cartid := authService.GetUser().CartID
	cart := &types.Cart{}
	cart.ID = cartid
	errs := url.Values{} 
	if  cart.ID == "" {
		applog.Debug("unable to find cart")
		errs.Add("id", "id is required") 
	}

	reqItem := types.Item{}
	if err := json.NewDecoder(r.Body).Decode(&reqItem); err != nil {
		errs.Add("data", "Invalid data") 
		applog.Errorf("invalid request for cart %s", cart.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return 
	} 
	if reqItem.ID == "" {
		applog.Debug("unable to find item")
		errs.Add("id", "item id is required") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return 
	} 


	applog.Info("adding item to cart") 
	crt := service.CartService{}
	cartService := crt.NewCartService()
	err = cartService.Validate(&reqItem, cart)
	if err!=nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}

	 err = cartService.AddToCart(cart)
	 if err !=nil { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": cart, "status": 1}
	json.NewEncoder(w).Encode(response)
	applog.Info("add to cart request completed")
}
// ViewCart Get All items in cart
func ViewCart(w http.ResponseWriter, r *http.Request) {
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	as := service.AuthService{}
	authService := as.NewAuthService()   
	cartid := authService.GetUser().CartID
	cart := &types.Cart{}
	cart.ID = cartid 
	applog.Info("get all items from cart")
	crt := service.CartService{}
	cartService := crt.NewCartService()
	err = cartService.ViewCart(cart)
	if err!=nil {  
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": err, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": cart.Items, "status": 1}
	json.NewEncoder(w).Encode(response)
}

// RemoveItem : delete item from cart
func RemoveItem(w http.ResponseWriter, r *http.Request) {
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	as := service.AuthService{}
	authService := as.NewAuthService() 
	params := mux.Vars(r)
	cartid := authService.GetUser().CartID
	cart := &types.Cart{}
    cart.ID = cartid 
	crt := service.CartService{}
	cartService := crt.NewCartService()
	err = cartService.RemoveItem(cart, params["itemid"])
	if err!=nil { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
	} 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": cart, "message": "Item Deleted Successfully", "status": 1}
	json.NewEncoder(w).Encode(response)
}
// ClearCart : delete item from cart
func ClearCart(w http.ResponseWriter, r *http.Request) {
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	as := service.AuthService{}
	authService := as.NewAuthService() 
	 
	cartid := authService.GetUser().CartID
	cart := &types.Cart{}
    cart.ID = cartid 
	crt := service.CartService{}
	cartService := crt.NewCartService()
	err= cartService.ClearCart(cart) 
	if err!= nil{ 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
	} 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": cart, "message": "Deleted Successfully", "status": 1}
	json.NewEncoder(w).Encode(response)
}