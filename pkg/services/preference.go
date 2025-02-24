package services

import (
	"errors"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
)

type PreferenceService struct {
	repo domain.IPreferenceRepo
}

func PreferenceServiceInstance(prefRepo domain.IPreferenceRepo) domain.IPreferenceService {
	return &PreferenceService{
		repo: prefRepo,
	}
}
func (service *PreferenceService) CreatePreferenceService(pref *models.Preference) error {
	if err := service.repo.CreatePreferenceRepo(pref); err != nil {
		return errors.New("preference was not created")
	}
	return nil
}

func (service *PreferenceService) GetPreferenceService() ([]*models.Preference, error) {
	pref, err := service.repo.GetPreference()
	if err != nil {
		return []*models.Preference{}, err
	}
	return pref, nil
}
