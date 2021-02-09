package services

import (
	"lana/flagship-store/models"
	"lana/flagship-store/services/commands"
	"lana/flagship-store/services/errors"
	"lana/flagship-store/utils/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddProductToCheckout(t *testing.T) {
	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG"},
	}
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", checkout.Id).Return(checkout, true)
	theCheckoutRepositoryMock.On("Persist", mock.AnythingOfType("models.Checkout"))
	theProductRepositoryMock := mocks.ProductRepositoryMock{}
	theProductRepositoryMock.On("SearchById", "PEN").Return(models.Product{}, true)
	addProductCommand := commands.AddProduct{Code: "PEN"}
	addProductToCheckout := AddProductToCheckout{&theCheckoutRepositoryMock, &theProductRepositoryMock}

	modifiedCheckout, _ := addProductToCheckout.Do(addProductCommand, checkout.Id)

	assert.NotNil(t, modifiedCheckout.Id)
	assert.EqualValues(t, "MUG", modifiedCheckout.Products[0])
	assert.EqualValues(t, "PEN", modifiedCheckout.Products[1])
	assert.EqualValues(t, 2, len(modifiedCheckout.Products))
	theCheckoutRepositoryMock.AssertNumberOfCalls(t, "Persist", 1)
	theCheckoutRepositoryMock.AssertExpectations(t)
}

func TestAddProductReturnCheckoutNotFoundErrorWhenCheckoutDoesnotExists(t *testing.T) {
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", mock.AnythingOfType("string")).Return(models.Checkout{}, false)
	theProductRepositoryMock := mocks.ProductRepositoryMock{}
	addProductCommand := commands.AddProduct{Code: "PEN"}
	addProductToCheckout := AddProductToCheckout{&theCheckoutRepositoryMock, &theProductRepositoryMock}

	_, err := addProductToCheckout.Do(addProductCommand, "a_fake_id")

	_, isCheckoutNotFoundError := err.(*errors.CheckoutNotFoundError)
	assert.EqualValues(t, true, isCheckoutNotFoundError)
}

func TestAddProductReturnProductNotFoundErrorWhenProductDoesnotExists(t *testing.T) {
	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG"},
	}
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", checkout.Id).Return(checkout, true)
	theProductRepositoryMock := mocks.ProductRepositoryMock{}
	theProductRepositoryMock.On("SearchById", "PEN").Return(models.Product{}, false)
	addProductCommand := commands.AddProduct{Code: "PEN"}
	addProductToCheckout := AddProductToCheckout{&theCheckoutRepositoryMock, &theProductRepositoryMock}

	_, err := addProductToCheckout.Do(addProductCommand, checkout.Id)

	_, isProductNotFoundError := err.(*errors.ProductNotFoundError)
	assert.EqualValues(t, true, isProductNotFoundError)
}
