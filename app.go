package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lana/flagship-store/models"
	"lana/flagship-store/persistence"
	"lana/flagship-store/services/commands"
	"lana/flagship-store/services/responses"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type App struct {
	Router                *mux.Router
	CheckoutRepository    persistence.CheckoutRepository
	ProductRepository     persistence.ProductRepository
	ProductsWithPromotion map[string]models.Product
	ProductsWithDiscount  map[string]models.Product
}

func (app *App) Initialize(checkoutRepository persistence.CheckoutRepository, productsRepository persistence.ProductRepository, productsWithPromotion map[string]models.Product, productsWithDiscount map[string]models.Product) {
	app.CheckoutRepository = checkoutRepository
	app.ProductRepository = productsRepository
	app.ProductsWithPromotion = productsWithPromotion
	app.ProductsWithDiscount = productsWithDiscount
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

	if _, existProduct := app.ProductRepository.SearchById(productCommand.Code); !existProduct {
		response.WriteHeader(http.StatusNotFound)
		productNotFound := responses.ProductNotFound{
			Message: "Product " + productCommand.Code + " not found",
		}
		json.NewEncoder(response).Encode(productNotFound)
		return
	}

	checkout := models.Checkout{
		Id:       uuid.NewString(),
		Products: []string{productCommand.Code},
	}

	app.CheckoutRepository.Persist(checkout)

	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(checkout)
}

func (app *App) addProductToCheckout(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	var addProductCommand commands.AddProduct
	json.Unmarshal(body, &addProductCommand)

	vars := mux.Vars(request)
	id := vars["id"]

	checkout, existCheckout := app.CheckoutRepository.SearchById(id)
	if !existCheckout {
		response.WriteHeader(http.StatusNotFound)
		checkoutNotFound := responses.CheckoutNotFound{
			Message: "Checkout " + id + " not found",
		}
		json.NewEncoder(response).Encode(checkoutNotFound)
		return
	}

	if _, existProduct := app.ProductRepository.SearchById(addProductCommand.Code); !existProduct {
		response.WriteHeader(http.StatusUnprocessableEntity)
		productNotFound := responses.ProductNotFound{
			Message: "Product " + addProductCommand.Code + " not found",
		}
		json.NewEncoder(response).Encode(productNotFound)
		return
	}

	checkout.Products = append(checkout.Products, addProductCommand.Code)
	app.CheckoutRepository.Persist(checkout)

	response.WriteHeader(http.StatusNoContent)
}

func (app *App) retrieveCheckoutAmount(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	checkout, existCheckout := app.CheckoutRepository.SearchById(id)
	if !existCheckout {
		response.WriteHeader(http.StatusNotFound)
		checkoutNotFound := responses.CheckoutNotFound{
			Message: "Checkout " + id + " not found",
		}
		json.NewEncoder(response).Encode(checkoutNotFound)
		return
	}

	checkoutAmount := calculateCheckoutAmount(checkout.Products, app.ProductRepository, app.ProductsWithPromotion, app.ProductsWithDiscount)
	responseCheckout := responses.Checkout{
		Amount: formatCheckoutAmount(checkoutAmount),
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(responseCheckout)
}

func calculateCheckoutAmount(checkoutProducts []string, productsRepository persistence.ProductRepository, productsWithPromotion map[string]models.Product, productsWithDiscount map[string]models.Product) int {
	productRealUnits := calculateRealProductUnits(checkoutProducts)
	productUnits := calculatePayableProductUnits(productRealUnits, productsWithPromotion)

	var amount int
	for productCode, quantity := range productUnits {
		product, _ := productsRepository.SearchById(productCode)
		if _, hasDiscount := productsWithDiscount[productCode]; hasDiscount {
			amount += calculateAmountWithDiscount(quantity, product.Price)
			continue
		}
		amount += (product.Price * quantity)
	}
	return amount
}

func calculateRealProductUnits(checkoutProducts []string) map[string]int {
	var productUnits = map[string]int{"PEN": 0, "TSHIRT": 0, "MUG": 0}
	for _, checkoutProductCode := range checkoutProducts {
		productUnits[checkoutProductCode] += 1
	}
	return productUnits
}

func calculatePayableProductUnits(productRealUnits map[string]int, productsWithPromotion map[string]models.Product) map[string]int {
	var productUnits = map[string]int{"PEN": 0, "TSHIRT": 0, "MUG": 0}
	for productCode, quantity := range productRealUnits {
		if _, found := productsWithPromotion[productCode]; found {
			productUnits[productCode] = calculatePayableUnitsApplying2X1Promotion(quantity)
			continue
		}
		productUnits[productCode] = quantity
	}
	return productUnits
}

func calculatePayableUnitsApplying2X1Promotion(quantity int) int {
	if quantity == 0 {
		return 0
	}
	if quantity%2 == 0 {
		return quantity / 2
	}
	return ((quantity - 1) / 2) + 1
}

func calculateAmountWithDiscount(quantity int, price int) int {
	if quantity < 3 {
		return price * quantity
	}
	unitPriceWithDiscount := (price * 75) / 100
	return unitPriceWithDiscount * quantity
}

func formatCheckoutAmount(amount int) string {
	amount_with_decimals := float64(amount) / 100
	amount_with_fixed_decimals := strconv.FormatFloat(amount_with_decimals, 'f', 2, 64)
	return amount_with_fixed_decimals + "€"
}

func (app *App) deleteCheckout(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	checkout, existCheckout := app.CheckoutRepository.SearchById(id)
	if !existCheckout {
		response.WriteHeader(http.StatusNotFound)
		checkoutNotFound := responses.CheckoutNotFound{
			Message: "Checkout " + id + " not found",
		}
		json.NewEncoder(response).Encode(checkoutNotFound)
		return
	}

	app.CheckoutRepository.Delete(checkout)

	response.WriteHeader(http.StatusNoContent)
}
