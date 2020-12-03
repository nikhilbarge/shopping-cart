package database

import (
	"shopping-cart/types"
	"time"

	"gopkg.in/mgo.v2/bson"
)
 
type DataService struct {
}
//GetAuthToken : fetch token from database
func (authdb *DataService) GetAuthToken(accessToken *types.AccessToken) error {
	now := time.Now().Local()
	return db.C("accesstoken").Find(bson.M{"token": accessToken.Token, "expires_at": bson.M{"$gt": now}}).One(&accessToken)
}

// InsertToken : insert token
func (authdb *DataService) InsertToken(accessToken *types.AccessToken) error {
	return db.C("accesstoken").Insert(&accessToken)
} 
// RemoveToken :
func (authdb *DataService) RemoveToken(accessToken *types.AccessToken) error {
	return  db.C("accesstoken").Remove(&accessToken)
} 

// GetUserByID : Get user from database
func (authdb *DataService) GetUserByID(accessToken *types.AccessToken,  user *types.User) error {
	return db.C("user").Find(bson.M{"_id": accessToken.UserID}).One(&user)
} 
// GetUserByName : Get user from database
func (authdb *DataService) GetUserByName(username string, user *types.User) error {
	return  db.C("user").Find(bson.M{"username": username}).One(&user)
} 

// GetUserByEmail :
func (authdb *DataService) GetUserByEmail(email string, user *types.User) error {
	return  db.C("user").Find(bson.M{"email": email}).One(&user)

}

// InsertUser :
func (authdb *DataService) InsertUser(user *types.User) error {
	return db.C("user").Insert(&user)
} 

// InsertCart :
func (authdb *DataService) InsertCart(cart *types.Cart) error {
	return db.C("cart").Insert(&cart)
} 
// GetCartByID : Get cart from database
func (authdb *DataService) GetCartByID(cartid bson.ObjectId, cart *types.Cart) error {
	return db.C("cart").Find(bson.M{"_id": cartid}).One(&cart)
} 

// UpdateCart : update cart
func (authdb *DataService) UpdateCart(cartid bson.ObjectId, cart *types.Cart) error {
	return db.C("cart").Update(bson.M{"_id": cartid}, bson.M{"$set": cart})
} 

// GetItemByID : Get cart from database
func (authdb *DataService) GetItemByID(itemid bson.ObjectId, item *types.Item) error {
	return db.C("inventory").Find(bson.M{"_id": itemid}).One(&item)
} 

// GetItemByName :
func (authdb *DataService) GetItemByName(itemname string, item *types.Item) error {
	return db.C("inventory").Find(bson.M{"name": itemname}).One(&item)
} 

// GetAllItems :
func (authdb *DataService) GetAllItems() ([]types.Item,error) {
	items := []types.Item{}
	err := db.C("inventory").Find(nil).All(&items)
	return items, err
} 

// UpdateItemByID :
func (authdb *DataService) UpdateItemByID(itemid bson.ObjectId, item *types.Item) error {
	return db.C("inventory").Update(bson.M{"_id": itemid}, bson.M{"$set": item})
} 

// InsertItem :
func (authdb *DataService) InsertItem(item *types.Item) error {
	return db.C("inventory").Insert(&item)
} 

// RemoveItem :
func (authdb *DataService) RemoveItem(itemid bson.ObjectId) error {
	return db.C("inventory").Remove(bson.M{"_id": itemid})
} 

// RemoveAllItem :
func (authdb *DataService) RemoveAllItem() error {
	_, err:=  db.C("inventory").RemoveAll(nil)
	return err
} 
