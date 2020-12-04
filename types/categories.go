package types

import "gopkg.in/mgo.v2/bson"

type Categories struct {
	ID  bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name  string          `bson:"name" json:"name,omitempty"`
}

type CategoryList struct {
	List []Categories `bson:"categories" json:"categories,omitempty"`
}