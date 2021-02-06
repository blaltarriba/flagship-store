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
}

func (app *App) Initialize(checkouts []models.Checkout) {
	app.Checkouts = checkouts
	app.Router = mux.NewRouter().StrictSlash(true)
	app.initializeRoutes()
}

func (app *App) Run(addr string) {
	fmt.Println("My first Golang application")
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/checkouts", app.createNewCheckout).Methods("POST")
}

func (app *App) createNewCheckout(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	var productCommand commands.Product
	json.Unmarshal(body, &productCommand)

	checkout := models.Checkout{
		Id:       uuid.New(),
		Products: []string{productCommand.Code},
	}

	app.Checkouts = append(app.Checkouts, checkout)

	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(checkout)
}
