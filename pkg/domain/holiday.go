package domain

import "meal-management/pkg/models"

type IHolidayRepo interface {
	CreateHoliday(holiday *models.Holiday) error
	GetHoliday() ([]models.Holiday, error)
	DeleteHoliday(date string) error
}

type IHolidayService interface {
	CreateHoliday(holiday []models.Holiday) ([]string, []string, error)
	GetHoliday() ([]models.Holiday, error)
	DeleteHoliday(date string) error
}
