package database

import (
	"shopping-cart/types"
	"time"

	"gopkg.in/mgo.v2/bson"
)
func NewDBService() *dataService {
	return &dataService{}
}
type IDataService interface {
	GetAuthToken(accessToken *types.AccessToken) error
	InsertToken(accessToken *types.AccessToken) error
	RemoveToken(accessToken *types.AccessToken) error
	
	GetUserByID(accessToken *types.AccessToken,  user *types.User) error
	GetUserByName(username string, user *types.User) error
	GetUserByEmail(email string, user *types.User) error
	
	InsertUser(user *types.User) error
	InsertCart(cart *types.Cart) error
	GetCartByID(cartid bson.ObjectId, cart *types.Cart) error
	GetAllCartsForUser(userid bson.ObjectId) ([]types.Cart, error)
	GetCartByName(cartname string, cart *types.Cart) error
	UpdateCart(cartid bson.ObjectId, cart *types.Cart) error
	DeleteCart(cartid bson.ObjectId) error

	GetItemByID(itemid bson.ObjectId, item *types.Item) error
	GetItemByName(itemname string, item *types.Item) error
 	GetAllItems() ([]types.Item,error)
	UpdateItemByID(itemid bson.ObjectId, item *types.Item) error 
	InsertItem(item *types.Item) error
	RemoveItem(itemid bson.ObjectId) error
	RemoveAllItem() error 

	GetCategoriesByName(categoryname string, category *types.Categories) error
	GetCategoriesByID(id bson.ObjectId, category *types.Categories) error
 	GetAllCategoriess() ([]types.Categories,error)
	UpdateCategoriesByID(categoryid bson.ObjectId, category *types.Categories) error 
	InsertCategories(category *types.Categories) error
	RemoveCategories(categoryid bson.ObjectId) error
	RemoveAllCategories() error 
}
type dataService struct {
	IDataService
}
//GetAuthToken : fetch token from database
func (ds *dataService) GetAuthToken(accessToken *types.AccessToken) error {
	now := time.Now().Local()
	return db.C("accesstoken").Find(bson.M{"token": accessToken.Token, "expires_at": bson.M{"$gt": now}}).One(&accessToken)
}

// InsertToken : insert token
func (ds *dataService) InsertToken(accessToken *types.AccessToken) error {
	return db.C("accesstoken").Insert(&accessToken)
} 
// RemoveToken :
func (ds *dataService) RemoveToken(accessToken *types.AccessToken) error {
	return  db.C("accesstoken").Remove(&accessToken)
} 

// GetUserByID : Get user from database
func (ds *dataService) GetUserByID(accessToken *types.AccessToken,  user *types.User) error {
	return db.C("user").Find(bson.M{"_id": accessToken.UserID}).One(&user)
} 
// GetUserByName : Get user from database
func (ds *dataService) GetUserByName(username string, user *types.User) error {
	return  db.C("user").Find(bson.M{"username": username}).One(&user)
} 

// GetUserByEmail :
func (ds *dataService) GetUserByEmail(email string, user *types.User) error {
	return  db.C("user").Find(bson.M{"email": email}).One(&user)

}

// InsertUser :
func (ds *dataService) InsertUser(user *types.User) error {
	return db.C("user").Insert(&user)
} 

// InsertCart :
func (ds *dataService) InsertCart(cart *types.Cart) error {
	return db.C("cart").Insert(&cart)
} 
// GetCartByID : Get cart from database
func (ds *dataService) GetCartByID(cartid bson.ObjectId, cart *types.Cart) error {
	return db.C("cart").Find(bson.M{"_id": cartid}).One(&cart)
} 

// GetAllCartsForUser : Get cart from database
func (ds *dataService) GetAllCartsForUser(userid bson.ObjectId) ([]types.Cart, error) {
	carts := []types.Cart{}
	err := db.C("cart").Find(bson.M{"userid": userid}).All(&carts)
	return carts, err
} 

// GetCartByName : Get cart from database by name
func (ds *dataService) GetCartByName(name string, cart *types.Cart) error {
	return db.C("cart").Find(bson.M{"name": name}).One(&cart)
} 
// UpdateCart : update cart
func (ds *dataService) UpdateCart(cartid bson.ObjectId, cart *types.Cart) error {
	return db.C("cart").Update(bson.M{"_id": cartid}, bson.M{"$set": cart})
} 

// RemoveItem :
func (ds *dataService) DeleteCart(cartid bson.ObjectId) error {
	return db.C("cart").Remove(bson.M{"_id": cartid})
} 

// GetItemByID : Get cart from database
func (ds *dataService) GetItemByID(itemid bson.ObjectId, item *types.Item) error {
	return db.C("inventory").Find(bson.M{"_id": itemid}).One(&item)
} 

// GetItemByName :
func (ds *dataService) GetItemByName(itemname string, item *types.Item) error {
	return db.C("inventory").Find(bson.M{"name": itemname}).One(&item)
} 

// GetAllItems :
func (ds *dataService) GetAllItems() ([]types.Item,error) {
	items := []types.Item{}
	err := db.C("inventory").Find(nil).All(&items)
	return items, err
} 

// UpdateItemByID :
func (ds *dataService) UpdateItemByID(itemid bson.ObjectId, item *types.Item) error {
	return db.C("inventory").Update(bson.M{"_id": itemid}, bson.M{"$set": item})
} 

// InsertItem :
func (ds *dataService) InsertItem(item *types.Item) error {
	return db.C("inventory").Insert(&item)
} 

// RemoveItem :
func (ds *dataService) RemoveItem(itemid bson.ObjectId) error {
	return db.C("inventory").Remove(bson.M{"_id": itemid})
} 

// RemoveAllItem :
func (ds *dataService) RemoveAllItem() error {
	_, err:=  db.C("inventory").RemoveAll(nil)
	return err
} 






// GetCategoriesByID : Get cart from database
func (ds *dataService) GetCategoriesByID(categoryid bson.ObjectId, category *types.Categories) error {
	return db.C("categories").Find(bson.M{"_id": categoryid}).One(&category)
} 

// GetCategoriesByName :
func (ds *dataService) GetCategoriesByName(categoryname string, category *types.Categories) error {
	return db.C("categories").Find(bson.M{"name": categoryname}).One(&category)
} 

// GetAllCategoriess :
func (ds *dataService) GetAllCategoriess() ([]types.Categories,error) {
	categorys := []types.Categories{}
	err := db.C("categories").Find(nil).All(&categorys)
	return categorys, err
} 

// UpdateCategoriesByID :
func (ds *dataService) UpdateCategoriesByID(categoryid bson.ObjectId, category *types.Categories) error {
	return db.C("categories").Update(bson.M{"_id": categoryid}, bson.M{"$set": category})
} 

// InsertCategories :
func (ds *dataService) InsertCategories(category *types.Categories) error {
	return db.C("categories").Insert(&category)
} 

// RemoveCategories :
func (ds *dataService) RemoveCategories(categoryid bson.ObjectId) error {
	return db.C("categories").Remove(bson.M{"_id": categoryid})
} 

// RemoveAllCategories :
func (ds *dataService) RemoveAllCategories() error {
	_, err:=  db.C("categories").RemoveAll(nil)
	return err
} 
