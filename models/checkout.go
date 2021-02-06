package models

import "github.com/google/uuid"

type Checkout struct {
	Id       uuid.UUID `json:"id"`
	Products []string  `json:"products"`
}
