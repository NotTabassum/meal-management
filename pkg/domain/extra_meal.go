package domain

import "meal-management/pkg/models"

type IExtraMealRepo interface {
	GenerateExtraMeal(date string) error
	UpdateExtraMeal(date string, LunchCount int, SnackCount int) error
	FetchExtraMeal(date string) (models.ExtraMeal, error)
}

type IExtraMealService interface {
	GenerateExtraMeal() error
	UpdateExtraMeal(date string, LunchCount int, SnackCount int) error
	FetchExtraMeal(date string) (models.ExtraMeal, error)
}
