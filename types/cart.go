package types

import (
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

 
 