package service

import (
	"shopping-cart/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)
 
func (suite *ServiceTestSuite) SetupTest() {
	suite.clientMock = new(MockApiClient)
	suite.serviceMock = new(MockApiService)
}

type ServiceTestSuite struct {
	suite.Suite
	clientMock  *MockApiClient
	serviceMock *MockApiService
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
func (suite *ServiceTestSuite) Test_Validate_Success() { 
	suite.clientMock.On("GetUserByEmail").Return(nil)
	suite.clientMock.On("GetUserByName").Return(nil)
	service := UserService{dbsrv: suite.clientMock}

	// Act 
	response := service.Validate(&types.User{Name: "nik",Email: "nik@gmail.com",Password: "nik",UserName: "nik"})

	// Assert
	assert.Nil(suite.T(), response, "Response should be nil")
	 
}


func (suite *ServiceTestSuite) Test_RegisterUser_Success() { 
	suite.clientMock.On("InsertCart").Return(nil)
	suite.clientMock.On("InsertUser").Return(nil)
	service := UserService{dbsrv: suite.clientMock}

	// Act 
	response := service.RegisterUser(&types.User{Name: "nik",Email: "nik@gmail.com",Password: "nik",UserName: "nik"})

	// Assert
	assert.Nil(suite.T(), response, "Response should be nil")
	 
}


func (suite *ServiceTestSuite) Test_ValidateCart_Success() { 
	suite.clientMock.On("GetCartByID").Return(nil)
	suite.clientMock.On("GetItemByID").Return(nil)
	service := CartService{dbsrv: suite.clientMock}

	cart := types.Cart{ID: "5fc966a74d278b000141901f"}
	item := types.Item{Name: "shoe", Price: 2000, Quantity: 1 }
	// Act 
	response := service.Validate(&item, &cart)

	// Assert
	assert.Nil(suite.T(), response, "Response should be nil")
	 
}


func (suite *ServiceTestSuite) Test_AddToCart_Success() { 
	
	suite.clientMock.On("GetCartByID").Return(nil) 
	suite.clientMock.On("GetItemByID").Return(nil)
	suite.clientMock.On("UpdateCart").Return(nil) 
	service := CartService{dbsrv: suite.clientMock} 

	cart := types.Cart{ID: "5fc966a74d278b000141901f"}
	item := types.Item{Name: "shoe", Price: 2000, Quantity: 1 }
	response := service.Validate(&item, &cart)
	// Act 
	response = service.AddToCart(&cart)

	// Assert
	assert.Nil(suite.T(), response, "Response should be nil")
}



func (suite *ServiceTestSuite) Test_ViewCart_Success() { 
	
	suite.clientMock.On("GetCartByID").Return(nil) 
	service := CartService{dbsrv: suite.clientMock} 
  
	// Act 
	_, err := service.ViewAllCarts("5fc966a74d278b000141901d")

	// Assert
	assert.Nil(suite.T(), err, "Error should be nil")
}

func (suite *ServiceTestSuite) Test_RemoveItem_Success() { 
	
	suite.clientMock.On("GetCartByID").Return(nil) 
	suite.clientMock.On("UpdateCart").Return(nil) 
	service := CartService{dbsrv: suite.clientMock} 

	cart := types.Cart{ID: "5fc966a74d278b000141901f"}
	 
	// Act 
	response := service.RemoveItem(&cart,"5fc966eb182897aed9b4bfa7")

	// Assert
	assert.Nil(suite.T(), response, "Response should be nil")
}


func (suite *ServiceTestSuite) Test_ClearCart_Success() { 
	
	suite.clientMock.On("GetCartByID").Return(nil) 
	suite.clientMock.On("DeleteCart").Return(nil) 
	service := CartService{dbsrv: suite.clientMock} 

	// Act 
	response := service.DeleteCart("5fc966a74d278b000141901f")

	// Assert
	assert.Nil(suite.T(), response, "Response should be nil")
}



func (suite *ServiceTestSuite) Test_AddToInventory_Success() { 
	
	suite.clientMock.On("GetItemByName").Return(nil) 
	suite.clientMock.On("UpdateItemByID").Return(nil) 
	suite.clientMock.On("InsertItem").Return(nil) 
	service := InventoryService{dbsrv: suite.clientMock} 

	item := types.Item{ID: "5fc966a74d278b000141901f"}
	 
	// Act 
	response := service.AddToInventory(&item)

	// Assert
	assert.Nil(suite.T(), response, "Response should be nil")
}

// ViewInvetory
func (suite *ServiceTestSuite) Test_ViewInvetory_Success() { 
	
	suite.clientMock.On("GetAllItems").Return(nil,nil)  
	service := InventoryService{dbsrv: suite.clientMock} 
 
	// Act 
	items := types.ItemList{}
	response := service.ViewInvetory(&items)

	// Assert
	assert.Nil(suite.T(), response, "Response should be nil")
}
// RemoveItem
func (suite *ServiceTestSuite) Test_Inventory_RemoveItem_Success() { 
	
	suite.clientMock.On("RemoveItem").Return(nil) 
	suite.clientMock.On("RemoveAllItem").Return(nil)  
	service := InventoryService{dbsrv: suite.clientMock} 

	item := types.Item{}
	 
	// Act 
	response := service.RemoveItem(&item, "5fc966a74d278b000141901f")

	// Assert
	assert.Nil(suite.T(), response, "Response should be nil")
}