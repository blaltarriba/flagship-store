package persistence

import (
	"lana/flagship-store/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSearchByIdReturnCheckoutWhenCheckoutExists(t *testing.T) {
	checkout_id := uuid.NewString()
	checkout := models.Checkout{
		Id:       checkout_id,
		Products: []string{"PEN"},
	}
	checkouts := map[string]models.Checkout{checkout.Id: checkout}

	inMemoryCheckoutRepository := &InMemoryCheckoutRepository{checkouts}

	checkoutRetrieved, exists := inMemoryCheckoutRepository.SearchById(checkout_id)

	assert.EqualValues(t, true, exists)
	assert.EqualValues(t, checkout, checkoutRetrieved)
}

func TestSearchByIdReturnEmptyCheckoutWhenCheckoutDoesNotExist(t *testing.T) {
	checkouts := make(map[string]models.Checkout)
	inMemoryCheckoutRepository := &InMemoryCheckoutRepository{checkouts}

	checkoutRetrieved, exists := inMemoryCheckoutRepository.SearchById("an_id")

	assert.EqualValues(t, false, exists)
	assert.EqualValues(t, "", checkoutRetrieved.Id)
}

func TestPersistCreateCheckoutWhenCheckoutDoesNotExist(t *testing.T) {
	checkout_id := uuid.NewString()
	checkout := models.Checkout{
		Id:       checkout_id,
		Products: []string{"PEN"},
	}
	checkouts := make(map[string]models.Checkout)
	inMemoryCheckoutRepository := &InMemoryCheckoutRepository{checkouts}

	inMemoryCheckoutRepository.Persist(checkout)

	assert.EqualValues(t, 1, len(checkouts))
	assert.EqualValues(t, checkout, checkouts[checkout_id])
}

func TestPersistUpdateCheckoutWhenCheckoutExist(t *testing.T) {
	checkout_id := uuid.NewString()
	checkout := models.Checkout{
		Id:       checkout_id,
		Products: []string{"PEN"},
	}
	checkouts := map[string]models.Checkout{checkout.Id: checkout}
	inMemoryCheckoutRepository := &InMemoryCheckoutRepository{checkouts}
	checkout.Products = append(checkout.Products, "MUG")

	inMemoryCheckoutRepository.Persist(checkout)

	modifiedCheckout := checkouts[checkout_id]
	assert.EqualValues(t, 1, len(checkouts))
	assert.EqualValues(t, 2, len(modifiedCheckout.Products))
	assert.EqualValues(t, "MUG", modifiedCheckout.Products[1])
}

func TestDeleteRemoveCheckoutWhenCheckoutExist(t *testing.T) {
	checkout_id := uuid.NewString()
	checkout := models.Checkout{
		Id:       checkout_id,
		Products: []string{"PEN"},
	}
	checkouts := map[string]models.Checkout{checkout.Id: checkout}
	inMemoryCheckoutRepository := &InMemoryCheckoutRepository{checkouts}

	inMemoryCheckoutRepository.Delete(checkout)

	assert.EqualValues(t, 0, len(checkouts))
}

func TestDeleteDoesNothingWhenCheckoutDoesNotExist(t *testing.T) {
	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"PEN"},
	}
	checkouts := make(map[string]models.Checkout)
	inMemoryCheckoutRepository := &InMemoryCheckoutRepository{checkouts}

	inMemoryCheckoutRepository.Delete(checkout)

	assert.EqualValues(t, 0, len(checkouts))
}

func TestCountReturnNumberOfCheckoutsWhenCheckoutExist(t *testing.T) {
	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"PEN"},
	}
	checkouts := map[string]models.Checkout{checkout.Id: checkout}
	inMemoryCheckoutRepository := &InMemoryCheckoutRepository{checkouts}

	count := inMemoryCheckoutRepository.Count()

	assert.EqualValues(t, 1, count)
}

func TestCountReturnZeroWhenCheckoutDoesNotExist(t *testing.T) {
	checkouts := make(map[string]models.Checkout)
	inMemoryCheckoutRepository := &InMemoryCheckoutRepository{checkouts}

	count := inMemoryCheckoutRepository.Count()

	assert.EqualValues(t, 0, count)
}
