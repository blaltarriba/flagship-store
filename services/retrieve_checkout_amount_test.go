package services

import (
	"lana/flagship-store/models"
	"lana/flagship-store/utils/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func ProductRepositoryMockWithAllProducts() mocks.ProductRepositoryMock {
	pen := models.Product{
		Code:  "PEN",
		Name:  "Lana Pen",
		Price: 500,
	}
	tshirt := models.Product{
		Code:  "TSHIRT",
		Name:  "Lana T-Shirt",
		Price: 2000,
	}
	mug := models.Product{
		Code:  "MUG",
		Name:  "Lana Coffee Mug",
		Price: 750,
	}
	theProductRepositoryMock := mocks.ProductRepositoryMock{}
	theProductRepositoryMock.On("SearchById", pen.Code).Return(pen, true)
	theProductRepositoryMock.On("SearchById", mug.Code).Return(mug, true)
	theProductRepositoryMock.On("SearchById", tshirt.Code).Return(tshirt, true)

	return theProductRepositoryMock
}

func ProductWithPromotionRepositoryMockWithProducts() mocks.ProductWithPromotionRepositoryMock {
	pen := models.Product{
		Code:  "PEN",
		Name:  "Lana Pen",
		Price: 500,
	}
	theProductWithPromotionRepositoryMock := mocks.ProductWithPromotionRepositoryMock{}
	theProductWithPromotionRepositoryMock.On("SearchById", pen.Code).Return(pen, true)
	theProductWithPromotionRepositoryMock.On("SearchById", "MUG").Return(models.Product{}, false)
	theProductWithPromotionRepositoryMock.On("SearchById", "TSHIRT").Return(models.Product{}, false)

	return theProductWithPromotionRepositoryMock
}

func ProductWithDiscountRepositoryMockWithProducts() mocks.ProductWithDiscountRepositoryMock {
	tshirt := models.Product{
		Code:  "TSHIRT",
		Name:  "Lana T-Shirt",
		Price: 2000,
	}
	theProductWithDiscountRepositoryMock := mocks.ProductWithDiscountRepositoryMock{}
	theProductWithDiscountRepositoryMock.On("SearchById", tshirt.Code).Return(tshirt, true)
	theProductWithDiscountRepositoryMock.On("SearchById", "MUG").Return(models.Product{}, false)
	theProductWithDiscountRepositoryMock.On("SearchById", "PEN").Return(models.Product{}, false)

	return theProductWithDiscountRepositoryMock
}

func TestRetrieveCheckoutAmountWhenCheckoutExists(t *testing.T) {
	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG"},
	}
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", checkout.Id).Return(checkout, true)
	theProductRepositoryMock := ProductRepositoryMockWithAllProducts()
	theProductWithDiscountRepositoryMock := ProductWithDiscountRepositoryMockWithProducts()
	theProductWithPromotionRepositoryMock := ProductWithPromotionRepositoryMockWithProducts()
	retrieveCheckoutAmountService := RetrieveCheckoutAmount{
		&theCheckoutRepositoryMock,
		&theProductRepositoryMock,
		&theProductWithPromotionRepositoryMock,
		&theProductWithDiscountRepositoryMock}

	checkoutAmount, _ := retrieveCheckoutAmountService.Do(checkout.Id)

	assert.EqualValues(t, 750, checkoutAmount)
}

func TestAmountWith2X1PromotionWhenCheckoutContainsTwoOfSameProductWithPromotion(t *testing.T) {
	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"PEN", "PEN"},
	}
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", checkout.Id).Return(checkout, true)
	theProductRepositoryMock := ProductRepositoryMockWithAllProducts()
	theProductWithDiscountRepositoryMock := ProductWithDiscountRepositoryMockWithProducts()
	theProductWithPromotionRepositoryMock := ProductWithPromotionRepositoryMockWithProducts()
	retrieveCheckoutAmountService := RetrieveCheckoutAmount{
		&theCheckoutRepositoryMock,
		&theProductRepositoryMock,
		&theProductWithPromotionRepositoryMock,
		&theProductWithDiscountRepositoryMock}

	checkoutAmount, _ := retrieveCheckoutAmountService.Do(checkout.Id)

	assert.EqualValues(t, 500, checkoutAmount)
}

func TestAmountWithNo2X1PromotionWhenCheckoutDoesNotContainsTwoOfSameProductWithPromotion(t *testing.T) {
	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG", "MUG"},
	}
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", checkout.Id).Return(checkout, true)
	theProductRepositoryMock := ProductRepositoryMockWithAllProducts()
	theProductWithDiscountRepositoryMock := ProductWithDiscountRepositoryMockWithProducts()
	theProductWithPromotionRepositoryMock := ProductWithPromotionRepositoryMockWithProducts()
	retrieveCheckoutAmountService := RetrieveCheckoutAmount{
		&theCheckoutRepositoryMock,
		&theProductRepositoryMock,
		&theProductWithPromotionRepositoryMock,
		&theProductWithDiscountRepositoryMock}

	checkoutAmount, _ := retrieveCheckoutAmountService.Do(checkout.Id)

	assert.EqualValues(t, 1500, checkoutAmount)
}

func TestAmountWithDiscountWhenCheckoutContainsThreeOfSameProductWithDiscount(t *testing.T) {
	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"TSHIRT", "TSHIRT", "TSHIRT"},
	}
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", checkout.Id).Return(checkout, true)
	theProductRepositoryMock := ProductRepositoryMockWithAllProducts()
	theProductWithDiscountRepositoryMock := ProductWithDiscountRepositoryMockWithProducts()
	theProductWithPromotionRepositoryMock := ProductWithPromotionRepositoryMockWithProducts()
	retrieveCheckoutAmountService := RetrieveCheckoutAmount{
		&theCheckoutRepositoryMock,
		&theProductRepositoryMock,
		&theProductWithPromotionRepositoryMock,
		&theProductWithDiscountRepositoryMock}

	checkoutAmount, _ := retrieveCheckoutAmountService.Do(checkout.Id)

	assert.EqualValues(t, 4500, checkoutAmount)
}

func TestAmountWithNoDiscountWhenCheckoutContainsLessThanThreeOfSameProductWithDiscount(t *testing.T) {
	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"TSHIRT", "TSHIRT"},
	}
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", checkout.Id).Return(checkout, true)
	theProductRepositoryMock := ProductRepositoryMockWithAllProducts()
	theProductWithDiscountRepositoryMock := ProductWithDiscountRepositoryMockWithProducts()
	theProductWithPromotionRepositoryMock := ProductWithPromotionRepositoryMockWithProducts()
	retrieveCheckoutAmountService := RetrieveCheckoutAmount{
		&theCheckoutRepositoryMock,
		&theProductRepositoryMock,
		&theProductWithPromotionRepositoryMock,
		&theProductWithDiscountRepositoryMock}

	checkoutAmount, _ := retrieveCheckoutAmountService.Do(checkout.Id)

	assert.EqualValues(t, 4000, checkoutAmount)
}

func TestAmountWithNoDiscountWhenCheckoutDoesNotContainsThreeOfSameProductWithDiscount(t *testing.T) {
	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG", "MUG", "MUG"},
	}
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", checkout.Id).Return(checkout, true)
	theProductRepositoryMock := ProductRepositoryMockWithAllProducts()
	theProductWithDiscountRepositoryMock := ProductWithDiscountRepositoryMockWithProducts()
	theProductWithPromotionRepositoryMock := ProductWithPromotionRepositoryMockWithProducts()
	retrieveCheckoutAmountService := RetrieveCheckoutAmount{
		&theCheckoutRepositoryMock,
		&theProductRepositoryMock,
		&theProductWithPromotionRepositoryMock,
		&theProductWithDiscountRepositoryMock}

	checkoutAmount, _ := retrieveCheckoutAmountService.Do(checkout.Id)

	assert.EqualValues(t, 2250, checkoutAmount)
}
