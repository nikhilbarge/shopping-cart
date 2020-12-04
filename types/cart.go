package types

import (
	time "time"

	"gopkg.in/mgo.v2/bson"
)

 
type CartItem struct { 
	Item 	 Item `bson:"item,omitempty" json:"item,omitempty"`
	Quantity int  `bson:"quantity" json:"quantity"`
	Price  int    `bson:"price" json:"price,omitempty"`
}

type Cart struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name	  string `bson:"name" json:"name"`
	Items     []CartItem    `bson:"cartitems" json:"cartitems"` 
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	Category  string  `bson:"category" json:"category"`
	UserID	  bson.ObjectId `bson:"userid,omitempty" json:"userid,omitempty"`
}

 
 