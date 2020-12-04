package service

import (
	"errors"
	"shopping-cart/pkg/database"
	"shopping-cart/types"
	"shopping-cart/utils/applog"
	"time"

	"github.com/asaskevich/govalidator"
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

//ViewAllCarts : Find Cart records
func (cs *CartService) ViewAllCarts(userid bson.ObjectId) ([]types.Cart, error) {
	applog.Infof("get all list of user %s", userid)
	carts, err := cs.dbsrv.GetAllCartsForUser(userid)
	if err != nil {
		applog.Errorf("error while fetching lists for user %v",userid)
		return nil, err 
	} 
	applog.Debugf("successfully fetched all carts for user %v", userid)
	return carts,nil 
}

// FindUserCart : find user cart
func (cs *CartService) FindUserCart(id string) (*types.Cart,error) {
	applog.Infof("get details for list id %s", id)
	cart := &types.Cart{}
	err := cs.dbsrv.GetCartByID(bson.ObjectIdHex(id), cart)
	if err != nil {
		applog.Errorf("error while fetching details for list %s", id)
		return nil, err 
	} 
	applog.Debugf("successfully fetched details for list %s ", cart.UserID)
	return cart, nil 
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
func (cs *CartService) ClearCart(cartid bson.ObjectId) error { 
	applog.Infof("remove all items in cart %s",cartid)
	oldCart := &types.Cart{}
	err := cs.dbsrv.GetCartByID(cartid, oldCart)
	oldCart.Items = []types.CartItem{}
	err =  cs.dbsrv.UpdateCart(cartid, oldCart)
	if err != nil {
		applog.Errorf("Failed to clear cart, err %v", err )
		return err
	} 
	return nil
}

// DeleteCart :  delete cart
func (cs *CartService) DeleteCart(cartid bson.ObjectId) error { 
	applog.Infof("remove cart %s",cartid) 
	err :=  cs.dbsrv.DeleteCart(cartid)
	if err != nil {
		applog.Errorf("Failed to delete cart, err %v", err )
		return err
	} 
	return nil
}


//CreateCart : create new cart for user
func (cs *CartService) CreateCart(cart *types.Cart) error { 
	if govalidator.IsNull(cart.Name) {
		return errors.New("cart name is required")
	}
	if cart.UserID == "" {
		return errors.New("user id is required")
	}
	
	applog.Infof("check if cart already exists %s",cart.Name)
	err := cs.dbsrv.GetCartByName(cart.Name, cart)
	if err!=nil && err.Error()!="not found" { 
		return err
	}
	if cart.ID != "" {
		return errors.New("cart already exists")
	}
	
	applog.Infof("validate category selected %s",cart.Name)
	category := types.Categories{} 
	
	err = cs.dbsrv.GetCategoriesByName(cart.Category, &category)
	if err!=nil { 
		applog.Errorf("Failed get selected category %v", err )
		return errors.New("selected category not found")
	}

	cart.Items = []types.CartItem{}
	cart.CreatedAt = time.Now().Local()
	applog.Infof("create new cart for user %s",cart.Name)
	err =  cs.dbsrv.InsertCart(cart)
	if err != nil {
		applog.Errorf("Failed to create cart, err %v", err )
		return errors.New("Failed to create cart")
	} 
	return nil
}

//UpdateCart : update cart details
func (cs *CartService) UpdateCart(cart *types.Cart) error { 
	applog.Infof("update cart for user %s",cart.UserID)
	if govalidator.IsNull(cart.Name) {
		return errors.New("cart name is required")
	}
	if cart.UserID == "" {
		return errors.New("invalid user id")
	}
	if cart.ID == "" {
		return errors.New("invalid cart id")
	}
	if cart.Category == "" {
		return errors.New("category is required")
	}
	oldCart:= types.Cart{}
	 
	err:= cs.dbsrv.GetCartByID(cart.ID, &oldCart)
	oldCart.Name = cart.Name  
	oldCart.Category = cart.Category
	// set other cart fields to update if required here
	err = cs.dbsrv.UpdateCart(cart.ID,&oldCart)
	if err != nil {
		applog.Errorf("failed to update cart %s ",cart.Name)
		return err
	}  
	return nil
}