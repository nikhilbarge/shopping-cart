package types

import (
	"encoding/json"
	"net/http"
	"net/url"
	"shopping-cart/pkg/database"
	"shopping-cart/utils/applog"
	time "time"

	"github.com/asaskevich/govalidator"
	"github.com/jameskeane/bcrypt"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)
var user User
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

//GetUser : return user object
func (accessToken *AccessToken) GetUser() User {
	user.Password = ""
	return user
}

//Remove : remove access token
func (accessToken *AccessToken) Remove() bool {

	db := database.Db
	err := db.C("accesstoken").Remove(&accessToken)
	if err != nil {
		return false
	} else {
		return true
	}

}

//AuthorizeByToken : Authorize api calls by token
func (accessToken *AccessToken) AuthorizeByToken(w http.ResponseWriter, r *http.Request) bool {

	keys, ok := r.URL.Query()["access_token"]
	token := ""
	if !ok || len(keys[0]) < 1 {
		token = r.Header.Get("access_token")
	} else {
		token = keys[0]
	}

	if govalidator.IsNull(token) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": map[string]interface{}{"access_token": "Access token is required"}, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	}

	db := database.Db
	now := time.Now().Local()
	err := db.C("accesstoken").Find(bson.M{"token": token, "expires_at": bson.M{"$gt": now}}).One(&accessToken)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": map[string]interface{}{"access_token": "Invalid Access token"}, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	}

	err = db.C("user").Find(bson.M{"_id": accessToken.UserID}).One(&user)
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
func (accessTokenRequest *AccessTokenRequest) GenerateAccessToken(w http.ResponseWriter) *AccessToken {
	accessToken := &AccessToken{}
	accessToken.ID = bson.NewObjectId()
	accessToken.Token = uuid.Must(uuid.NewV4(),nil).String()
	accessToken.UserID = user.ID
	accessToken.ExpiresAt = time.Now().Local().Add(time.Hour*time.Duration((24*60)) +
		time.Minute*time.Duration(0) +
		time.Second*time.Duration(0))
	accessToken.CreatedAt = time.Now().Local()
	accessToken.UpdatedAt = time.Now().Local()

	db := database.Db
	c := db.C("accesstoken")
	// Insert
	insertionErrors := c.Insert(&accessToken)

	if insertionErrors != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": insertionErrors.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)

	}

	return accessToken
}
 
func (req *AccessTokenRequest) Validate(w http.ResponseWriter, r *http.Request) bool {
	errs := url.Values{}
	db := database.Db

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		applog.Errorf("Error validating user credentials %v", err)
		errs.Add("data", "Invalid data")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return false
	}

	if govalidator.IsNull(req.Username) {
		errs.Add("username", "Username is required")
		applog.Error("In valid user username")
	}
	if govalidator.IsNull(req.Password) {
		errs.Add("password", "Password is required")
		applog.Error("In valid user password")
	}

	if !govalidator.IsNull(req.Password) && !govalidator.IsNull(req.Username) {

		err := db.C("user").Find(bson.M{"username": req.Username}).One(&user)
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