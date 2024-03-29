package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lana/flagship-store/services"
	"lana/flagship-store/services/commands"
	"lana/flagship-store/services/errors"
	"lana/flagship-store/services/responses"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type App struct {
	Router                        *mux.Router
	CreateCheckoutService         services.CreateCheckout
	AddProductToCheckoutService   services.AddProductToCheckout
	RetrieveCheckoutAmountService services.RetrieveCheckoutAmount
	DeleteCheckoutService         services.DeleteCheckout
}

func (app *App) Initialize(createCheckoutService services.CreateCheckout, addProductToCheckoutService services.AddProductToCheckout, deleteCheckoutService services.DeleteCheckout, retrieveCheckoutAmountService services.RetrieveCheckoutAmount) {
	app.CreateCheckoutService = createCheckoutService
	app.AddProductToCheckoutService = addProductToCheckoutService
	app.RetrieveCheckoutAmountService = retrieveCheckoutAmountService
	app.DeleteCheckoutService = deleteCheckoutService
	app.Router = mux.NewRouter().StrictSlash(true)
	app.initializeRoutes()
}

func (app *App) Run(addr string) {
	fmt.Println("My first Golang application")
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/checkouts", app.createCheckout).Methods("POST")
	app.Router.HandleFunc("/checkouts/{id}", app.addProductToCheckout).Methods("PATCH")
	app.Router.HandleFunc("/checkouts/{id}", app.deleteCheckout).Methods("DELETE")
	app.Router.HandleFunc("/checkouts/{id}/amount", app.retrieveCheckoutAmount).Methods("GET")
}

func (app *App) createCheckout(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	var productCommand commands.Product
	json.Unmarshal(body, &productCommand)

	checkout, err := app.CreateCheckoutService.Do(productCommand)

	if _, ok := err.(*errors.ProductNotFoundError); ok {
		response.WriteHeader(http.StatusNotFound)
		productNotFound := responses.ProductNotFound{
			Message: "Product " + productCommand.Code + " not found",
		}
		json.NewEncoder(response).Encode(productNotFound)
		return
	}

	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(checkout)
}

func (app *App) addProductToCheckout(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	var addProductCommand commands.AddProduct
	json.Unmarshal(body, &addProductCommand)

	vars := mux.Vars(request)
	id := vars["id"]

	_, err := app.AddProductToCheckoutService.Do(addProductCommand, id)

	if _, isThisError := err.(*errors.CheckoutNotFoundError); isThisError {
		response.WriteHeader(http.StatusNotFound)
		checkoutNotFound := responses.CheckoutNotFound{
			Message: "Checkout " + id + " not found",
		}
		json.NewEncoder(response).Encode(checkoutNotFound)
		return
	}

	if _, isThisError := err.(*errors.ProductNotFoundError); isThisError {
		response.WriteHeader(http.StatusUnprocessableEntity)
		productNotFound := responses.ProductNotFound{
			Message: "Product " + addProductCommand.Code + " not found",
		}
		json.NewEncoder(response).Encode(productNotFound)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func (app *App) retrieveCheckoutAmount(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	amount, err := app.RetrieveCheckoutAmountService.Do(id)

	if _, ok := err.(*errors.CheckoutNotFoundError); ok {
		response.WriteHeader(http.StatusNotFound)
		checkoutNotFound := responses.CheckoutNotFound{
			Message: "Checkout " + id + " not found",
		}
		json.NewEncoder(response).Encode(checkoutNotFound)
		return
	}

	responseCheckout := responses.Checkout{
		Amount: formatCheckoutAmount(amount),
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(responseCheckout)
}

func formatCheckoutAmount(amount int) string {
	amount_with_decimals := float64(amount) / 100
	amount_with_fixed_decimals := strconv.FormatFloat(amount_with_decimals, 'f', 2, 64)
	return amount_with_fixed_decimals + "€"
}

func (app *App) deleteCheckout(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	_, err := app.DeleteCheckoutService.Do(id)

	if _, ok := err.(*errors.CheckoutNotFoundError); ok {
		response.WriteHeader(http.StatusNotFound)
		checkoutNotFound := responses.CheckoutNotFound{
			Message: "Checkout " + id + " not found",
		}
		json.NewEncoder(response).Encode(checkoutNotFound)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}
