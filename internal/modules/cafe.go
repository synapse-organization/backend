package modules

import (
	"barista/pkg/log"
	"barista/pkg/models"
	"barista/pkg/repo"
	"context"
)

type CafeHandler struct {
	CafeRepo  repo.CafesRepo
	Rating    repo.RatingsRepo
	ImageRepo repo.ImageRepo
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

func (c CafeHandler) SearchCafe(ctx context.Context, name string, address string, location string, category string) ([]models.Cafe, error) {
	cafes, err := c.CafeRepo.SearchCafe(ctx, name, address, location, category)
	if err != nil {
		log.GetLog().Errorf("Unable to search cafe. error: %v", err)
		return nil, err
	}
	for i, cafe := range cafes {
		cafes[i].Rating, err = c.Rating.GetCafesRating(ctx, cafe.ID)
		if err != nil {
			log.GetLog().Errorf("Unable to get cafe rating. error: %v", err)
		}
	}

	for i := range cafes {
		images, err := c.ImageRepo.GetByCafeID(ctx, cafes[i].ID)
		if err != nil {
			log.GetLog().Errorf("Unable to get cafe images. error: %v", err)
			continue
		}
		for _, image := range images {
			cafes[i].Images = append(cafes[i].Images, image.ID)
		}
	}

	return cafes, err
}
