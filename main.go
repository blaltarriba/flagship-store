package main

import "lana/flagship-store/models"

func main() {
	a := App{}
	a.Initialize(Checkouts)

	a.Run(":10000")
}

var Checkouts []models.Checkout
