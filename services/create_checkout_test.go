package services

import (
	"lana/flagship-store/models"
	"lana/flagship-store/services/commands"
	"lana/flagship-store/services/errors"
	"lana/flagship-store/utils/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCheckout(t *testing.T) {
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("Persist", mock.AnythingOfType("models.Checkout"))

	theProductRepositoryMock := mocks.ProductRepositoryMock{}
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
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("Persist", mock.AnythingOfType("models.Checkout"))

	theProductRepositoryMock := mocks.ProductRepositoryMock{}
	theProductRepositoryMock.On("SearchById", "PEN").Return(models.Product{}, false)

	productCommand := commands.Product{Code: "PEN"}
	createCheckout := CreateCheckout{&theCheckoutRepositoryMock, &theProductRepositoryMock}

	_, err := createCheckout.Do(productCommand)

	_, isProductNotFoundError := err.(*errors.ProductNotFoundError)
	assert.EqualValues(t, true, isProductNotFoundError)
	theCheckoutRepositoryMock.AssertNumberOfCalls(t, "Persist", 0)
}
