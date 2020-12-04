package user

import (
	"encoding/json"
	"net/http"
	"net/url"
	"shopping-cart/pkg/controllers/common"
	"shopping-cart/pkg/service"
	"shopping-cart/types"
	"shopping-cart/utils/applog"
)

// RegisterUser : Register user account
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	user := &types.User{}
	applog.Info("Register new user")
	
	errs := url.Values{}  
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		applog.Errorf("invalid user info err %v", err)
		errs.Add("detail", "Invalid data") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return 
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
	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}

	us := service.UserService{}
	userService := us.NewUserService()
	err := userService.Validate(user)
	if err!=nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotAcceptable)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}
	err = userService.RegisterUser( user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return 
	} 
	applog.Debugf("User '%s' created successfully",user.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": user, "status": 1}
	json.NewEncoder(w).Encode(response)
	applog.Info("Register of new user completed")
}

// Login : login to account
func Login(w http.ResponseWriter, r *http.Request) {
	auth := &types.AccessTokenRequest{}
	applog.Info("login in user")
	as := service.AuthService{}
	authService := as.NewAuthService() 
	errs := url.Values{} 
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		applog.Errorf("Error validating user credentials %v", err)
		errs.Add("data", "Invalid data") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return  
	}
	err := authService.Validate(auth) 
	if err!=nil{
		applog.Debug("unable to allow login")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": err, "status": 0}
		json.NewEncoder(w).Encode(response)
		return 
	}
	applog.Info("generating token for user")
	accesstoken,err := authService.GenerateAccessToken(auth)
	if err!=nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": map[string]interface{}{
		"access_token": accesstoken.Token,
		"expires_at":   accesstoken.ExpiresAt,
	}, "status": 1}
	json.NewEncoder(w).Encode(response) 
	applog.Debug("user logged in success fully")
}

// LogOut : logout from account
func LogOut(w http.ResponseWriter, r *http.Request) {
	// authenticating user
	accessToken,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	as := service.AuthService{}
	authService := as.NewAuthService() 
 
	user := authService.GetUser()
	authService.Remove(accessToken)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": map[string]interface{}{"user_id": user.ID, "message": "LoggedOut Successfully"}, "status": 1}
	json.NewEncoder(w).Encode(response)
}
 
func AppDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": map[string]interface{}{"app": "shopping cart", "version": "1.0.1"}, "status": 1}
	json.NewEncoder(w).Encode(response)
}