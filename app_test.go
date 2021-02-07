package main

import (
	"bytes"
	"encoding/json"
	"lana/flagship-store/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var a App

func TestMain(m *testing.M) {
	var checkouts []models.Checkout
	var products []models.Product

	a = App{}
	a.Initialize(checkouts, products)

	code := m.Run()

	clearCheckouts()

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func clearCheckouts() {
	a.Checkouts = nil
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

	assert.EqualValues(t, 1, len(a.Checkouts))
	assert.NotNil(t, createdCheckout.Id)
	assert.EqualValues(t, "PEN", createdCheckout.Products[0])
	assert.EqualValues(t, 1, len(createdCheckout.Products))
}

func TestReturn204WhenAddProductToCheckout(t *testing.T) {
	clearCheckouts()

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{"MUG"},
	}
	a.Checkouts = append(a.Checkouts, checkout)

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
	a.Checkouts = append(a.Checkouts, checkout)

	payload := []byte(`{"product":"PEN"}`)

	req, _ := http.NewRequest("PATCH", "/checkouts/"+checkout.Id, bytes.NewBuffer(payload))
	executeRequest(req)

	modifiedCheckout := a.Checkouts[0]
	assert.EqualValues(t, 2, len(modifiedCheckout.Products))
	assert.EqualValues(t, "PEN", modifiedCheckout.Products[1])
}
