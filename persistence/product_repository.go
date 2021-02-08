package persistence

import "lana/flagship-store/models"

type ProductRepository interface {
	SearchById(id string) (models.Product, bool)
}
