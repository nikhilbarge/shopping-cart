package service

import (
	"errors"
	"shopping-cart/pkg/database"
	"shopping-cart/types"
	"shopping-cart/utils/applog"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jameskeane/bcrypt"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

type AuthService struct{
	dbsrv database.IDataService
} 
func (as *AuthService) NewAuthService() *AuthService {
	dbservice := database.NewDBService()
	return &AuthService{dbsrv: dbservice}
}
var user types.User
//GetUser : return user object
func (auth *AuthService) GetUser() types.User {
	user.Password = ""
	return user
}

//Remove : remove access token
func (auth *AuthService) Remove(accessToken *types.AccessToken) bool {
	err := auth.dbsrv.RemoveToken(accessToken)
	if err != nil {
		return false
	} else {
		return true
	}
}

//AuthorizeByToken : Authorize api calls by token
func (auth *AuthService) AuthorizeByToken(accessToken *types.AccessToken) error {
	err := auth.dbsrv.GetAuthToken(accessToken) 
	if err != nil { 
		return err
	} 
	err = auth.dbsrv.GetUserByID(accessToken, &user)
	if err != nil { 
		return err
	} 
	return nil
}

//GenerateAccessToken : generate and return accessToken Object
func (auth *AuthService) GenerateAccessToken(accessTokenRequest *types.AccessTokenRequest) (*types.AccessToken,error) {
	accessToken := &types.AccessToken{}
	accessToken.ID = bson.NewObjectId()
	accessToken.Token = uuid.Must(uuid.NewV4(),nil).String()
	accessToken.UserID = user.ID
	accessToken.ExpiresAt = time.Now().Local().Add(time.Hour*time.Duration((24*60)) +
		time.Minute*time.Duration(0) +
		time.Second*time.Duration(0))
	accessToken.CreatedAt = time.Now().Local()
	accessToken.UpdatedAt = time.Now().Local()

	
	// Insert
	insertionErrors := auth.dbsrv.InsertToken(accessToken)

	if insertionErrors != nil {
		applog.Errorf("error while generating token, err %v",insertionErrors)
		return nil, insertionErrors
	}

	return accessToken, nil
}
 
func (auth *AuthService) Validate(req *types.AccessTokenRequest) error {
	
	applog.Info("validate user login details")  
	if govalidator.IsNull(req.Username) { 
		applog.Error("In valid user username")
		return errors.New("in valida username")
	}
	if govalidator.IsNull(req.Password) { 
		applog.Error("In valid user password")
		return errors.New("in valida password")
	}

	applog.Infof("username and password provided")
	if !govalidator.IsNull(req.Password) && !govalidator.IsNull(req.Username) {
		applog.Infof("get username and password from database")
		err := auth.dbsrv.GetUserByName(req.Username, &user)
		if err != nil {
			applog.Errorf("Error while fetching user details %v", err)
			return errors.New("Username or Password is wrong")
		} else {
			if !bcrypt.Match(req.Password, user.Password) {
				applog.Error("Password mismatch,In valid user password")
				return errors.New("Username or Password is wrong")
			}
		}
	}
	return nil
}