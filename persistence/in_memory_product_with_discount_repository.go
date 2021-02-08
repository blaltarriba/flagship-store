package persistence

import "lana/flagship-store/models"

type InMemoryProductWithDiscountRepository struct {
	products map[string]models.Product
}

func NewProductWithDiscountRepository(products map[string]models.Product) *InMemoryProductWithDiscountRepository {
	return &InMemoryProductWithDiscountRepository{products}
}

func (repository *InMemoryProductWithDiscountRepository) SearchById(id string) (models.Product, bool) {
	product, exists := repository.products[id]
	return product, exists
}
