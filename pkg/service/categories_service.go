package service

import (
	"errors"
	"shopping-cart/pkg/database"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"gopkg.in/mgo.v2/bson"
)

type CategoriesService struct{
	dbsrv database.IDataService
} 
func (is *CategoriesService) NewCategoriesService() *CategoriesService {
	dbservice := database.NewDBService()
	return &CategoriesService{dbsrv: dbservice}
}
// AddToCategories : adding/updating Categories
func (categories *CategoriesService) AddCategories(catagory *types.Categories) error {
	applog.Info("validating add to categories request")
	oldCategory:= types.Categories{}
	err := categories.dbsrv.GetCategoriesByName(catagory.Name,&oldCategory)
	if err!=nil && err.Error()!="not found" {
		applog.Errorf("failed to fetch catagory %s err %v", catagory.Name,err)
		return err
	}
	if oldCategory.ID != "" {
		applog.Debugf("category %s already exists ",catagory.Name)
		return errors.New("already exists")
	} else {
		err := categories.dbsrv.InsertCategories(catagory)
		if err != nil { 
			applog.Errorf("failed to add catagory %s to invetory err %v ",catagory.Name, err)
			return err
		} 
	}
	categories.dbsrv.GetCategoriesByName(catagory.Name, catagory)   
	applog.Info("categories updated successfully")
	return nil
}
 

//ViewInvetory : Find All Invetory records
func (categories *CategoriesService) ViewCategories(category *types.CategoryList) error { 
 
	applog.Info("get all category in invetory")
	 
	catogoryList, err := categories.dbsrv.GetAllCategoriess()
	if err != nil { 
		applog.Errorf("error while fetching category in categories")
		return  err
	} 
	category.List = catogoryList
	applog.Debugf("successfully fetched category in invetory")
	return nil
}

// RemoveItem : Remove Invetory record
func (categories *CategoriesService) RemoveCategory(catagory *types.Categories, id string) error { 
	applog.Infof("remove all category in categories")
	if id != "" && bson.IsObjectIdHex(id) {
		err := categories.dbsrv.RemoveCategories(catagory.ID)
		if err != nil { 
			applog.Errorf("error while removing category from categories %s",catagory.ID)
			return err
		} 
		applog.Infof("removed all catagory %s from categories",catagory.ID)
		return nil
	} else if id == "" {
		err := categories.dbsrv.RemoveAllCategories()
		if err != nil { 
			applog.Errorf("error while removing category from categories")
			return err
		} 
		applog.Infof("removed all category from categories")
		return nil
	} else {
		applog.Debugf("invalid catagory id %s",catagory.ID)
		return errors.New("Invalid catagory id") 
	}    
}
 