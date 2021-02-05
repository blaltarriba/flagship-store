package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lana/flagship-store/models"
	"lana/flagship-store/services/commands"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router    *mux.Router
	Checkouts []models.Checkout
}

func (a *App) Initialize(checkouts []models.Checkout) {

	a.Checkouts = checkouts

	a.Router = mux.NewRouter().StrictSlash(true)
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	fmt.Println("My first Golang application")
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/checkouts", a.createNewCheckout).Methods("POST")
}

func (a *App) createNewCheckout(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var productCommand commands.Product
	json.Unmarshal(reqBody, &productCommand)

	products := []string{
		productCommand.Code,
	}
	var checkout models.Checkout
	checkout.Id = 1234
	checkout.Products = products
	a.Checkouts = append(a.Checkouts, checkout)

	json.NewEncoder(w).Encode(checkout)
}
