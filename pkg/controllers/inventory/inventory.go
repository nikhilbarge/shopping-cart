package inventory

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

// AddItemToInventory : handler function for POST /v1/inventory call
func AddItemToInventory(w http.ResponseWriter, r *http.Request) {
	// authenticating user
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	as := service.AuthService{}
	authService := as.NewAuthService()  
	errs := url.Values{}
	if authService.GetUser().UserName!= "admin" {
		errs.Add("data", "Forbiden, user is not 'admin'") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}
	 
	item := &types.Item{} 
	if err := json.NewDecoder(r.Body).Decode(item); err != nil {
		errs.Add("data", "Invalid data") 
		applog.Errorf("invalid request for add item to inventory, %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return 
	} 
	applog.Info("adding item to item")
	is := service.InventoryService{} 
	inventoryService := is.NewInventoryService() 


	err = inventoryService.AddToInventory(item)
	if err!=nil { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": item, "status": 1}
	json.NewEncoder(w).Encode(response)

	applog.Info("add to inventory request completed")
}
// ViewInventory Get All items in inventory
func ViewInventory(w http.ResponseWriter, r *http.Request) {
	// authenticating user  
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	items := &types.ItemList{} 
	applog.Info("get all items from item")
	is := service.InventoryService{}
	inventoryService := is.NewInventoryService()
	err = inventoryService.ViewInvetory(items) 
	if err!=nil{ 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": err, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": items, "status": 1}
	json.NewEncoder(w).Encode(response)
}

// RemoveItem : delete item from item
func RemoveItem(w http.ResponseWriter, r *http.Request) {
	// authenticating user
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	as := service.AuthService{}
	authService := as.NewAuthService()  
	errs := url.Values{}
	if authService.GetUser().UserName!= "admin" {
		errs.Add("data", "User does not have access to update inventory") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}
	params := mux.Vars(r)
	item := &types.Item{}

	is := service.InventoryService{}
	inventoryService := is.NewInventoryService()
	err = inventoryService.RemoveItem(item, params["itemid"])
	if err!=nil { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	} 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": item, "message": "Item Deleted Successfully", "status": 1}
	json.NewEncoder(w).Encode(response)

}
 