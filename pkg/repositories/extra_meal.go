package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
)

type ExtraMealRepo struct {
	db *gorm.DB
}

func ExtraMealDBInstance(DB *gorm.DB) domain.IExtraMealRepo {
	return &ExtraMealRepo{
		db: DB,
	}
}

func (repo *ExtraMealRepo) GenerateExtraMeal(date string) error {
	var cnt int64 = 0
	err := repo.db.Model(&models.ExtraMeal{}).Where("date = ?", date).Count(&cnt).Error
	if err != nil {
		fmt.Println(err)
		return err
	}
	if cnt == 0 {
		repo.db.Create(&models.ExtraMeal{
			Date:       date,
			LunchCount: 0,
			SnackCount: 0,
		})
		log.Println("Extra Meals generated for date:", date)
	}
	return nil
}

//
//func (repo *ExtraMealRepo) UpdateExtraMeal(date string, count int) error {
//	if err := repo.db.Model(&models.ExtraMeal{}).
//		Where("date = ?", date).
//		Updates(models.ExtraMeal{
//			Date:  date,
//			Count: count,
//		}).Error; err != nil {
//		fmt.Println(err)
//		return err
//	}
//	return nil
//}

func (repo *ExtraMealRepo) UpdateExtraMeal(date string, LunchCount int, SnackCount int) error {
	if err := repo.db.Model(&models.ExtraMeal{}).
		Where("date = ?", date).
		Updates(map[string]interface{}{
			"date":        date,
			"lunch_count": LunchCount,
			"snack_count": SnackCount,
		}).Error; err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (repo *ExtraMealRepo) FetchExtraMeal(date string) (models.ExtraMeal, error) {
	var extraMeal models.ExtraMeal
	err := repo.db.Where("date = ?", date).First(&extraMeal).Error
	if err != nil {
		return models.ExtraMeal{}, err
	}
	return extraMeal, err
}
