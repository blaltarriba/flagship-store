package persistence

import "lana/flagship-store/models"

type InMemoryCheckoutRepository struct {
	checkouts map[string]models.Checkout
}

func NewCheckoutRepository(checkouts map[string]models.Checkout) *InMemoryCheckoutRepository {
	return &InMemoryCheckoutRepository{checkouts}
}

func (repository *InMemoryCheckoutRepository) SearchById(id string) (models.Checkout, bool) {
	checkout, exists := repository.checkouts[id]
	return checkout, exists
}

func (repository *InMemoryCheckoutRepository) Persist(checkout models.Checkout) {
	repository.checkouts[checkout.Id] = checkout
}

func (repository *InMemoryCheckoutRepository) Delete(checkout models.Checkout) {
	delete(repository.checkouts, checkout.Id)
}

func (repository *InMemoryCheckoutRepository) Count() int {
	return len(repository.checkouts)
}
