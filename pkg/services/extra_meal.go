package services

import (
	"errors"
	"fmt"
	"meal-management/pkg/consts"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
	"time"
)

type ExtraMealService struct {
	repo domain.IExtraMealRepo
}

func ExtraMealServiceInstance(extraMealRepo domain.IExtraMealRepo) domain.IExtraMealService {
	return &ExtraMealService{
		repo: extraMealRepo,
	}
}

func (service *ExtraMealService) GenerateExtraMeal() error {
	now := time.Now()
	date := now.Format(consts.DateFormat)
	dates, err := getNext30Dates(date)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, date := range dates {
		err = service.repo.GenerateExtraMeal(date)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func (service *ExtraMealService) UpdateExtraMeal(date string, LunchCount int, SnackCount int) error {
	if err := service.repo.UpdateExtraMeal(date, LunchCount, SnackCount); err != nil {
		return errors.New("failed to update extra meal activity")
	}
	return nil
}

func (service *ExtraMealService) FetchExtraMeal(date string) (models.ExtraMeal, error) {
	extraMeal, err := service.repo.FetchExtraMeal(date)
	if err != nil {
		return models.ExtraMeal{}, err
	}
	return extraMeal, nil
}
