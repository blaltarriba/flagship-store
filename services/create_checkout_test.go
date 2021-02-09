package services

import (
	"lana/flagship-store/models"
	"lana/flagship-store/services/commands"
	"lana/flagship-store/services/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type checkoutRepositoryMock struct {
	mock.Mock
}

func (repository *checkoutRepositoryMock) SearchById(id string) (models.Checkout, bool) {
	args := repository.Called(id)
	return args.Get(0).(models.Checkout), args.Bool(1)
}

func (repository *checkoutRepositoryMock) Persist(checkout models.Checkout) {
	repository.Called(checkout)
	return
}

func (repository *checkoutRepositoryMock) Delete(checkout models.Checkout) {
	repository.Called(checkout)
	return
}

func (repository *checkoutRepositoryMock) Count() int {
	args := repository.Called()
	return args.Int(0)
}

type productRepositoryMock struct {
	mock.Mock
}

func (repository *productRepositoryMock) SearchById(id string) (models.Product, bool) {
	args := repository.Called(id)
	return args.Get(0).(models.Product), args.Bool(1)
}

func TestCreateCheckout(t *testing.T) {
	theCheckoutRepositoryMock := checkoutRepositoryMock{}
	theCheckoutRepositoryMock.On("Persist", mock.AnythingOfType("models.Checkout"))

	theProductRepositoryMock := productRepositoryMock{}
	theProductRepositoryMock.On("SearchById", "PEN").Return(models.Product{}, true)

	productCommand := commands.Product{Code: "PEN"}
	createCheckout := CreateCheckout{&theCheckoutRepositoryMock, &theProductRepositoryMock}

	createdCheckout, _ := createCheckout.Do(productCommand)

	assert.NotNil(t, createdCheckout.Id)
	assert.EqualValues(t, "PEN", createdCheckout.Products[0])
	assert.EqualValues(t, 1, len(createdCheckout.Products))
	theCheckoutRepositoryMock.AssertNumberOfCalls(t, "Persist", 1)
	theCheckoutRepositoryMock.AssertExpectations(t)
}

func TestReturnProductNotFoundErrorWhenProductDoesnotExists(t *testing.T) {
	theCheckoutRepositoryMock := checkoutRepositoryMock{}
	theCheckoutRepositoryMock.On("Persist", mock.AnythingOfType("models.Checkout"))

	theProductRepositoryMock := productRepositoryMock{}
	theProductRepositoryMock.On("SearchById", "PEN").Return(models.Product{}, false)

	productCommand := commands.Product{Code: "PEN"}
	createCheckout := CreateCheckout{&theCheckoutRepositoryMock, &theProductRepositoryMock}

	_, err := createCheckout.Do(productCommand)

	_, isProductNotFoundError := err.(*errors.ProductNotFoundError)
	assert.EqualValues(t, true, isProductNotFoundError)
	theCheckoutRepositoryMock.AssertNumberOfCalls(t, "Persist", 0)
}
