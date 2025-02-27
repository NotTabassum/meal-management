package repositories

import (
	"gorm.io/gorm"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
)

type MealPlanRepo struct {
	db *gorm.DB
}

func MealPlanDBInstance(d *gorm.DB) domain.IMealPlanRepo {
	return &MealPlanRepo{
		db: d,
	}
}

func (repo *MealPlanRepo) CreateMealPlan(mealPlan *models.MealPlan) error {
	if err := repo.db.Create(mealPlan).Error; err != nil {
		return err
	}
	return nil
}

func (repo *MealPlanRepo) CreateOrUpdateMealPlan(mealPlan *models.MealPlan) error {
	var existingMealPlan models.MealPlan

	result := repo.db.Where("date = ? AND meal_type = ?", mealPlan.Date, mealPlan.MealType).Find(&existingMealPlan)

	if result.RowsAffected > 0 {
		existingMealPlan.Food = mealPlan.Food
		existingMealPlan.PreferenceFood = mealPlan.PreferenceFood
		if err := repo.db.Save(&existingMealPlan).Error; err != nil {
			return err
		}
	} else {
		if err := repo.db.Create(mealPlan).Error; err != nil {
			return err
		}
	}

	return nil
}

func (repo *MealPlanRepo) GetMealPlanByPrimaryKey(Date string, MealType string) (*models.MealPlan, error) {
	var mealPlan models.MealPlan
	err := repo.db.Where("date = ? AND meal_type = ?", Date, MealType).First(&mealPlan).Error
	if err != nil {
		return nil, err
	}
	return &mealPlan, nil
}

func (repo *MealPlanRepo) GetMealPlan(startDate, endDate string) []models.MealPlan {
	var mealPlans []models.MealPlan
	var err error

	err = repo.db.Where("date >= ? AND date <= ?", startDate, endDate).Find(&mealPlans).Error

	if err != nil {
		return []models.MealPlan{}
	}
	return mealPlans
}
func (repo *MealPlanRepo) UpdateMealPlan(mealPlan *models.MealPlan) error {
	if err := repo.db.Save(mealPlan).Error; err != nil {
		return err
	}
	return nil
}
func (repo *MealPlanRepo) DeleteMealPlan(date string, mealType string) error {
	var meal models.MealPlan
	if err := repo.db.Where("date = ? AND meal_type = ?", date, mealType).Delete(&meal).Error; err != nil {
		return err
	}
	return nil
}
