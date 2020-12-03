package types

import (
	"encoding/json"
	"net/http"
	"net/url"
	"shopping-cart/pkg/database"
	"shopping-cart/utils/applog"

	"gopkg.in/mgo.v2/bson"
)

type Item struct {
	ID  bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string          `bson:"name" json:"name,omitempty"`
	Price  int          `bson:"price" json:"price,omitempty"`
	Quantity int  `bson:"quantity" json:"quantity,omitempty"`
}  

type ItemList struct {
	Items []Item `bson:"items" json:"items,omitempty"`
}
// AddToInventory : adding/updating Inventory
func (item *Item) AddToInventory(w http.ResponseWriter, r *http.Request) bool {
	applog.Info("validating add to inventory request")
	errs := url.Values{}
	db := database.Db 
	oldItem:= Item{}
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
	err := db.C("inventory").Find(bson.M{"name": item.Name}).One(&oldItem)
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
		err := db.C("inventory").Update(bson.M{"_id": oldItem.ID}, bson.M{"$set": oldItem})
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
		err := db.C("inventory").Insert(&item)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			response := map[string]interface{}{"errors": err.Error(), "status": 0}
			json.NewEncoder(w).Encode(response)
			applog.Errorf("failed to add item %s to invetory err %v ",oldItem.Name, err)
			return false
		} 
	}
	err = db.C("inventory").Find(bson.M{"name": item.Name}).One(&item)	 
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
func (items *ItemList) ViewInvetory(w http.ResponseWriter) bool { 
	db := database.Db 
	applog.Info("get all items in invetory")
	itemList := []Item{}
	err := db.C("inventory").Find(nil).All(&itemList)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": err, "status": 0}
		json.NewEncoder(w).Encode(response)
		applog.Errorf("error while fetching items in inventory")
		return false 
	} 
	items.Items = itemList
	applog.Debugf("successfully fetched items in invetory")
	return true 
}

// RemoveItem : Remove Invetory record
func (item *Item) RemoveItem(w http.ResponseWriter, id string) bool { 
	errs := url.Values{}
	db := database.Db 
	applog.Infof("remove all items in inventory")
	if id != "" && bson.IsObjectIdHex(id) {
		err := db.C("inventory").Remove(bson.M{"_id": item.ID})
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
		_,err := db.C("inventory").RemoveAll(nil)
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
 