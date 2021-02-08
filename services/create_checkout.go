package services

import (
	"lana/flagship-store/models"
	"lana/flagship-store/persistence"
	"lana/flagship-store/services/commands"
	"lana/flagship-store/services/errors"

	"github.com/google/uuid"
)

type CreateCheckout struct {
	CheckoutRepository persistence.CheckoutRepository
	ProductRepository  persistence.ProductRepository
}

func NewCreateCheckout(checkoutRepository persistence.CheckoutRepository, productRepository persistence.ProductRepository) CreateCheckout {
	return CreateCheckout{checkoutRepository, productRepository}
}

func (service *CreateCheckout) Do(productCommand commands.Product) (models.Checkout, error) {
	if _, existProduct := service.ProductRepository.SearchById(productCommand.Code); !existProduct {
		emptyCheckout := models.Checkout{}
		return emptyCheckout, errors.NewProductNotFoundError()
	}

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{productCommand.Code},
	}
	service.CheckoutRepository.Persist(checkout)

	return checkout, nil
}
