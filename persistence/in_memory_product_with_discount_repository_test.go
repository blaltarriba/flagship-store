package persistence

import (
	"lana/flagship-store/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchByIdReturnProductWhenProductHasDiscount(t *testing.T) {
	pen := models.Product{
		Code:  "PEN",
		Name:  "Lana Pen",
		Price: 500,
	}
	products := make(map[string]models.Product)
	products[pen.Code] = pen
	inMemoryProductWithDiscountRepository := &InMemoryProductWithDiscountRepository{products}

	product, exists := inMemoryProductWithDiscountRepository.SearchById("PEN")

	assert.EqualValues(t, true, exists)
	assert.EqualValues(t, pen, product)
}

func TestSearchByIdReturnEmptyProductWhenProductDoesNotHaveDiscount(t *testing.T) {
	products := make(map[string]models.Product)
	inMemoryProductWithDiscountRepository := &InMemoryProductWithDiscountRepository{products}

	product, exists := inMemoryProductWithDiscountRepository.SearchById("PEN")

	assert.EqualValues(t, false, exists)
	assert.EqualValues(t, "", product.Code)
	assert.EqualValues(t, "", product.Name)
	assert.EqualValues(t, 0, product.Price)
}
