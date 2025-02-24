package repositories

import (
	"gorm.io/gorm"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
)

type HolidayRepo struct {
	db *gorm.DB
}

func HolidayDBInstance(d *gorm.DB) domain.IHolidayRepo {
	return &HolidayRepo{
		db: d,
	}
}

func (repo *HolidayRepo) CreateHoliday(holiday *models.Holiday) error {
	if err := repo.db.Create(holiday).Error; err != nil {
		return err
	}
	return nil
}
