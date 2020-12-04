package service

import (
	"encoding/json"
	"net/http"
	"net/url"
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
func (auth *AuthService) AuthorizeByToken(w http.ResponseWriter, accessToken *types.AccessToken) bool {
	if govalidator.IsNull(accessToken.Token) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": map[string]interface{}{"access_token": "Access token is required"}, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	} 
	err := auth.dbsrv.GetAuthToken(accessToken) 
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": map[string]interface{}{"access_token": "Invalid Access token"}, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	}

	err = auth.dbsrv.GetUserByID(accessToken, &user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": map[string]interface{}{"access_token": "Invalid User Record for this Access token"}, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	}

	return true
}

//GenerateAccessToken : generate and return accessToken Object
func (auth *AuthService) GenerateAccessToken(w http.ResponseWriter, accessTokenRequest *types.AccessTokenRequest) *types.AccessToken {
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": insertionErrors.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)

	}

	return accessToken
}
 
func (auth *AuthService) Validate(w http.ResponseWriter, r *http.Request,req *types.AccessTokenRequest) bool {
	errs := url.Values{} 
	applog.Info("validate user login details")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		applog.Errorf("Error validating user credentials %v", err)
		errs.Add("data", "Invalid data")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	}
	applog.Infof("username and password found")

	if govalidator.IsNull(req.Username) {
		errs.Add("username", "Username is required")
		applog.Error("In valid user username")
	}
	if govalidator.IsNull(req.Password) {
		errs.Add("password", "Password is required")
		applog.Error("In valid user password")
	}

	applog.Infof("username and password provided")
	if !govalidator.IsNull(req.Password) && !govalidator.IsNull(req.Username) {
		applog.Infof("get username and password from database")
		err := auth.dbsrv.GetUserByName(req.Username, &user)
		if err != nil {
			applog.Errorf("Error while fetching user details %v", err)
			errs.Add("password", "Username or Password is wrong")
		} else {
			if !bcrypt.Match(req.Password, user.Password) {
				applog.Error("Password mismatch,In valid user password")
				errs.Add("password", "Username or Password is wrong")
			}
		}
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