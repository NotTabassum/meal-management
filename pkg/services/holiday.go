package services

import (
	"fmt"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
	"time"
)

type HolidayService struct {
	repo domain.IHolidayRepo
}

func HolidayServiceInstance(holidayRepo domain.IHolidayRepo) domain.IHolidayService {
	return &HolidayService{
		repo: holidayRepo,
	}
}

func (service *HolidayService) CreateHoliday(holiday []models.Holiday) ([]string, error) {
	var upcomingHolidays []string
	for _, reqHoliday := range holiday {
		if isHolidayWithinNext30Days(reqHoliday.Date) {
			upcomingHolidays = append(upcomingHolidays, reqHoliday.Date)
		}
		holiday := &models.Holiday{
			Date:    reqHoliday.Date,
			Remarks: reqHoliday.Remarks,
		}

		if err := service.repo.CreateHoliday(holiday); err != nil {
			return []string{}, err
		}
	}
	return upcomingHolidays, nil
}

func isHolidayWithinNext30Days(holidayDate string) bool {
	holiday, err := time.Parse("2006-01-02", holidayDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return false
	}

	today := time.Now()
	thirtyDaysFromToday := today.Add(30 * 24 * time.Hour)

	return holiday.After(today) && holiday.Before(thirtyDaysFromToday)
}
