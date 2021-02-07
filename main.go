package main

import (
	"lana/flagship-store/models"
)

func main() {
	app := App{}
	app.Initialize(checkouts, populate_products(products))

	app.Run(":10000")
}

func populate_products(products []models.Product) []models.Product {
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
	products = append(products, pen)
	products = append(products, tshirt)
	products = append(products, mug)
	return products
}

var checkouts []models.Checkout
var products []models.Product
