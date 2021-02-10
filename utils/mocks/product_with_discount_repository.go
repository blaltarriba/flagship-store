package mocks

import (
	"lana/flagship-store/models"

	"github.com/stretchr/testify/mock"
)

type ProductWithDiscountRepositoryMock struct {
	mock.Mock
}

func (repository *ProductWithDiscountRepositoryMock) SearchById(id string) (models.Product, bool) {
	args := repository.Called(id)
	return args.Get(0).(models.Product), args.Bool(1)
}
