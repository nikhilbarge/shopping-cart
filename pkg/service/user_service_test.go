package service

import (
	"net/http"
	"shopping-cart/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func (suite *ApiTestSuite) SetupTest() {
	suite.clientMock = new(MockApiClient)
	suite.serviceMock = new(MockApiService)
}

type ApiTestSuite struct {
	suite.Suite
	clientMock  *MockApiClient
	serviceMock *MockApiService
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
func (suite *ApiTestSuite) Test_Validate_Success() { 
	expectedResponse := true
	suite.clientMock.On("GetUserByEmail").Return(expectedResponse, nil)
	service := UserService{dbsrv: suite.clientMock}

	// Act 
	response := service.Validate(nil ,&http.Request{}, &types.User{})

	// Assert
	assert.NotNil(suite.T(), response, "Response should not be nil")
	 
}