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

func main() {
	fmt.Println("My first Golang application")

	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/checkouts", createNewCheckout).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func createNewCheckout(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var productCommand commands.Product
	json.Unmarshal(reqBody, &productCommand)

	products := []string{
		productCommand.Code,
	}
	var checkout models.Checkout
	checkout.Id = 1234
	checkout.Products = products
	Checkouts = append(Checkouts, checkout)

	json.NewEncoder(w).Encode(checkout)
}

var Checkouts []models.Checkout
