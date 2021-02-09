package mocks

import (
	"lana/flagship-store/models"

	"github.com/stretchr/testify/mock"
)

type CheckoutRepositoryMock struct {
	mock.Mock
}

func (repository *CheckoutRepositoryMock) SearchById(id string) (models.Checkout, bool) {
	args := repository.Called(id)
	return args.Get(0).(models.Checkout), args.Bool(1)
}

func (repository *CheckoutRepositoryMock) Persist(checkout models.Checkout) {
	repository.Called(checkout)
	return
}

func (repository *CheckoutRepositoryMock) Delete(checkout models.Checkout) {
	repository.Called(checkout)
	return
}

func (repository *CheckoutRepositoryMock) Count() int {
	args := repository.Called()
	return args.Int(0)
}
