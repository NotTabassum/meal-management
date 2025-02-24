package domain

import "meal-management/pkg/models"

type IHolidayRepo interface {
	CreateHoliday(holiday *models.Holiday) error
}

type IHolidayService interface {
	CreateHoliday(holiday []models.Holiday) ([]string, error)
}
