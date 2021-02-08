package persistence

import "lana/flagship-store/models"

type InMemoryProductsRepository struct {
	products map[string]models.Product
}

func NewProductsRepository(products map[string]models.Product) *InMemoryProductsRepository {
	return &InMemoryProductsRepository{products}
}

func (repository *InMemoryProductsRepository) SearchById(id string) (models.Product, bool) {
	product, exists := repository.products[id]
	return product, exists
}
