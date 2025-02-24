package services

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"meal-management/pkg/consts"
	"meal-management/pkg/domain"
	"meal-management/pkg/models"
	"sort"
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

func (service *HolidayService) CreateHoliday(holiday []models.Holiday) ([]string, []string, error) {
	var upcomingHolidays []string
	failedHolidays := make([]string, 0)
	for _, reqHoliday := range holiday {
		if isHolidayWithinNext30Days(reqHoliday.Date) {
			upcomingHolidays = append(upcomingHolidays, reqHoliday.Date)
		}
		holiday := &models.Holiday{
			Date:    reqHoliday.Date,
			Remarks: reqHoliday.Remarks,
		}

		err := service.repo.CreateHoliday(holiday)
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			fmt.Printf("Holiday with date %s already exists (Duplicate Entry)\n", reqHoliday.Date)
			failedHolidays = append(failedHolidays, reqHoliday.Date)
			continue
		}
		if err != nil {
			return nil, nil, err
		}
	}
	return failedHolidays, upcomingHolidays, nil
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

func (service *HolidayService) GetHoliday() ([]models.Holiday, error) {
	holidays, err := service.repo.GetHoliday()
	if err != nil {
		return []models.Holiday{}, err
	}
	sort.SliceStable(holidays, func(i, j int) bool {
		dateI, errI := time.Parse(consts.DateFormat, holidays[i].Date)
		dateJ, errJ := time.Parse(consts.DateFormat, holidays[j].Date)

		if errI != nil || errJ != nil {
			return false
		}

		return dateI.Before(dateJ)
	})
	return holidays, nil
}
