package service

import (
	"encoding/json"
	"net/http"
	"net/url"
	"shopping-cart/pkg/database"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"gopkg.in/mgo.v2/bson"
)

type InventoryService struct{
	DbService database.DataService
}
// AddToInventory : adding/updating Inventory
func (inventory *InventoryService) AddToInventory(w http.ResponseWriter, r *http.Request, item *types.Item) bool {
	applog.Info("validating add to inventory request")
	errs := url.Values{} 
	oldItem:= types.Item{}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errs.Add("data", "Invalid data") 
		applog.Errorf("invalid request for add item to inventory, %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	} 

	applog.Infof("processing request to add item %s to inventory", item.Name)
	err := inventory.DbService.GetItemByName(item.Name,&oldItem)
	if err!=nil && err.Error()!="not found" {
		errs.Add("database", "unable to fetch item details") 
		applog.Errorf("failed to fetch item %s err %v", item.Name,err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	}
	if oldItem.ID != "" {
		applog.Debugf("item %s already exists in inventory, increase in quantity ",item.Name)
		oldItem.Quantity += item.Quantity
		err := inventory.DbService.UpdateItemByID(oldItem.ID,&oldItem)
		if err != nil {
			applog.Errorf("failed to update item %s err %v", item.Name,err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			response := map[string]interface{}{"errors": err.Error(), "status": 0}
			json.NewEncoder(w).Encode(response)
			applog.Errorf("failed to update item %s to invetory ",oldItem.Name)
			return false
		} 
	} else {
		err := inventory.DbService.InsertItem(item)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			response := map[string]interface{}{"errors": err.Error(), "status": 0}
			json.NewEncoder(w).Encode(response)
			applog.Errorf("failed to add item %s to invetory err %v ",oldItem.Name, err)
			return false
		} 
	}
	err = inventory.DbService.GetItemByName(item.Name, item)  
	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		applog.Error("failed to complete add to inventory request")
		return false
	}

	applog.Info("inventory updated successfully")
	return true
}
 

//ViewInvetory : Find All Invetory records
func (inventory *InventoryService) ViewInvetory(w http.ResponseWriter,items *types.ItemList) bool { 
 
	applog.Info("get all items in invetory")
	 
	itemList, err := inventory.DbService.GetAllItems()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": err, "status": 0}
		json.NewEncoder(w).Encode(response)
		applog.Errorf("error while fetching items in inventory")
		return  false
	} 
	items.Items = itemList
	applog.Debugf("successfully fetched items in invetory")
	return true
}

// RemoveItem : Remove Invetory record
func (inventory *InventoryService) RemoveItem(w http.ResponseWriter,item *types.Item, id string) bool { 
	errs := url.Values{} 
	applog.Infof("remove all items in inventory")
	if id != "" && bson.IsObjectIdHex(id) {
		err := inventory.DbService.RemoveItem(item.ID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			response := map[string]interface{}{"errors": err.Error(), "status": 0}
			json.NewEncoder(w).Encode(response)
			applog.Errorf("error while removing items from inventory %s",item.ID)
			return false
		} 
		applog.Infof("removed all item %s from inventory",item.ID)
		return true
	} else if id == "" {
		err := inventory.DbService.RemoveAllItem()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			response := map[string]interface{}{"errors": err.Error(), "status": 0}
			json.NewEncoder(w).Encode(response)
			applog.Errorf("error while removing items from inventory")
			return false
		} 
		applog.Infof("removed all items from inventory")
		return true
	} else {
		errs.Add("id", "Invalid Document ID")
		applog.Debugf("invalid item id %s",item.ID)
		return false
	}   
}
 