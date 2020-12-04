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
	"gopkg.in/mgo.v2/bson"
)

// AddItem : handler function for PATCH /v1/cart call
func AddItem(w http.ResponseWriter, r *http.Request) {
	// authenticating user
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	} 
	errs := url.Values{} 
	// proccess add to cart request
	params := mux.Vars(r) 
	cartid := params["listid"]
	if  cartid == "" && !bson.IsObjectIdHex(cartid){ 
		errs.Add("listid", "list id is required") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	} 
	reqItem := types.Item{}
	if err := json.NewDecoder(r.Body).Decode(&reqItem); err != nil {
		errs.Add("data", "Invalid data") 
		applog.Errorf("invalid add to cart request for list %s", cartid)
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
	
	cart, err := cartService.FindUserCart(cartid)
	if err!=nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}

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
		return
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
	param := mux.Vars(r)
	cartid := param["listid"]    
	applog.Info("get all items in list %s ", cartid)
	crt := service.CartService{}
	cartService := crt.NewCartService()
	cart, err := cartService.FindUserCart(cartid)
	if err!=nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
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
	errs := url.Values{} 
	params := mux.Vars(r)  
	
	
	cartid := params["listid"]
	if  cartid == "" && !bson.IsObjectIdHex(cartid){ 
		errs.Add("listid", "list id is required") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	} 
	itemid := params["itemid"]
	if itemid == "" && !bson.IsObjectIdHex(itemid){
		errs.Add("itemid", "item id is required") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	} 
	crt := service.CartService{}
	cartService := crt.NewCartService()
	cart, err := cartService.FindUserCart(cartid)
	if err!=nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}
	err = cartService.RemoveItem(cart, itemid)
	if err!=nil { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	} 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": cart, "message": "Item Deleted Successfully", "status": 1}
	json.NewEncoder(w).Encode(response)
}
// DeleteCart : delete item from cart
func DeleteCart(w http.ResponseWriter, r *http.Request) {
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	errs := url.Values{} 
	params := mux.Vars(r)
	cartid := params["listid"]
	if  cartid == "" && !bson.IsObjectIdHex(cartid){ 
		errs.Add("listid", "list id is required") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	} 
	crt := service.CartService{}
	cartService := crt.NewCartService() 
	err= cartService.DeleteCart(bson.ObjectIdHex(cartid)) 
	if err!= nil{ 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	} 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": cartid, "message": "Deleted Successfully", "status": 1}
	json.NewEncoder(w).Encode(response)
}


// CreateNewCart : create new empty cart for user
func CreateNewCart(w http.ResponseWriter, r *http.Request) {
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	errs := url.Values{} 
	cart := &types.Cart{}
	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		errs.Add("data", "Invalid data") 
		applog.Error("invalid request for create cart")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return 
	} 
	as := service.AuthService{}
	authService := as.NewAuthService()  
	cart.UserID = authService.GetUser().ID
	crt := service.CartService{}
	cartService := crt.NewCartService()
	err = cartService.CreateCart(cart) 
	if err!= nil{ 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		return 
	} 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": cart, "message": "Cart Create Successfully", "status": 1}
	json.NewEncoder(w).Encode(response)
}


// UpdateUserCart : update user cart
func UpdateUserCart(w http.ResponseWriter, r *http.Request) {
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	errs := url.Values{} 
	params := mux.Vars(r) 
	cartid := params["listid"]
	if  cartid == "" && !bson.IsObjectIdHex(cartid){ 
		errs.Add("listid", "list id is required") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	} 
	cart := &types.Cart{} 
	if err := json.NewDecoder(r.Body).Decode(&cart); err != nil {
		errs.Add("data", "Invalid cart details") 
		applog.Error("invalid request for update cart")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return 
	}  
	cart.ID = bson.ObjectIdHex(cartid)
	
	as := service.AuthService{}
	authService := as.NewAuthService() 
	cart.UserID = authService.GetUser().ID
	
	crt := service.CartService{}
	cartService := crt.NewCartService()
	err = cartService.UpdateCart(cart) 
	if err!= nil{ 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
	} 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": cart, "message": "Cart Updated Successfully", "status": 1}
	json.NewEncoder(w).Encode(response)
}

// GetAllUserCarts : get all list
func GetAllUserCarts(w http.ResponseWriter, r *http.Request) {
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	} 
	crt := service.CartService{}
	cartService := crt.NewCartService()
	as := service.AuthService{}
	authService := as.NewAuthService() 
	userid := authService.GetUser().ID
	carts ,err := cartService.ViewAllCarts(userid)
	if err!=nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	} 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": carts, "status": 1}
	json.NewEncoder(w).Encode(response)
}