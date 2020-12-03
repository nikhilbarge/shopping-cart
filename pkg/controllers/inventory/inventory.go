package inventory

import (
	"encoding/json"
	"net/http"
	"net/url"
	"shopping-cart/pkg/service"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"github.com/gorilla/mux"
)

// AddItemToInventory : handler function for POST /v1/inventory call
func AddItemToInventory(w http.ResponseWriter, r *http.Request) {
	// authenticating user
	accessToken := &types.AccessToken{}
	accessToken.GetTokenFromRequest(r)
	authService := &service.AuthService{}
	if !authService.AuthorizeByToken(w, accessToken) {
		return
	}
	errs := url.Values{}
	if authService.GetUser().UserName!= "admin" {
		errs.Add("data", "User does not have access to update inventory") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
	}
	item := &types.Item{}
	applog.Info("adding item to item")
	inventoryService := service.InventoryService{}
	if inventoryService.AddToInventory(w, r, item) { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": item, "status": 1}
		json.NewEncoder(w).Encode(response)
	}
	applog.Info("add to inventory request completed")
}
// ViewInventory Get All items in inventory
func ViewInventory(w http.ResponseWriter, r *http.Request) {
	// authenticating user
	accessToken := &types.AccessToken{}
	accessToken.GetTokenFromRequest(r)
	authService := &service.AuthService{}
	if !authService.AuthorizeByToken(w, accessToken) {
		return
	}
	items := &types.ItemList{} 
	applog.Info("get all items from item")
	inventoryService := service.InventoryService{}
	if inventoryService.ViewInvetory(w,items) { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": items, "status": 1}
		json.NewEncoder(w).Encode(response)
	}

}

// RemoveItem : delete item from item
func RemoveItem(w http.ResponseWriter, r *http.Request) {
	// authenticating user
	accessToken := &types.AccessToken{}
	accessToken.GetTokenFromRequest(r)
	authService := &service.AuthService{}
	if !authService.AuthorizeByToken(w, accessToken) {
		return
	}
	errs := url.Values{}
	if authService.GetUser().UserName!= "admin" {
		errs.Add("data", "User does not have access to update inventory") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
	}
	params := mux.Vars(r)
	item := &types.Item{}

	inventoryService := service.InventoryService{}
	if inventoryService.RemoveItem(w,item, params["itemid"]) { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": item, "message": "Item Deleted Successfully", "status": 1}
		json.NewEncoder(w).Encode(response)
	} 
}
 