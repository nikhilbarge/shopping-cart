package service

import (
	"encoding/json"
	"net/http"
	"net/url"
	"shopping-cart/pkg/database"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"github.com/jameskeane/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type UserService struct {
	DbService database.DataService
}
func (userservices *UserService) Validate(w http.ResponseWriter, r *http.Request,user *types.User) bool {
	errs := url.Values{} 
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
	usr := types.User{}
 	err := userservices.DbService.GetUserByEmail(user.Email, &usr) 
	if err != nil && err.Error() !=  "not found" {
		applog.Error("failed to fetch user details by email", err)
		errs.Add("user", "unable to get user data")
	} else if usr.ID != ""{
		errs.Add("email", "user with this email is already exists")
	}
	usr2 := types.User{}
	err = userservices.DbService.GetUserByName(user.UserName, &usr2)  
	if err != nil && err.Error() !=  "not found" {
		applog.Error("failed to fetch user details by username", err)
		errs.Add("user", "unable to get user data")
	} else if usr2.ID != "" {
		errs.Add("username", "username is already exists")
	}
 
	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	}
	applog.Info("user data validated successfully")
	return true
}


func (userservices *UserService) RegisterUser(w http.ResponseWriter, r *http.Request,user *types.User) bool {
	// Insert

	applog.Info("registering user")
	cart := &types.Cart{}
	cart.ID = bson.NewObjectId()
	cart.Items= []types.CartItem{}

	applog.Info("create new cart for user")
	cartErr := userservices.DbService.InsertCart(cart)
	if cartErr != nil {
		applog.Errorf("failed to create cart for user, err %v ", cartErr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": cartErr.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		return false

	}
	user.ID = bson.NewObjectId()
	user.CartID = cart.ID
	salt, _ := bcrypt.Salt(10)
	user.Password, _ = bcrypt.Hash(user.Password, salt)
	applog.Info("create new user")
	insertionErrors := userservices.DbService.InsertUser(user)
	if insertionErrors != nil {
		applog.Error("error occured while registration of new user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": insertionErrors.Error(), "status": 0}
		json.NewEncoder(w).Encode(response) 
		return false
	}
	applog.Info("user resgistered successfully")
	return true
}