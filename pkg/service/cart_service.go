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

type CartService struct {
	DbService database.DataService
}
// Validate : Validate Cart creation/Updation
func (cartservice *CartService) Validate(w http.ResponseWriter, r *http.Request, cart *types.Cart) bool {
	applog.Info("validating add to cart request")
	errs := url.Values{}
	 
	
	if  cart.ID == "" {
		applog.Debug("unable to find cart")
		errs.Add("id", "id is required") 
	}
	if  cart.ID != "" {
		oldCart := types.Cart{}
		err :=  cartservice.DbService.GetCartByID(cart.ID, &oldCart)
		if err != nil {
			applog.Errorf("invalid cart selected %s", cart.ID)
			errs.Add("id", "Invalid Document ID")
		}  
	}
	reqItem := types.Item{}
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
	inventoryItem := types.Item{}
	if reqItem.ID!=""{
		err := cartservice.DbService.GetItemByID(reqItem.ID,&inventoryItem) 
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
	err := cartservice.DbService.GetCartByID(cart.ID, cart)
	if err!=nil {
		errs.Add("database", "unable to fetch cart details") 
		applog.Errorf("invalid request for cart %s", cart.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	}
	newCartItem := types.CartItem{}
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
func (cartservice *CartService) AddToCart(w http.ResponseWriter, cart *types.Cart) bool { 

 	applog.Infof("updating cart %s with new items",cart.ID)
	err := cartservice.DbService.UpdateCart(cart.ID,cart)
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
func (cartservice *CartService) ViewCart(w http.ResponseWriter, cart *types.Cart) bool { 
 
	applog.Infof("get all items in cart %v", cart.ID)
	err := cartservice.DbService.GetCartByID(cart.ID, cart)
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
func (cartservice *CartService) RemoveItem(w http.ResponseWriter,cart *types.Cart, id string) bool { 
	errs := url.Values{}
 
	applog.Infof("remove items in cart %s",cart.ID)
	if bson.IsObjectIdHex(id) {
		err := cartservice.DbService.GetCartByID(cart.ID, cart)
		idx:= -1; 
		for i ,crtItm := range cart.Items {
			if crtItm.Item.ID == bson.ObjectIdHex(id) {
				idx = i 
				break
			}
		} 		
		if idx != -1 {
			cart.Items = append(cart.Items[:idx], cart.Items[idx+1:]...)
			err = cartservice.DbService.UpdateCart(cart.ID,cart)
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
func (cartservice *CartService) ClearCart(w http.ResponseWriter, cart *types.Cart) bool { 
	applog.Infof("remove all items in cart %s",cart.ID)
	err := cartservice.DbService.GetCartByID(cart.ID, cart)
	cart.Items = []types.CartItem{}
	err =  cartservice.DbService.UpdateCart(cart.ID, cart)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	} 
	return true
}