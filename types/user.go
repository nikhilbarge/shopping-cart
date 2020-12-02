package types

import (
	"encoding/json"
	"net/http"
	"net/url"
	"shopping-cart/pkg/database"
	"shopping-cart/utils/applog"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string        `bson:"name" json:"name"`
	UserName  string        `bson:"username" json:"username"`
	Email     string        `bson:"email" json:"email"`
	Password  string        `bson:"password" json:"password,omitempty"`
	CartID    bson.ObjectId  `bson:"cartid" json:"cartid"`
}

func (user *User) Validate(w http.ResponseWriter, r *http.Request) bool {
	errs := url.Values{}
	db := database.Db
	applog.Info("validating user info")
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		applog.Errorf("invalid user info err %v", err)
		errs.Add("detail", "Invalid data") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	}

	if user.Name == "" {
		errs.Add("name", "Name is not provided")
		applog.Error("Invalid info, empty user's name")
	}
	if user.UserName == "" {
		errs.Add("username", "Username is not provided")
		applog.Error("Invalid info, empty user's username")
	}
	if user.Email == "" {
		errs.Add("email", "E-mail is not provided")
		applog.Error("Invalid info, empty user's email")
	}
	if user.Password == "" {
		errs.Add("password", "Password is not provided")
		applog.Error("Invalid info, empty user's password")
	}

	count, _ := db.C("user").Find(bson.M{"email": user.Email}).Count()
	if count > 0 {
		errs.Add("email", "user with email-id already exists")
	}
	count, _ = db.C("user").Find(bson.M{"username": user.UserName}).Count()
	if count > 0 {
		errs.Add("username", "username is already exists")
	}

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	}
	return true
}
