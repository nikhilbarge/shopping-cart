package types

import (
	"encoding/json"
	"net/http"
	"net/url"
	"shopping-cart/pkg/database"

	"gopkg.in/mgo.v2/bson"
)

type Item struct {
	ID  bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string          `bson:"name" json:"name,omitempty"`
	Price  int          `bson:"price" json:"price,omitempty"`
}
type CartItem struct { 
	Item 	 Item `bson:"item,omitempty" json:"item,omitempty"`
	Quantity int  `bson:"quantity" json:"quantity"`
	Price  int  `bson:"price" json:"price,omitempty"`
}

type Cart struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Items     []CartItem      `bson:"cartitems" json:"cartitems"` 
}

 
var err error 
// Validate : Validate Cart creation/Updation
func (cart *Cart) Validate(w http.ResponseWriter, r *http.Request) bool {
	errs := url.Values{}
	db := database.Db
	item := Item{}
	if  cart.ID == "" {
		errs.Add("id", "id is required") 
	}
	if  cart.ID != "" {
		oldCart := Cart{}
		err = db.C("cart").Find(bson.M{"_id": cart.ID}).One(&oldCart) 
		if err != nil {
			errs.Add("id", "Invalid Document ID")
		}  
	}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errs.Add("data", "Invalid data") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	} 
	err = db.C("cart").Find(bson.M{"_id": cart.ID}).One(&cart)
	foundInCart:= false;
	for _,crtItm:= range cart.Items {
		if crtItm.Item.ID == item.ID {
			crtItm.Quantity++
			crtItm.Price += item.Price 
			foundInCart =true
			break;
		}
	}
	if !foundInCart{
		newCartItem := CartItem{}
		newCartItem.Item = item
		newCartItem.Price = item.Price
		newCartItem.Quantity = 1 
		cart.Items = append(cart.Items, newCartItem)
	}
 
	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	}
	return true
}

//AddToCart : Update Cart record
func (cart *Cart) AddToCart(w http.ResponseWriter) bool { 
	db := database.Db
	c := db.C("cart") 
	err = c.Update(bson.M{"_id": cart.ID}, bson.M{"$set": cart})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	} 
	return true
} 

//ViewCart : Find Cart records
func (cart *Cart) ViewCart(w http.ResponseWriter) bool { 
	db := database.Db 
	err = db.C("cart").Find(bson.M{"_id": cart.ID}).One(&cart)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": err, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	} 
	return true 
}

// RemoveItem : Find Cart record
func (cart *Cart) RemoveItem(w http.ResponseWriter, id string) bool { 
	errs := url.Values{}
	db := database.Db 
	if bson.IsObjectIdHex(id) {
		err = db.C("cart").Find(bson.M{"_id": cart.ID}).One(&cart)
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
				return false
			} 
			return true
		} 
	} else {
		errs.Add("id", "Invalid Document ID")
	} 
	return true 
}
 
//ClearCart : remove all item from cart
func (cart *Cart) ClearCart(w http.ResponseWriter) bool { 
	db := database.Db
	err = db.C("cart").Find(bson.M{"_id": cart.ID}).One(&cart)
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