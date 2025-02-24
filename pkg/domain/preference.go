package domain

import "meal-management/pkg/models"

type IPreferenceRepo interface {
	CreatePreferenceRepo(pref *models.Preference) error
	GetPreference() ([]*models.Preference, error)
}

type IPreferenceService interface {
	CreatePreferenceService(pref *models.Preference) error
	GetPreferenceService() ([]*models.Preference, error)
}
