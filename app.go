package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lana/flagship-store/models"
	"lana/flagship-store/services/commands"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type App struct {
	Router    *mux.Router
	Checkouts []models.Checkout
	Products  []models.Product
}

func (app *App) Initialize(checkouts []models.Checkout, products []models.Product) {
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

	var checkout models.Checkout

	for i := range app.Checkouts {
		if app.Checkouts[i].Id == id {
			checkout = app.Checkouts[i]
			checkout.Products = append(checkout.Products, addProductCommand.Code)
			app.Checkouts[i] = checkout
			break
		}
	}

	response.WriteHeader(http.StatusNoContent)
}
