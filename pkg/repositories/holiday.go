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
	if err := repo.db.Save(holiday).Error; err != nil {
		return err
	}
	return nil
}

func (repo *HolidayRepo) GetHoliday() ([]models.Holiday, error) {
	var holiday []models.Holiday
	if err := repo.db.Find(&holiday).Error; err != nil {
		return nil, err
	}
	return holiday, nil
}

func (repo *HolidayRepo) DeleteHoliday(date string) error {
	if err := repo.db.Delete(&models.Holiday{}, "date = ?", date).Error; err != nil {
		return err
	}
	return nil
}
