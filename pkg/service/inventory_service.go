package service

import (
	"errors"
	"shopping-cart/pkg/database"
	"shopping-cart/types"
	"shopping-cart/utils/applog"

	"gopkg.in/mgo.v2/bson"
)

type InventoryService struct{
	dbsrv database.IDataService
} 
func (is *InventoryService) NewInventoryService() *InventoryService {
	dbservice := database.NewDBService()
	return &InventoryService{dbsrv: dbservice}
}
// AddToInventory : adding/updating Inventory
func (inventory *InventoryService) AddToInventory(item *types.Item) error {
	applog.Info("validating add to inventory request")
	applog.Infof("processing request to add item %s to inventory", item.Name)
	oldItem:= types.Item{}
	err := inventory.dbsrv.GetItemByName(item.Name,&oldItem)
	if err!=nil && err.Error()!="not found" {
		applog.Errorf("failed to fetch item %s err %v", item.Name,err)
		return err
	}
	if oldItem.ID != "" {
		applog.Debugf("item %s already exists in inventory, increase in quantity ",item.Name)
		oldItem.Quantity += item.Quantity
		err := inventory.dbsrv.UpdateItemByID(oldItem.ID,&oldItem)
		if err != nil { 
			applog.Errorf("failed to update item %s to invetory ",oldItem.Name)
			return err
		} 
	} else {
		err := inventory.dbsrv.InsertItem(item)
		if err != nil { 
			applog.Errorf("failed to add item %s to invetory err %v ",item.Name, err)
			return err
		} 
	}
	inventory.dbsrv.GetItemByName(item.Name, item)   
	applog.Info("inventory updated successfully")
	return nil
}
 

//ViewInvetory : Find All Invetory records
func (inventory *InventoryService) ViewInvetory(items *types.ItemList) error { 
 
	applog.Info("get all items in invetory")
	 
	itemList, err := inventory.dbsrv.GetAllItems()
	if err != nil { 
		applog.Errorf("error while fetching items in inventory")
		return  err
	} 
	items.Items = itemList
	applog.Debugf("successfully fetched items in invetory")
	return nil
}

// RemoveItem : Remove Invetory record
func (inventory *InventoryService) RemoveItem(item *types.Item, id string) error { 
	applog.Infof("remove all items in inventory")
	if id != "" && bson.IsObjectIdHex(id) {
		err := inventory.dbsrv.RemoveItem(item.ID)
		if err != nil { 
			applog.Errorf("error while removing items from inventory %s",item.ID)
			return err
		} 
		applog.Infof("removed all item %s from inventory",item.ID)
		return nil
	} else if id == "" {
		err := inventory.dbsrv.RemoveAllItem()
		if err != nil { 
			applog.Errorf("error while removing items from inventory")
			return err
		} 
		applog.Infof("removed all items from inventory")
		return nil
	} else {
		applog.Debugf("invalid item id %s",item.ID)
		return errors.New("Invalid item id") 
	}    
}
 