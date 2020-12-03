package types

import (
	"encoding/json"
	"net/http"
	"net/url"
	"shopping-cart/pkg/database"
	"shopping-cart/utils/applog"

	"gopkg.in/mgo.v2/bson"
)

 
type CartItem struct { 
	Item 	 Item `bson:"item,omitempty" json:"item,omitempty"`
	Quantity int  `bson:"quantity" json:"quantity"`
	Price  int    `bson:"price" json:"price,omitempty"`
}

type Cart struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Items     []CartItem    `bson:"cartitems" json:"cartitems"` 
}

 
 
// Validate : Validate Cart creation/Updation
func (cart *Cart) Validate(w http.ResponseWriter, r *http.Request) bool {
	applog.Info("validating add to cart request")
	errs := url.Values{}
	db := database.Db
	
	if  cart.ID == "" {
		applog.Debug("unable to find cart")
		errs.Add("id", "id is required") 
	}
	if  cart.ID != "" {
		oldCart := Cart{}
		err := db.C("cart").Find(bson.M{"_id": cart.ID}).One(&oldCart) 
		if err != nil {
			applog.Errorf("invalid cart selected %s", cart.ID)
			errs.Add("id", "Invalid Document ID")
		}  
	}
	reqItem := Item{}
	if err := json.NewDecoder(r.Body).Decode(&reqItem); err != nil {
		errs.Add("data", "Invalid data") 
		applog.Errorf("invalid request for cart %s", cart.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	} 
	if reqItem.ID == "" {
		applog.Debug("unable to find item")
		errs.Add("id", "item id is required") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	} 

	// get item from inventory
	inventoryItem := Item{}
	if reqItem.ID!=""{
		err := db.C("inventory").Find(bson.M{"_id": reqItem.ID}).One(&inventoryItem) 
		if err != nil {
			applog.Errorf("invalid item selected %s", reqItem.ID)
			errs.Add("id", "Invalid item ID") 
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{"errors": errs, "status": 0}
			json.NewEncoder(w).Encode(response)
			return false
		}  
	}
	
	applog.Infof("processing request to add item %s to cart %s", inventoryItem.Name, cart.ID)
	err := db.C("cart").Find(bson.M{"_id": cart.ID}).One(&cart)
	if err!=nil {
		errs.Add("database", "unable to fetch cart details") 
		applog.Errorf("invalid request for cart %s", cart.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	}
	newCartItem := CartItem{}
	newCartItem.Item = inventoryItem
	newCartItem.Item.Quantity=0
	newCartItem.Price = inventoryItem.Price
	newCartItem.Quantity = 1 
	idx:= -1;
	for i ,crtItm:= range cart.Items {
		if crtItm.Item.ID == inventoryItem.ID {
			idx = i
			applog.Debugf("item %s already exists in cart %s, increase in quantity ",inventoryItem.Name, cart.ID)
			break;
		}
	}
	if idx==-1{ 
		cart.Items = append(cart.Items, newCartItem)
		applog.Debugf("item %s will be added to cart %s",inventoryItem.Name, cart.ID)
	} else {
		cart.Items[idx].Quantity ++
		cart.Items[idx].Price += inventoryItem.Price 
	}
 
	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		applog.Error("failed to complete add to cart request")
		return false
	}
	applog.Info("cart is updated successfully")
	return true
}

//AddToCart : Update Cart record
func (cart *Cart) AddToCart(w http.ResponseWriter) bool { 
	db := database.Db
	c := db.C("cart") 
	applog.Infof("updating cart %s with new items",cart.ID)
	err := c.Update(bson.M{"_id": cart.ID}, bson.M{"$set": cart})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		applog.Errorf("failed to add items to cart %s ",cart.ID)
		return false
	} 
	applog.Infof("updated cart %s with new items",cart.ID)
	return true
} 

//ViewCart : Find Cart records
func (cart *Cart) ViewCart(w http.ResponseWriter) bool { 
	db := database.Db 
	applog.Infof("get all items in cart %v", cart.ID)
	err := db.C("cart").Find(bson.M{"_id": cart.ID}).One(&cart)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": err, "status": 0}
		json.NewEncoder(w).Encode(response)
		applog.Errorf("error while fetching items in cart %v",cart.ID)
		return false 
	} 
	applog.Debugf("successfully fetched items in cart %v", cart.ID)
	return true 
}

// RemoveItem : Find Cart record
func (cart *Cart) RemoveItem(w http.ResponseWriter, id string) bool { 
	errs := url.Values{}
	db := database.Db 
	applog.Infof("remove items in cart %s",cart.ID)
	if bson.IsObjectIdHex(id) {
		err := db.C("cart").Find(bson.M{"_id": cart.ID}).One(&cart)
		idx:= -1; 
		for i ,crtItm := range cart.Items {
			if crtItm.Item.ID == bson.ObjectIdHex(id) {
				idx = i 
				break
			}
		} 		
		if idx != -1 {
			cart.Items = append(cart.Items[:idx], cart.Items[idx+1:]...)
			err = db.C("cart").Update(bson.M{"_id": cart.ID}, bson.M{"$set": cart})
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				response := map[string]interface{}{"errors": err.Error(), "status": 0}
				json.NewEncoder(w).Encode(response)
				applog.Errorf("error while removing items from cart %s",cart.ID.Hex())
				return false
			} 
			return true
		} 
	} else {
		errs.Add("id", "Invalid Document ID")
		applog.Debugf("invalid cart id %s",cart.ID)
		return false
	} 
	applog.Infof("removed all items from cart %s",cart.ID)
	return true 
}
 
//ClearCart : remove all item from cart
func (cart *Cart) ClearCart(w http.ResponseWriter) bool { 
	db := database.Db
	applog.Infof("remove all items in cart %s",cart.ID)
	err := db.C("cart").Find(bson.M{"_id": cart.ID}).One(&cart)
	cart.Items = []CartItem{}
	err = db.C("cart").Update(bson.M{"_id": cart.ID}, bson.M{"$set": cart})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	} 
	return true
}