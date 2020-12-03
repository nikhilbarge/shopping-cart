package inventory

import (
	"encoding/json"
	"net/http"
	"net/url"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"github.com/gorilla/mux"
)

// AddItemToInventory : handler function for POST /v1/inventory call
func AddItemToInventory(w http.ResponseWriter, r *http.Request) {
	accessToken := &types.AccessToken{}
	if !accessToken.AuthorizeByToken(w, r) {
		return
	} 
	errs := url.Values{}
	if accessToken.GetUser().UserName!= "admin" {
		errs.Add("data", "User does not have access to update inventory") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
	}
	item := &types.Item{}
	applog.Info("adding item to item")
	if item.AddToInventory(w, r) { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": item, "status": 1}
		json.NewEncoder(w).Encode(response)
	}
	applog.Info("add to inventory request completed")
}
// ViewInventory Get All items in inventory
func ViewInventory(w http.ResponseWriter, r *http.Request) {
	accessToken := &types.AccessToken{}
	if !accessToken.AuthorizeByToken(w, r) {
		return
	} 
	items := &types.ItemList{} 
	applog.Info("get all items from item")
	if items.ViewInvetory(w) { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": items, "status": 1}
		json.NewEncoder(w).Encode(response)
	}

}

// RemoveItem : delete item from item
func RemoveItem(w http.ResponseWriter, r *http.Request) {
	accessToken := &types.AccessToken{}
	if !accessToken.AuthorizeByToken(w, r) {
		return
	} 
	errs := url.Values{}
	if accessToken.GetUser().UserName!= "admin" {
		errs.Add("data", "User does not have access to update inventory") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
	}
	params := mux.Vars(r)
	item := &types.Item{}


	if item.RemoveItem(w, params["itemid"]) { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": item, "message": "Item Deleted Successfully", "status": 1}
		json.NewEncoder(w).Encode(response)
	} 
}
 