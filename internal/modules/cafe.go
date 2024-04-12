package modules

import (
	"barista/pkg/models"
	"barista/pkg/repo"
	"context"
)

type CafeHandler struct {
	CafeRepo repo.CafesRepo
}

func (c CafeHandler) Create(ctx context.Context, cafe *models.Cafe) error {
	panic("implement me")
}

func (c CafeHandler) GetCafes() {
	panic("implement me")
}

func (c CafeHandler) GetCafeByID() {
	panic("implement me")
}
