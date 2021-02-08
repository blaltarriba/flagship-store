package persistence

import "lana/flagship-store/models"

type InMemoryProductWithPromotionRepository struct {
	products map[string]models.Product
}

func NewProductWithPromotionRepository(products map[string]models.Product) *InMemoryProductWithPromotionRepository {
	return &InMemoryProductWithPromotionRepository{products}
}

func (repository *InMemoryProductWithPromotionRepository) SearchById(id string) (models.Product, bool) {
	product, exists := repository.products[id]
	return product, exists
}
