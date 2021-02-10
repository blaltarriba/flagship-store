package mocks

import (
	"lana/flagship-store/models"

	"github.com/stretchr/testify/mock"
)

type ProductWithPromotionRepositoryMock struct {
	mock.Mock
}

func (repository *ProductWithPromotionRepositoryMock) SearchById(id string) (models.Product, bool) {
	args := repository.Called(id)
	return args.Get(0).(models.Product), args.Bool(1)
}
