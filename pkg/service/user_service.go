package service

import (
	"errors"
	"shopping-cart/pkg/database"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"github.com/jameskeane/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type UserService struct {
	dbsrv database.IDataService
}

func (us *UserService) NewUserService() *UserService {
	dbservice := database.NewDBService()
	return &UserService{dbsrv: dbservice}
}
func (us *UserService) Validate(user *types.User) error {
	applog.Info("validating user info") 
	usr := types.User{}
 	err := us.dbsrv.GetUserByEmail(user.Email, &usr) 
	if err != nil && err.Error() !=  "not found" {
		applog.Error("failed to fetch user details by email", err) 
		return err
	} else if usr.ID != ""{
		applog.Error("user with this email is already exists") 
		return errors.New("user with this email is already exists")
	}
	usr2 := types.User{}
	err = us.dbsrv.GetUserByName(user.UserName, &usr2)  
	if err != nil && err.Error() !=  "not found" {
		applog.Error("failed to fetch user details by username", err) 
		return err
	} else if usr2.ID != "" { 
		applog.Error("username is already exists") 
		return errors.New("username is already exists")
	} 
	applog.Info("user data validated successfully")
	return nil
}


func (us *UserService) RegisterUser(user *types.User) error {
	// Insert 
	applog.Info("registering user")
	cart := &types.Cart{}
	cart.ID = bson.NewObjectId()
	cart.Items= []types.CartItem{}

	applog.Info("create new cart for user")
	cartErr := us.dbsrv.InsertCart(cart)
	if cartErr != nil {
		applog.Errorf("failed to create cart for user, err %v ", cartErr)
		return cartErr    
	}
	user.ID = bson.NewObjectId()
	user.CartID = cart.ID
	salt, _ := bcrypt.Salt(10)
	user.Password, _ = bcrypt.Hash(user.Password, salt)
	applog.Info("create new user")
	err := us.dbsrv.InsertUser(user)
	if err != nil {
		applog.Error("error occured while registration of new user")
		 
		return err
	}
	applog.Info("user resgistered successfully")
	return nil
}