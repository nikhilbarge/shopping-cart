package types

import (
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
