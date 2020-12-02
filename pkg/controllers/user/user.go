package user

import (
	"encoding/json"
	"net/http"
	"shopping-cart/pkg/database"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"gopkg.in/mgo.v2/bson"

	"github.com/jameskeane/bcrypt"
)

// RegisterUser : Register user account
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	user := &types.User{}
	if !user.Validate(w, r) {
		return
	}
	db := database.Db
	userDb := db.C("user")
	cart := &types.Cart{}
	cartDb:= db.C("cart")
	// Insert
	cart.ID = bson.NewObjectId()
	cart.Items= []types.CartItem{}
	cartErr := cartDb.Insert(&cart)
	if cartErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": cartErr.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)

	}
	user.ID = bson.NewObjectId()
	user.CartID = cart.ID
	salt, _ := bcrypt.Salt(10)
	user.Password, _ = bcrypt.Hash(user.Password, salt)

	insertionErrors := userDb.Insert(&user)

	if insertionErrors != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": insertionErrors.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)

	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": user, "status": 1}
		json.NewEncoder(w).Encode(response)
	}

}

// Login : login to account
func Login(w http.ResponseWriter, r *http.Request) {
	auth := &types.AccessTokenRequest{}
	applog.Info("login in user")
	if !auth.Validate(w, r) {
		return
	}
	accesstoken := auth.GenerateAccessToken(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": map[string]interface{}{
		"access_token": accesstoken.Token,
		"expires_at":   accesstoken.ExpiresAt,
	}, "status": 1}
	json.NewEncoder(w).Encode(response) 
}

// LogOut : logout from account
func LogOut(w http.ResponseWriter, r *http.Request) {
	accessToken := &types.AccessToken{}
	if !accessToken.AuthorizeByToken(w, r) {
		return
	}
	user := accessToken.GetUser()
	accessToken.Remove()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": map[string]interface{}{"user_id": user.ID, "message": "LoggedOut Successfully"}, "status": 1}
	json.NewEncoder(w).Encode(response)
}
 