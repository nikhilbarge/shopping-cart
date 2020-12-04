package service

import (
	"shopping-cart/pkg/database"
	"shopping-cart/types"

	"github.com/stretchr/testify/mock"
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
