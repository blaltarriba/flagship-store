package main

import (
	"lana/flagship-store/models"
)

func main() {
	app := App{}
	checkouts = make(map[string]models.Checkout)
	app.Initialize(checkouts, populate_products(products), populate_products_with_promotion(productsWithPromotion), populate_products_with_discount(productsWithDiscount))

	app.Run(":10000")
}

func populate_products(products map[string]models.Product) map[string]models.Product {
	products = make(map[string]models.Product)
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

	products[pen.Code] = pen
	products[tshirt.Code] = tshirt
	products[mug.Code] = mug
	return products
}

func populate_products_with_promotion(products map[string]models.Product) map[string]models.Product {
	products = make(map[string]models.Product)
	pen := models.Product{
		Code:  "PEN",
		Name:  "Lana Pen",
		Price: 500,
	}
	products[pen.Code] = pen
	return products
}

func populate_products_with_discount(products map[string]models.Product) map[string]models.Product {
	products = make(map[string]models.Product)
	tshirt := models.Product{
		Code:  "TSHIRT",
		Name:  "Lana T-Shirt",
		Price: 2000,
	}
	products[tshirt.Code] = tshirt
	return products
}

var checkouts map[string]models.Checkout
var products map[string]models.Product
var productsWithPromotion map[string]models.Product
var productsWithDiscount map[string]models.Product
