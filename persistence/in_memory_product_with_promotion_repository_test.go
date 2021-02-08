package persistence

import (
	"lana/flagship-store/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchByIdReturnProductWhenProductHasPromotion(t *testing.T) {
	pen := models.Product{
		Code:  "PEN",
		Name:  "Lana Pen",
		Price: 500,
	}
	products := make(map[string]models.Product)
	products[pen.Code] = pen
	inMemoryProductWithPromotionRepository := &InMemoryProductWithPromotionRepository{products}

	product, exists := inMemoryProductWithPromotionRepository.SearchById("PEN")

	assert.EqualValues(t, true, exists)
	assert.EqualValues(t, pen, product)
}

func TestSearchByIdReturnEmptyProductWhenProductDoesNotHavePromotion(t *testing.T) {
	products := make(map[string]models.Product)
	inMemoryProductWithPromotionRepository := &InMemoryProductWithPromotionRepository{products}

	product, exists := inMemoryProductWithPromotionRepository.SearchById("PEN")

	assert.EqualValues(t, false, exists)
	assert.EqualValues(t, "", product.Code)
	assert.EqualValues(t, "", product.Name)
	assert.EqualValues(t, 0, product.Price)
}
