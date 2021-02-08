package main

import (
	"bytes"
	"encoding/json"
	"lana/flagship-store/models"
	"lana/flagship-store/persistence"
	"lana/flagship-store/services/responses"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var app App

func TestMain(m *testing.M) {
	checkoutRepository := initialize_checkout_repository()
	productRepository := initialize_product_repository()
	productsWithPromotionRepository := initialize_products_with_promotions()
	productsWithDiscount := initialize_products_with_discount()

	app = App{}
	app.Initialize(checkoutRepository, productRepository, productsWithPromotionRepository, productsWithDiscount)

	code := m.Run()

	clearCheckouts()

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}

func clearCheckouts() {
	checkouts := make(map[string]models.Checkout)
	app.CheckoutRepository = persistence.NewCheckoutRepository(checkouts)
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
	tshirt := models.Product{
		Code:  "TSHIRT",
		Name:  "Lana T-Shirt",
		Price: 2000,
	}
	mug := models.Product{
		Code:  "MUG",
		Name:  "Lana Coffee Mug",
		Price: 750,
	}

	products[pen.Code] = pen
	products[tshirt.Code] = tshirt
	products[mug.Code] = mug
	return persistence.NewProductsRepository(products)
}

func initialize_products_with_promotions() persistence.ProductRepository {
	products := make(map[string]models.Product)

	pen := models.Product{
		Code:  "PEN",
		Name:  "Lana Pen",
		Price: 500,
	}

	products[pen.Code] = pen
	return persistence.NewProductWithPromotionRepository(products)
}

func initialize_products_with_discount() map[string]models.Product {
	products := make(map[string]models.Product)

	tshirt := models.Product{
		Code:  "TSHIRT",
		Name:  "Lana T-Shirt",
		Price: 2000,
	}
	products[tshirt.Code] = tshirt
	return products
}

func TestReturn200WhenCreateCheckout(t *testing.T) {
	clearCheckouts()

	payload := []byte(`{"product-code":"PEN"}`)

	req, _ := http.NewRequest("POST", "/checkouts", bytes.NewBuffer(payload))
	response := executeRequest(req)

	assert.EqualValues(t, 201, response.Code)
}

func TestCreateCheckout(t *testing.T) {
	clearCheckouts()

	payload := []byte(`{"product-code":"PEN"}`)

	req, _ := http.NewRequest("POST", "/checkouts", bytes.NewBuffer(payload))
	response := executeRequest(req)

	var createdCheckout models.Checkout
	json.Unmarshal(response.Body.Bytes(), &createdCheckout)

	assert.EqualValues(t, 1, app.CheckoutRepository.Count())
	assert.NotNil(t, createdCheckout.Id)
	assert.EqualValues(t, "PEN", createdCheckout.Products[0])
	assert.EqualValues(t, 1, len(createdCheckout.Products))
}

func TestReturn404WhenCreateCheckoutWithNotValidProduct(t *testing.T) {
	clearCheckouts()

	payload := []byte(`{"product-code":"FAKE"}`)

	req, _ := http.NewRequest("POST", "/checkouts", bytes.NewBuffer(payload))
	response := executeRequest(req)

	var productNotFound responses.ProductNotFound
	json.Unmarshal(response.Body.Bytes(), &productNotFound)

	assert.EqualValues(t, 404, response.Code)
	assert.EqualValues(t, "Product FAKE not found", productNotFound.Message)
}

func TestReturn204AddingProductToCheckoutWhenCheckoutExists(t *testing.T) {
	clearCheckouts()

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG"},
	}
	app.CheckoutRepository.Persist(checkout)

	payload := []byte(`{"product":"PEN"}`)

	req, _ := http.NewRequest("PATCH", "/checkouts/"+checkout.Id, bytes.NewBuffer(payload))
	response := executeRequest(req)

	assert.EqualValues(t, 204, response.Code)
}

func TestAddProductToCheckoutWhenCheckoutExists(t *testing.T) {
	clearCheckouts()

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG"},
	}
	app.CheckoutRepository.Persist(checkout)

	payload := []byte(`{"product":"PEN"}`)

	req, _ := http.NewRequest("PATCH", "/checkouts/"+checkout.Id, bytes.NewBuffer(payload))
	executeRequest(req)

	modifiedCheckout, _ := app.CheckoutRepository.SearchById(checkout.Id)
	assert.EqualValues(t, 2, len(modifiedCheckout.Products))
	assert.EqualValues(t, "PEN", modifiedCheckout.Products[1])
}

func TestReturn404AddingProductToCheckoutWhenCheckoutDoesNotExists(t *testing.T) {
	clearCheckouts()

	payload := []byte(`{"product":"PEN"}`)

	req, _ := http.NewRequest("PATCH", "/checkouts/a_fake_checkout", bytes.NewBuffer(payload))
	response := executeRequest(req)

	var checkoutNotFound responses.CheckoutNotFound
	json.Unmarshal(response.Body.Bytes(), &checkoutNotFound)

	assert.EqualValues(t, 404, response.Code)
	assert.EqualValues(t, "Checkout a_fake_checkout not found", checkoutNotFound.Message)
}

func TestReturn422AddingProductToCheckoutWhenProductDoesNotExists(t *testing.T) {
	clearCheckouts()

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG"},
	}
	app.CheckoutRepository.Persist(checkout)

	payload := []byte(`{"product":"FAKE"}`)

	req, _ := http.NewRequest("PATCH", "/checkouts/"+checkout.Id, bytes.NewBuffer(payload))
	response := executeRequest(req)

	var productNotFound responses.ProductNotFound
	json.Unmarshal(response.Body.Bytes(), &productNotFound)

	assert.EqualValues(t, 422, response.Code)
	assert.EqualValues(t, "Product FAKE not found", productNotFound.Message)
}

func TestReturn200RetrievingCheckoutAmountWhenCheckoutExists(t *testing.T) {
	clearCheckouts()

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"PEN"},
	}
	app.CheckoutRepository.Persist(checkout)

	req, _ := http.NewRequest("GET", "/checkouts/"+checkout.Id+"/amount", nil)
	response := executeRequest(req)

	assert.EqualValues(t, 200, response.Code)
}

func TestAmountWhenCheckoutExists(t *testing.T) {
	clearCheckouts()

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG"},
	}
	app.CheckoutRepository.Persist(checkout)

	req, _ := http.NewRequest("GET", "/checkouts/"+checkout.Id+"/amount", nil)
	response := executeRequest(req)

	var responseCheckout responses.Checkout
	json.Unmarshal(response.Body.Bytes(), &responseCheckout)

	assert.EqualValues(t, "7.50€", responseCheckout.Amount)
}

func TestReturn404RetrievingCheckoutAmountWhenCheckoutDoesNotExists(t *testing.T) {
	clearCheckouts()

	req, _ := http.NewRequest("GET", "/checkouts/a_fake_checkout/amount", nil)
	response := executeRequest(req)

	var checkoutNotFound responses.CheckoutNotFound
	json.Unmarshal(response.Body.Bytes(), &checkoutNotFound)

	assert.EqualValues(t, 404, response.Code)
	assert.EqualValues(t, "Checkout a_fake_checkout not found", checkoutNotFound.Message)
}

func TestAmountWith2X1PromotionWhenCheckoutContainsTwoOfSameProductWithPromotion(t *testing.T) {
	clearCheckouts()

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"PEN", "PEN"},
	}
	app.CheckoutRepository.Persist(checkout)

	req, _ := http.NewRequest("GET", "/checkouts/"+checkout.Id+"/amount", nil)
	response := executeRequest(req)

	var responseCheckout responses.Checkout
	json.Unmarshal(response.Body.Bytes(), &responseCheckout)

	assert.EqualValues(t, "5.00€", responseCheckout.Amount)
}

func TestAmountWithNo2X1PromotionWhenCheckoutDoesNotContainsTwoOfSameProductWithPromotion(t *testing.T) {
	clearCheckouts()

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG", "MUG"},
	}
	app.CheckoutRepository.Persist(checkout)

	req, _ := http.NewRequest("GET", "/checkouts/"+checkout.Id+"/amount", nil)
	response := executeRequest(req)

	var responseCheckout responses.Checkout
	json.Unmarshal(response.Body.Bytes(), &responseCheckout)

	assert.EqualValues(t, "15.00€", responseCheckout.Amount)
}

func TestAmountWithDiscountWhenCheckoutContainsThreeOfSameProductWithDiscount(t *testing.T) {
	clearCheckouts()

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"TSHIRT", "TSHIRT", "TSHIRT"},
	}
	app.CheckoutRepository.Persist(checkout)

	req, _ := http.NewRequest("GET", "/checkouts/"+checkout.Id+"/amount", nil)
	response := executeRequest(req)

	var responseCheckout responses.Checkout
	json.Unmarshal(response.Body.Bytes(), &responseCheckout)

	assert.EqualValues(t, "45.00€", responseCheckout.Amount)
}

func TestAmountWithNoDiscountWhenCheckoutContainsLessThanThreeOfSameProductWithDiscount(t *testing.T) {
	clearCheckouts()

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"TSHIRT", "TSHIRT"},
	}
	app.CheckoutRepository.Persist(checkout)

	req, _ := http.NewRequest("GET", "/checkouts/"+checkout.Id+"/amount", nil)
	response := executeRequest(req)

	var responseCheckout responses.Checkout
	json.Unmarshal(response.Body.Bytes(), &responseCheckout)

	assert.EqualValues(t, "40.00€", responseCheckout.Amount)
}

func TestAmountWithNoDiscountWhenCheckoutDoesNotContainsThreeOfSameProductWithDiscount(t *testing.T) {
	clearCheckouts()

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG", "MUG", "MUG"},
	}
	app.CheckoutRepository.Persist(checkout)

	req, _ := http.NewRequest("GET", "/checkouts/"+checkout.Id+"/amount", nil)
	response := executeRequest(req)

	var responseCheckout responses.Checkout
	json.Unmarshal(response.Body.Bytes(), &responseCheckout)

	assert.EqualValues(t, "22.50€", responseCheckout.Amount)
}

func TestReturn204WhenDeleteCheckout(t *testing.T) {
	clearCheckouts()

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"PEN"},
	}
	app.CheckoutRepository.Persist(checkout)

	req, _ := http.NewRequest("DELETE", "/checkouts/"+checkout.Id, nil)
	response := executeRequest(req)

	assert.EqualValues(t, 204, response.Code)
}

func TestDeleteCheckout(t *testing.T) {
	clearCheckouts()

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"PEN"},
	}
	app.CheckoutRepository.Persist(checkout)

	req, _ := http.NewRequest("DELETE", "/checkouts/"+checkout.Id, nil)
	executeRequest(req)

	assert.EqualValues(t, 0, app.CheckoutRepository.Count())
}

func TestReturn404DeletingCheckoutWhenCheckoutDoesNotExists(t *testing.T) {
	clearCheckouts()

	req, _ := http.NewRequest("DELETE", "/checkouts/a_fake_checkout", nil)
	response := executeRequest(req)

	var checkoutNotFound responses.CheckoutNotFound
	json.Unmarshal(response.Body.Bytes(), &checkoutNotFound)

	assert.EqualValues(t, 404, response.Code)
	assert.EqualValues(t, "Checkout a_fake_checkout not found", checkoutNotFound.Message)
}
