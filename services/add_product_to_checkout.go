package services

import (
	"lana/flagship-store/models"
	"lana/flagship-store/persistence"
	"lana/flagship-store/services/commands"
	"lana/flagship-store/services/errors"
)

type AddProductToCheckout struct {
	CheckoutRepository persistence.CheckoutRepository
	ProductRepository  persistence.ProductRepository
}

func NewAddProductToCheckout(checkoutRepository persistence.CheckoutRepository, productRepository persistence.ProductRepository) AddProductToCheckout {
	return AddProductToCheckout{checkoutRepository, productRepository}
}

func (service *AddProductToCheckout) Do(addProductCommand commands.AddProduct, checkoutId string) (models.Checkout, error) {
	checkout, existCheckout := service.CheckoutRepository.SearchById(checkoutId)
	if !existCheckout {
		return models.Checkout{}, errors.NewCheckoutNotFoundError()
	}

	if _, existProduct := service.ProductRepository.SearchById(addProductCommand.Code); !existProduct {
		return models.Checkout{}, errors.NewProductNotFoundError()
	}

	checkout.Products = append(checkout.Products, addProductCommand.Code)
	service.CheckoutRepository.Persist(checkout)

	return checkout, nil

}
