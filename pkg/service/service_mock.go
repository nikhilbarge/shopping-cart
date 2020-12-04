package service

import (
	"shopping-cart/pkg/database"
	"shopping-cart/types"

	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
)

type MockApiService struct {
	mock.Mock 
}

type MockApiClient struct {
	mock.Mock
	database.IDataService
}

//GetUserByEmail : mock for GetUserByEmail
func (m *MockApiClient) GetUserByEmail(email string, user *types.User) (error) {
	args := m.Called() 
	err, _ := args.Get(0).(error)
	return  err
} 

func (m *MockApiClient) GetUserByName(email string, user *types.User) (error) {
	args := m.Called() 
	err, _ := args.Get(0).(error)
	return  err
}



func (m *MockApiClient) InsertCart(cart *types.Cart) error {
	args := m.Called() 
	err, _ := args.Get(0).(error)
	return  err
} 

func (m *MockApiClient) InsertUser(user *types.User) error {
	args := m.Called() 
	err, _ := args.Get(0).(error)
	return  err
}





func (m *MockApiClient) GetCartByID(cartid bson.ObjectId, cart *types.Cart) error {
	args := m.Called() 
	err, _ := args.Get(0).(error)
	return  err
} 


func (m *MockApiClient) UpdateCart(cartid bson.ObjectId, cart *types.Cart) error {
	args := m.Called() 
	err, _ := args.Get(0).(error)
	return  err
}

func (m *MockApiClient) GetItemByID(itemid bson.ObjectId, item *types.Item) error{
	args := m.Called() 
	err, _ := args.Get(0).(error)
	return  err
}

func (m *MockApiClient) GetItemByName(itemname string, item *types.Item) error {
	args := m.Called() 
	err, _ := args.Get(0).(error)
	return  err
}

func (m *MockApiClient) InsertItem(item *types.Item) error {
	args := m.Called() 
	err, _ := args.Get(0).(error)
	return  err
}

func (m *MockApiClient) UpdateItemByID(itemid bson.ObjectId, item *types.Item) error {
	args := m.Called() 
	err, _ := args.Get(0).(error)
	return  err
}


func (m *MockApiClient) RemoveItem(itemid bson.ObjectId) error {
	args := m.Called() 
	err, _ := args.Get(0).(error)
	return  err
}

func (m *MockApiClient) RemoveAllItem() error  {
	args := m.Called() 
	err, _ := args.Get(0).(error)
	return  err
}


func (m *MockApiClient) GetAllItems() ([]types.Item,error)  {
	args := m.Called()
	items, _ := args.Get(0).([]types.Item) 
	err, _ := args.Get(1).(error)
	return  items,err
}