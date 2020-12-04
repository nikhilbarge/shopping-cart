package service

import (
	"errors"
	"shopping-cart/pkg/database"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"gopkg.in/mgo.v2/bson"
)

type CartService struct {
	dbsrv database.IDataService
} 
func (cs *CartService) NewCartService() *CartService {
	dbservice := database.NewDBService()
	return &CartService{dbsrv: dbservice}
}
// Validate : Validate Cart creation/Updation
func (cs *CartService) Validate(reqItem *types.Item, cart *types.Cart) error {
	applog.Info("validating add to cart request")
	if  cart.ID != "" {
		oldCart := types.Cart{}
		err :=  cs.dbsrv.GetCartByID(cart.ID, &oldCart)
		if err != nil {
			applog.Errorf("invalid cart selected %s", cart.ID)
			return err
		}  
	}
	
	// get item from inventory
	inventoryItem := types.Item{}
	if reqItem.ID!=""{
		err := cs.dbsrv.GetItemByID(reqItem.ID,&inventoryItem) 
		if err != nil {
			applog.Errorf("invalid item selected %s", reqItem.ID) 
			return err
		}  
	}
	
	applog.Infof("processing request to add item %s to cart %s", inventoryItem.Name, cart.ID)
	err := cs.dbsrv.GetCartByID(cart.ID, cart)
	if err!=nil {  
		return err
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
	applog.Info("cart is updated successfully")
	return nil
}

//AddToCart : Update Cart record
func (cs *CartService) AddToCart(cart *types.Cart) error { 
 	applog.Infof("updating cart %s with new items",cart.ID)
	err := cs.dbsrv.UpdateCart(cart.ID,cart)
	if err != nil {
		applog.Errorf("failed to add items to cart %s ",cart.ID)
		return err
	} 
	applog.Infof("updated cart %s with new items",cart.ID)
	return nil
} 

//ViewCart : Find Cart records
func (cs *CartService) ViewCart(cart *types.Cart) error { 
 
	applog.Infof("get all items in cart %v", cart.ID)
	err := cs.dbsrv.GetCartByID(cart.ID, cart)
	if err != nil {
		applog.Errorf("error while fetching items in cart %v",cart.ID)
		return err 
	} 
	applog.Debugf("successfully fetched items in cart %v", cart.ID)
	return nil 
}

// RemoveItem : Find Cart record
func (cs *CartService) RemoveItem(cart *types.Cart, id string) error {  
	applog.Infof("remove items in cart %s",cart.ID)
	if bson.IsObjectIdHex(id) {
		err := cs.dbsrv.GetCartByID(cart.ID, cart)
		idx:= -1; 
		for i ,crtItm := range cart.Items {
			if crtItm.Item.ID == bson.ObjectIdHex(id) {
				idx = i 
				break
			}
		} 		
		if idx != -1 {
			cart.Items = append(cart.Items[:idx], cart.Items[idx+1:]...)
			err = cs.dbsrv.UpdateCart(cart.ID,cart)
			if err != nil { 
				applog.Errorf("error while removing items from cart %s",cart.ID.Hex())
				return err
			} 
			return nil
		} 
	} else {
		applog.Debugf("invalid item id %s", id)
		return errors.New("invalid item id")
	} 
	applog.Infof("removed all items from cart %s",cart.ID)
	return nil 
}
 
//ClearCart : remove all item from cart
func (cs *CartService) ClearCart(cart *types.Cart) error { 
	applog.Infof("remove all items in cart %s",cart.ID)
	err := cs.dbsrv.GetCartByID(cart.ID, cart)
	cart.Items = []types.CartItem{}
	err =  cs.dbsrv.UpdateCart(cart.ID, cart)
	if err != nil {
		applog.Errorf("Failed to clear cart, err %v", err )
		return err
	} 
	return nil
}