package types

import (
	"net/http"
	time "time"

	"gopkg.in/mgo.v2/bson"
)

//AccessTokenRequestBody : AccessTokenRequestBody structure
type AccessTokenRequest struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

//AccessToken : Authcode structure
type AccessToken struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Token     string        `bson:"token" json:"token"`
	ExpiresAt time.Time     `bson:"expires_at" json:"expires_at"`
	UserID    bson.ObjectId `bson:"user_id" json:"user_id"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

// var user User
// //GetUser : return user object
// func (accessToken *AccessToken) GetUser() User {
// 	user.Password = ""
// 	return user
// }

// func (accessToken *AccessToken) SetUser(usr User)  {
// 	user = usr 
// }

func (accessToken *AccessToken) GetTokenFromRequest(r *http.Request) {
	keys, ok := r.URL.Query()["access_token"] 
	if !ok || len(keys[0]) < 1 {
		accessToken.Token = r.Header.Get("access_token")
	} else {
		accessToken.Token = keys[0]
	}
}