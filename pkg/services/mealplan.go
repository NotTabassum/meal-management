package services

import (
	"errors"
	"meal-management/pkg/consts"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
	"meal-management/pkg/types"
	"time"
)

type MealPlanService struct {
	repo domain.IMealPlanRepo
}

func MealPlanServiceInstance(MealPlanRepo domain.IMealPlanRepo) domain.IMealPlanService {
	return &MealPlanService{
		repo: MealPlanRepo,
	}
}

func (service *MealPlanService) CreateMealPlan(MealPlan *models.MealPlan) error {
	if err := service.repo.CreateOrUpdateMealPlan(MealPlan); err != nil {
		return err
	}
	return nil
}

func (service *MealPlanService) GetMealPlanByPrimaryKey(Date string, MealType string) (models.MealPlan, error) {
	mealPlan, err := service.repo.GetMealPlanByPrimaryKey(Date, MealType)
	if err != nil {
		return models.MealPlan{}, err
	}

	response := models.MealPlan{
		Date:           mealPlan.Date,
		MealType:       mealPlan.MealType,
		Food:           mealPlan.Food,
		PreferenceFood: mealPlan.PreferenceFood,
	}

	return response, nil
}
func (service *MealPlanService) GetMealPlan(startDate string, days int) ([]types.GetMealPlanResponse, error) {
	var mealPlans []types.GetMealPlanResponse
	tempStDate, err := time.Parse(consts.DateFormat, startDate)
	if err != nil {
		return nil, err
	}
	tmpEndDate := tempStDate.AddDate(0, 0, days-1)
	endDate := tmpEndDate.Format(consts.DateFormat)
	meal := service.repo.GetMealPlan(startDate, endDate)

	//if len(meal) == 0 {
	//	return nil, errors.New("No Meal Plan Found.")
	//}

	groupedMeals := make(map[string][]types.Menu)

	for _, meal := range meal {
		menu := types.Menu{
			MealType:       meal.MealType,
			Food:           meal.Food,
			PreferenceFood: meal.PreferenceFood,
		}
		groupedMeals[meal.Date] = append(groupedMeals[meal.Date], menu)
	}
	for date, menu := range groupedMeals {
		mealPlans = append(mealPlans, types.GetMealPlanResponse{
			Date: date,
			Menu: menu,
		})
	}
	return mealPlans, nil
}
func (service *MealPlanService) UpdateMealPlan(meal *models.MealPlan) error {
	if err := service.repo.UpdateMealPlan(meal); err != nil {
		return errors.New("mealplan update was unsuccesful")
	}
	return nil
}
func (service *MealPlanService) DeleteMealPlan(date string, MealType string) error {
	if err := service.repo.DeleteMealPlan(date, MealType); err != nil {
		return errors.New("MealPlan was not deleted")
	}
	return nil
}
