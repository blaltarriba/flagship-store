package services

import (
	"lana/flagship-store/models"
	"lana/flagship-store/persistence"
	"lana/flagship-store/services/errors"
)

type DeleteCheckout struct {
	CheckoutRepository persistence.CheckoutRepository
}

func NewDeleteCheckout(checkoutRepository persistence.CheckoutRepository) DeleteCheckout {
	return DeleteCheckout{checkoutRepository}
}

func (service *DeleteCheckout) Do(checkoutId string) (models.Checkout, error) {
	checkout, existCheckout := service.CheckoutRepository.SearchById(checkoutId)
	if !existCheckout {
		return models.Checkout{}, errors.NewCheckoutNotFoundError()
	}

	service.CheckoutRepository.Delete(checkout)

	return checkout, nil
}
