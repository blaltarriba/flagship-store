package main

import (
	"bytes"
	"encoding/json"
	"lana/flagship-store/models"
	"lana/flagship-store/services"
	"lana/flagship-store/services/responses"
	"lana/flagship-store/utils/mocks"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var app App

func TestMain(m *testing.M) {
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theProductRepositoryMock := mocks.ProductRepositoryMock{}
	theProductWithDiscountRepositoryMock := mocks.ProductWithDiscountRepositoryMock{}
	theProductWithPromotionRepositoryMock := mocks.ProductWithPromotionRepositoryMock{}
	createCheckoutService := services.NewCreateCheckout(&theCheckoutRepositoryMock, &theProductRepositoryMock)
	addProductToCheckoutService := services.NewAddProductToCheckout(&theCheckoutRepositoryMock, &theProductRepositoryMock)
	retrieveCheckoutAmountService := services.NewRetrieveCheckoutAmount(&theCheckoutRepositoryMock, &theProductRepositoryMock, &theProductWithDiscountRepositoryMock, &theProductWithPromotionRepositoryMock)
	deleteCheckoutService := services.NewDeleteCheckout(&theCheckoutRepositoryMock)

	app = App{}
	app.Initialize(createCheckoutService, addProductToCheckoutService, deleteCheckoutService, retrieveCheckoutAmountService)

	code := m.Run()

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}

func ACheckout() models.Checkout {
	return models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG"},
	}
}

func ProductRepositoryMockWithAllProducts() mocks.ProductRepositoryMock {
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
	theProductRepositoryMock := mocks.ProductRepositoryMock{}
	theProductRepositoryMock.On("SearchById", pen.Code).Return(pen, true)
	theProductRepositoryMock.On("SearchById", mug.Code).Return(mug, true)
	theProductRepositoryMock.On("SearchById", tshirt.Code).Return(tshirt, true)

	return theProductRepositoryMock
}

func TestReturn200WhenCreateCheckout(t *testing.T) {
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("Persist", mock.AnythingOfType("models.Checkout"))
	theProductRepositoryMock := mocks.ProductRepositoryMock{}
	theProductRepositoryMock.On("SearchById", "PEN").Return(models.Product{}, true)
	app.CreateCheckoutService = services.NewCreateCheckout(&theCheckoutRepositoryMock, &theProductRepositoryMock)
	payload := []byte(`{"product-code":"PEN"}`)

	req, _ := http.NewRequest("POST", "/checkouts", bytes.NewBuffer(payload))
	response := executeRequest(req)

	var createdCheckout models.Checkout
	json.Unmarshal(response.Body.Bytes(), &createdCheckout)
	assert.EqualValues(t, 201, response.Code)
	assert.NotNil(t, createdCheckout.Id)
	assert.EqualValues(t, "PEN", createdCheckout.Products[0])
	assert.EqualValues(t, 1, len(createdCheckout.Products))
	theCheckoutRepositoryMock.AssertExpectations(t)
	theProductRepositoryMock.AssertExpectations(t)
}

func TestReturn404WhenCreateCheckoutWithNotValidProduct(t *testing.T) {
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theProductRepositoryMock := mocks.ProductRepositoryMock{}
	theProductRepositoryMock.On("SearchById", "FAKE").Return(models.Product{}, false)
	app.CreateCheckoutService = services.NewCreateCheckout(&theCheckoutRepositoryMock, &theProductRepositoryMock)
	payload := []byte(`{"product-code":"FAKE"}`)

	req, _ := http.NewRequest("POST", "/checkouts", bytes.NewBuffer(payload))
	response := executeRequest(req)

	var productNotFound responses.ProductNotFound
	json.Unmarshal(response.Body.Bytes(), &productNotFound)
	assert.EqualValues(t, 404, response.Code)
	assert.EqualValues(t, "Product FAKE not found", productNotFound.Message)
	theCheckoutRepositoryMock.AssertExpectations(t)
	theProductRepositoryMock.AssertExpectations(t)
}

func TestReturn204AddingProductToCheckoutWhenCheckoutExists(t *testing.T) {
	checkout := ACheckout()
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", checkout.Id).Return(checkout, true)
	theCheckoutRepositoryMock.On("Persist", mock.AnythingOfType("models.Checkout"))
	theProductRepositoryMock := mocks.ProductRepositoryMock{}
	theProductRepositoryMock.On("SearchById", "PEN").Return(models.Product{}, true)
	app.AddProductToCheckoutService = services.NewAddProductToCheckout(&theCheckoutRepositoryMock, &theProductRepositoryMock)
	payload := []byte(`{"product":"PEN"}`)

	req, _ := http.NewRequest("PATCH", "/checkouts/"+checkout.Id, bytes.NewBuffer(payload))
	response := executeRequest(req)
	assert.EqualValues(t, 204, response.Code)
	theCheckoutRepositoryMock.AssertExpectations(t)
	theProductRepositoryMock.AssertExpectations(t)
}

func TestReturn404AddingProductToCheckoutWhenCheckoutDoesNotExists(t *testing.T) {
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", mock.AnythingOfType("string")).Return(models.Checkout{}, false)
	theProductRepositoryMock := mocks.ProductRepositoryMock{}
	app.AddProductToCheckoutService = services.NewAddProductToCheckout(&theCheckoutRepositoryMock, &theProductRepositoryMock)
	payload := []byte(`{"product":"PEN"}`)

	req, _ := http.NewRequest("PATCH", "/checkouts/a_fake_checkout", bytes.NewBuffer(payload))
	response := executeRequest(req)

	var checkoutNotFound responses.CheckoutNotFound
	json.Unmarshal(response.Body.Bytes(), &checkoutNotFound)
	assert.EqualValues(t, 404, response.Code)
	assert.EqualValues(t, "Checkout a_fake_checkout not found", checkoutNotFound.Message)
	theCheckoutRepositoryMock.AssertExpectations(t)
	theProductRepositoryMock.AssertExpectations(t)
}

func TestReturn422AddingProductToCheckoutWhenProductDoesNotExists(t *testing.T) {
	checkout := ACheckout()
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", checkout.Id).Return(checkout, true)
	theProductRepositoryMock := mocks.ProductRepositoryMock{}
	theProductRepositoryMock.On("SearchById", "FAKE").Return(models.Product{}, false)
	app.AddProductToCheckoutService = services.NewAddProductToCheckout(&theCheckoutRepositoryMock, &theProductRepositoryMock)
	payload := []byte(`{"product":"FAKE"}`)

	req, _ := http.NewRequest("PATCH", "/checkouts/"+checkout.Id, bytes.NewBuffer(payload))
	response := executeRequest(req)

	var productNotFound responses.ProductNotFound
	json.Unmarshal(response.Body.Bytes(), &productNotFound)
	assert.EqualValues(t, 422, response.Code)
	assert.EqualValues(t, "Product FAKE not found", productNotFound.Message)
	theCheckoutRepositoryMock.AssertExpectations(t)
	theProductRepositoryMock.AssertExpectations(t)
}

func TestReturn200RetrievingCheckoutAmountWhenCheckoutExists(t *testing.T) {
	checkout := ACheckout()
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", checkout.Id).Return(checkout, true)
	theProductRepositoryMock := ProductRepositoryMockWithAllProducts()
	theProductWithDiscountRepositoryMock := mocks.ProductWithDiscountRepositoryMock{}
	theProductWithDiscountRepositoryMock.On("SearchById", mock.AnythingOfType("string")).Return(models.Product{}, false)
	theProductWithPromotionRepositoryMock := mocks.ProductWithPromotionRepositoryMock{}
	theProductWithPromotionRepositoryMock.On("SearchById", mock.AnythingOfType("string")).Return(models.Product{}, false)
	app.RetrieveCheckoutAmountService = services.NewRetrieveCheckoutAmount(
		&theCheckoutRepositoryMock,
		&theProductRepositoryMock,
		&theProductWithPromotionRepositoryMock,
		&theProductWithDiscountRepositoryMock)

	req, _ := http.NewRequest("GET", "/checkouts/"+checkout.Id+"/amount", nil)
	response := executeRequest(req)

	var responseCheckout responses.Checkout
	json.Unmarshal(response.Body.Bytes(), &responseCheckout)
	assert.EqualValues(t, 200, response.Code)
	assert.EqualValues(t, "7.50€", responseCheckout.Amount)
	theCheckoutRepositoryMock.AssertExpectations(t)
	theProductRepositoryMock.AssertExpectations(t)
	theProductWithDiscountRepositoryMock.AssertExpectations(t)
	theProductWithPromotionRepositoryMock.AssertExpectations(t)
}

func TestReturn404RetrievingCheckoutAmountWhenCheckoutDoesNotExists(t *testing.T) {
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", "a_fake_checkout").Return(models.Checkout{}, false)
	theProductRepositoryMock := mocks.ProductRepositoryMock{}
	theProductWithDiscountRepositoryMock := mocks.ProductWithDiscountRepositoryMock{}
	theProductWithPromotionRepositoryMock := mocks.ProductWithPromotionRepositoryMock{}
	app.RetrieveCheckoutAmountService = services.NewRetrieveCheckoutAmount(
		&theCheckoutRepositoryMock,
		&theProductRepositoryMock,
		&theProductWithPromotionRepositoryMock,
		&theProductWithDiscountRepositoryMock)

	req, _ := http.NewRequest("GET", "/checkouts/a_fake_checkout/amount", nil)
	response := executeRequest(req)

	var checkoutNotFound responses.CheckoutNotFound
	json.Unmarshal(response.Body.Bytes(), &checkoutNotFound)

	assert.EqualValues(t, 404, response.Code)
	assert.EqualValues(t, "Checkout a_fake_checkout not found", checkoutNotFound.Message)
	theCheckoutRepositoryMock.AssertExpectations(t)
}

func TestReturn204WhenDeleteCheckout(t *testing.T) {
	checkout := ACheckout()
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", checkout.Id).Return(checkout, true)
	theCheckoutRepositoryMock.On("Delete", checkout)
	app.DeleteCheckoutService = services.NewDeleteCheckout(&theCheckoutRepositoryMock)

	req, _ := http.NewRequest("DELETE", "/checkouts/"+checkout.Id, nil)
	response := executeRequest(req)

	assert.EqualValues(t, 204, response.Code)
	theCheckoutRepositoryMock.AssertExpectations(t)
}

func TestReturn404DeletingCheckoutWhenCheckoutDoesNotExists(t *testing.T) {
	theCheckoutRepositoryMock := mocks.CheckoutRepositoryMock{}
	theCheckoutRepositoryMock.On("SearchById", "a_fake_checkout").Return(models.Checkout{}, false)
	app.DeleteCheckoutService = services.NewDeleteCheckout(&theCheckoutRepositoryMock)

	req, _ := http.NewRequest("DELETE", "/checkouts/a_fake_checkout", nil)
	response := executeRequest(req)

	var checkoutNotFound responses.CheckoutNotFound
	json.Unmarshal(response.Body.Bytes(), &checkoutNotFound)
	assert.EqualValues(t, 404, response.Code)
	assert.EqualValues(t, "Checkout a_fake_checkout not found", checkoutNotFound.Message)
}
