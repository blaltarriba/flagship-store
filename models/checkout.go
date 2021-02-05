package models

type Checkout struct {
	Id       int      `json:"Id"`
	Products []string `json:"Products"`
}
