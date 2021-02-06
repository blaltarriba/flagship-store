package models

type Checkout struct {
	Id       string   `json:"id"`
	Products []string `json:"products"`
}
