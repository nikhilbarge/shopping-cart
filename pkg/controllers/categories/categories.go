package categories

import (
	"encoding/json"
	"net/http"
	"net/url"
	"shopping-cart/pkg/controllers/common"
	"shopping-cart/pkg/service"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"github.com/gorilla/mux"
)

// AddItemToCategories : handler function for POST /v1/categories call
func AddCategories(w http.ResponseWriter, r *http.Request) {
	// authenticating user
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	as := service.AuthService{}
	authService := as.NewAuthService()  
	errs := url.Values{}
	if authService.GetUser().UserName!= "admin" {
		errs.Add("data", "User does not have access to update categories") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}
	 
	catagory := &types.Categories{} 
	if err := json.NewDecoder(r.Body).Decode(catagory); err != nil {
		errs.Add("data", "Invalid data") 
		applog.Errorf("invalid request for add catagory to categories, %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return 
	} 
	applog.Info("adding catagory to catagory")
	is := service.CategoriesService{} 
	categoriesService := is.NewCategoriesService() 


	err = categoriesService.AddCategories(catagory)
	if err!=nil { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": catagory, "status": 1}
	json.NewEncoder(w).Encode(response)

	applog.Info("add to categories request completed")
}
// ViewCategories Get All catagorys in categories
func ViewCategories(w http.ResponseWriter, r *http.Request) {
	// authenticating user  
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	catagoryList := &types.CategoryList{} 
	applog.Info("get all catagories")
	is := service.CategoriesService{}
	categoriesService := is.NewCategoriesService()
	err = categoriesService.ViewCategories(catagoryList) 
	if err!=nil{ 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{"errors": err, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": catagoryList, "status": 1}
	json.NewEncoder(w).Encode(response)
}

// RemoveItem : delete catagory from catagory
func RemoveCategory(w http.ResponseWriter, r *http.Request) {
	// authenticating user
	_,err := common.CheckAuthorized(w, r)
	if err!=nil {
		return
	}
	as := service.AuthService{}
	authService := as.NewAuthService()  
	errs := url.Values{}
	if authService.GetUser().UserName!= "admin" {
		errs.Add("data", "Forbiden, user is not 'admin'") 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		response := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	}
	params := mux.Vars(r)
	catagory := &types.Categories{}

	is := service.CategoriesService{}
	categoriesService := is.NewCategoriesService()
	err = categoriesService.RemoveCategory(catagory, params["catagoryid"])
	if err!=nil { 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		return
	} 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"data": catagory, "message": "Item Deleted Successfully", "status": 1}
	json.NewEncoder(w).Encode(response)

}
 