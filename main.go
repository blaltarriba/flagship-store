package main

import (
	"lana/flagship-store/models"
	"lana/flagship-store/persistence"
	"lana/flagship-store/services"
)

func main() {
	app := App{}
	checkoutRepository := populate_checkouts()
	productRepository := populate_products()
	createCheckoutService := services.NewCreateCheckout(checkoutRepository, productRepository)
	app.Initialize(checkoutRepository, productRepository, populate_products_with_promotion(), populate_products_with_discount(), createCheckoutService)

	app.Run(":10000")
}

func populate_checkouts() persistence.CheckoutRepository {
	checkouts := make(map[string]models.Checkout)
	return persistence.NewCheckoutRepository(checkouts)
}

func populate_products() persistence.ProductRepository {
	products := make(map[string]models.Product)
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
	return persistence.NewProductsRepository(products)
}

func populate_products_with_promotion() persistence.ProductRepository {
	products := make(map[string]models.Product)
	pen := models.Product{
		Code:  "PEN",
		Name:  "Lana Pen",
		Price: 500,
	}
	products[pen.Code] = pen
	return persistence.NewProductWithPromotionRepository(products)
}

func populate_products_with_discount() persistence.ProductRepository {
	products := make(map[string]models.Product)
	tshirt := models.Product{
		Code:  "TSHIRT",
		Name:  "Lana T-Shirt",
		Price: 2000,
	}
	products[tshirt.Code] = tshirt
	return persistence.NewProductWithDiscountRepository(products)
}
