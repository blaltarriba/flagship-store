package services

import (
	"lana/flagship-store/models"
	"lana/flagship-store/services/errors"
	"lana/flagship-store/utils/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCheckout(t *testing.T) {
	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG"},
	}
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", checkout.Id).Return(checkout, true)
	theCheckoutRepositoryMock.On("Delete", checkout)
	deleteCheckout := DeleteCheckout{&theCheckoutRepositoryMock}

	_, err := deleteCheckout.Do(checkout.Id)

	assert.Nil(t, err)
	theCheckoutRepositoryMock.AssertNumberOfCalls(t, "Delete", 1)
	theCheckoutRepositoryMock.AssertExpectations(t)
}

func TestDeleteReturnCheckoutNotFoundErrorWhenCheckoutDoesnotExists(t *testing.T) {
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", "a_fake_id").Return(models.Checkout{}, false)
	deleteCheckout := DeleteCheckout{&theCheckoutRepositoryMock}

	_, err := deleteCheckout.Do("a_fake_id")

	_, isCheckoutNotFoundError := err.(*errors.CheckoutNotFoundError)
	assert.EqualValues(t, true, isCheckoutNotFoundError)
	theCheckoutRepositoryMock.AssertExpectations(t)
}
