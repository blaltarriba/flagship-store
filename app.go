package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lana/flagship-store/models"
	"lana/flagship-store/services/commands"
	"lana/flagship-store/services/responses"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type App struct {
	Router    *mux.Router
	Checkouts []models.Checkout
	Products  map[string]models.Product
}

func (app *App) Initialize(checkouts []models.Checkout, products map[string]models.Product) {
	app.Checkouts = checkouts
	app.Products = products
	app.Router = mux.NewRouter().StrictSlash(true)
	app.initializeRoutes()
}

func (app *App) Run(addr string) {
	fmt.Println("My first Golang application")
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/checkouts", app.createNewCheckout).Methods("POST")
	app.Router.HandleFunc("/checkouts/{id}", app.addProductToCheckout).Methods("PATCH")
	app.Router.HandleFunc("/checkouts/{id}/amount", app.retrieveCheckoutAmount).Methods("GET")
}

func (app *App) createNewCheckout(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	var productCommand commands.Product
	json.Unmarshal(body, &productCommand)

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{productCommand.Code},
	}

	app.Checkouts = append(app.Checkouts, checkout)

	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(checkout)
}

func (app *App) addProductToCheckout(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	var addProductCommand commands.AddProduct
	json.Unmarshal(body, &addProductCommand)

	vars := mux.Vars(request)
	id := vars["id"]

	for index, checkout := range app.Checkouts {
		if checkout.Id == id {
			checkout.Products = append(checkout.Products, addProductCommand.Code)
			app.Checkouts[index] = checkout
			break
		}
	}

	response.WriteHeader(http.StatusNoContent)
}

func (app *App) retrieveCheckoutAmount(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	checkout := searchCheckoutById(id, app.Checkouts)
	checkoutAmount := calculateCheckoutAmount(checkout.Products, app.Products)
	responseCheckout := responses.Checkout{
		Amount: formatCheckoutAmount(checkoutAmount),
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(responseCheckout)
}

func searchCheckoutById(id string, checkouts []models.Checkout) models.Checkout {
	var checkout models.Checkout

	for _, currentCheckout := range checkouts {
		if currentCheckout.Id == id {
			checkout = currentCheckout
			break
		}
	}

	return checkout
}

func calculateCheckoutAmount(checkoutProducts []string, products map[string]models.Product) int {
	var amount int
	for _, checkoutProductCode := range checkoutProducts {
		product := products[checkoutProductCode]
		amount += product.Price
	}
	return amount
}

func formatCheckoutAmount(amount int) float64 {
	return float64(amount) / 100
}
