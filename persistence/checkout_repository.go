package persistence

import "lana/flagship-store/models"

type CheckoutRepository interface {
	SearchById(id string) (models.Checkout, bool)
	Persist(checkout models.Checkout)
	Delete(checkout models.Checkout)
	Count() int
}
