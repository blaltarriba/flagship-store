package services

import (
	"lana/flagship-store/models"
	"lana/flagship-store/persistence"
	"lana/flagship-store/services/commands"
	"lana/flagship-store/services/errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var createCheckout CreateCheckout

func TestMain(m *testing.M) {
	checkoutRepository := initialize_checkout_repository()
	productRepository := initialize_product_repository()

	createCheckout = NewCreateCheckout(checkoutRepository, productRepository)

	code := m.Run()

	clearCheckouts()

	os.Exit(code)
}

func clearCheckouts() {
	checkouts := make(map[string]models.Checkout)
	createCheckout.CheckoutRepository = persistence.NewCheckoutRepository(checkouts)
}

func initialize_checkout_repository() persistence.CheckoutRepository {
	checkouts := make(map[string]models.Checkout)
	return persistence.NewCheckoutRepository(checkouts)
}

func initialize_product_repository() persistence.ProductRepository {
	products := make(map[string]models.Product)

	pen := models.Product{
		Code:  "PEN",
		Name:  "Lana Pen",
		Price: 500,
	}

	products[pen.Code] = pen
	return persistence.NewProductsRepository(products)
}

func TestCreateCheckout(t *testing.T) {
	clearCheckouts()

	productCommand := commands.Product{Code: "PEN"}

	createdCheckout, _ := createCheckout.Do(productCommand)

	assert.NotNil(t, createdCheckout.Id)
	assert.EqualValues(t, "PEN", createdCheckout.Products[0])
	assert.EqualValues(t, 1, len(createdCheckout.Products))
}

func TestReturnProductNotFoundErrorWhenProductDoesnotExists(t *testing.T) {
	clearCheckouts()

	productCommand := commands.Product{Code: "FAKE"}

	_, err := createCheckout.Do(productCommand)

	_, isProductNotFoundError := err.(*errors.ProductNotFoundError)
	assert.EqualValues(t, true, isProductNotFoundError)
}
