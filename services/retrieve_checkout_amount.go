package services

import (
	"lana/flagship-store/persistence"
	"lana/flagship-store/services/errors"
)

type RetrieveCheckoutAmount struct {
	CheckoutRepository             persistence.CheckoutRepository
	ProductRepository              persistence.ProductRepository
	ProductWithPromotionRepository persistence.ProductRepository
	ProductWithDiscountRepository  persistence.ProductRepository
}

func NewRetrieveCheckoutAmount(checkoutRepository persistence.CheckoutRepository, productRepository persistence.ProductRepository, productWithPromotionRepository persistence.ProductRepository, productWithDiscountRepository persistence.ProductRepository) RetrieveCheckoutAmount {
	return RetrieveCheckoutAmount{checkoutRepository, productRepository, productWithPromotionRepository, productWithDiscountRepository}
}

func (service *RetrieveCheckoutAmount) Do(checkoutId string) (int, error) {
	checkout, existCheckout := service.CheckoutRepository.SearchById(checkoutId)
	if !existCheckout {
		return 0, errors.NewCheckoutNotFoundError()
	}

	checkoutAmount := calculateCheckoutAmount(checkout.Products, service.ProductRepository, service.ProductWithPromotionRepository, service.ProductWithDiscountRepository)

	return checkoutAmount, nil
}

func calculateCheckoutAmount(checkoutProducts []string, productsRepository persistence.ProductRepository, productsWithPromotionRepository persistence.ProductRepository, productsWithDiscountRepository persistence.ProductRepository) int {
	productRealUnits := calculateRealProductUnits(checkoutProducts)
	productUnits := calculatePayableProductUnits(productRealUnits, productsWithPromotionRepository)

	var amount int
	for productCode, quantity := range productUnits {
		product, _ := productsRepository.SearchById(productCode)
		if _, hasDiscount := productsWithDiscountRepository.SearchById(productCode); hasDiscount {
			amount += calculateAmountWithDiscount(quantity, product.Price)
			continue
		}
		amount += (product.Price * quantity)
	}
	return amount
}

func calculateRealProductUnits(checkoutProducts []string) map[string]int {
	var productUnits = map[string]int{"PEN": 0, "TSHIRT": 0, "MUG": 0}
	for _, checkoutProductCode := range checkoutProducts {
		productUnits[checkoutProductCode] += 1
	}
	return productUnits
}

func calculatePayableProductUnits(productRealUnits map[string]int, productsWithPromotionRepository persistence.ProductRepository) map[string]int {
	var productUnits = map[string]int{"PEN": 0, "TSHIRT": 0, "MUG": 0}
	for productCode, quantity := range productRealUnits {
		if _, found := productsWithPromotionRepository.SearchById(productCode); found {
			productUnits[productCode] = calculatePayableUnitsApplying2X1Promotion(quantity)
			continue
		}
		productUnits[productCode] = quantity
	}
	return productUnits
}

func calculatePayableUnitsApplying2X1Promotion(quantity int) int {
	if quantity == 0 {
		return 0
	}
	if quantity%2 == 0 {
		return quantity / 2
	}
	return ((quantity - 1) / 2) + 1
}

func calculateAmountWithDiscount(quantity int, price int) int {
	if quantity < 3 {
		return price * quantity
	}
	unitPriceWithDiscount := (price * 75) / 100
	return unitPriceWithDiscount * quantity
}
