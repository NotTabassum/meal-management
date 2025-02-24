package repositories

import (
	"gorm.io/gorm"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
)

type PreferenceRepo struct {
	db *gorm.DB
}

func PreferenceDBInstance(d *gorm.DB) domain.IPreferenceRepo {
	return &PreferenceRepo{
		db: d,
	}
}

func (repo *PreferenceRepo) CreatePreferenceRepo(pref *models.Preference) error {
	if err := repo.db.Create(pref).Error; err != nil {
		return err
	}
	return nil
}

func (repo *PreferenceRepo) GetPreference() ([]*models.Preference, error) {
	var pref []*models.Preference
	if err := repo.db.Find(&pref).Error; err != nil {
		return nil, err
	}
	return pref, nil
}
