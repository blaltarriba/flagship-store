package main

import "lana/flagship-store/models"

func main() {
	app := App{}
	app.Initialize(checkouts)

	app.Run(":10000")
}

var checkouts []models.Checkout
