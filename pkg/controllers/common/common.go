package common

import (
	"encoding/json"
	"errors"
	"net/http"
	"shopping-cart/pkg/service"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"github.com/asaskevich/govalidator"
)


func CheckAuthorized(w http.ResponseWriter, r *http.Request) (*types.AccessToken, error) {
	accessToken := &types.AccessToken{}
	keys, ok := r.URL.Query()["access_token"] 
	if !ok || len(keys[0]) < 1 {
		accessToken.Token = r.Header.Get("access_token")
	} else {
		accessToken.Token = keys[0]
	}
	as := service.AuthService{}
	authService := as.NewAuthService()  
	if govalidator.IsNull(accessToken.Token) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": map[string]interface{}{"access_token": "Access token is required"}, "status": 0}
		json.NewEncoder(w).Encode(response)
		return nil, errors.New("access token required")
	} 
	 
	err := authService.AuthorizeByToken(accessToken) 
	if err!=nil{
		applog.Errorf("error while authorize with token, err %v",err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": map[string]interface{}{"access_token": "Invalid Access token"}, "status": 0}
		json.NewEncoder(w).Encode(response)
		return nil, errors.New("invalid access token")
	}
	return accessToken, nil
}