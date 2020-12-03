package user

import (
	"encoding/json"
	"net/http"
	"shopping-cart/pkg/service"
	"shopping-cart/types"
	"shopping-cart/utils/applog"
)

// RegisterUser : Register user account
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	user := &types.User{}
	applog.Info("Register new user")
	userService := service.UserService{} 
	if userService.Validate(w, r, user) && userService.RegisterUser(w, r, user){
		applog.Debugf("User '%s' created successfully",user.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{"data": user, "status": 1}
		json.NewEncoder(w).Encode(response)
	} 
	applog.Info("Register of new user completed")
}

// Login : login to account
func Login(w http.ResponseWriter, r *http.Request) {
	auth := &types.AccessTokenRequest{}
	applog.Info("login in user")
	authService:= service.AuthService{}
	if !authService.Validate(w, r, auth) {
		applog.Debug("unable to allow login")
		return
	}
	applog.Info("generating token for user")
	accesstoken := authService.GenerateAccessToken(w, auth)
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
	accessToken := &types.AccessToken{}
	accessToken.GetTokenFromRequest(r)
	authService := &service.AuthService{}
	if !authService.AuthorizeByToken(w, accessToken) {
		return
	} 
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