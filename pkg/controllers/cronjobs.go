package controllers

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func CronJob() {
	loc, _ := time.LoadLocation("Asia/Dhaka")
	time.Local = loc

	c := cron.New()

	var err error

	_, err = c.AddFunc("0 0 * * *", func() {
		log.Println("Running GenerateMealActivities at:", time.Now())
		if err := MealActivityService.GenerateMealActivities(); err != nil {
			log.Printf("Error generating meal activities: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule GenerateMealActivities: %v", err)
	}

	_, err = c.AddFunc("0 0 * * *", func() {
		log.Println("Running Extra Meal at:", time.Now())
		if err := ExtraMealService.GenerateExtraMeal(); err != nil {
			log.Printf(" Error generating extra meal activities: %v", err)
		}
	})
	if err != nil {
		log.Fatalf(" Failed to schedule ExtraMeal: %v", err)
	}

	_, err = c.AddFunc("07 17 * * *", func() {
		log.Println("Sending Lunch Summary at:", time.Now())
		if err := MealActivityService.LunchSummaryForEmail(); err != nil {
			log.Printf(" Error generating lunch summary: %v", err)
		}
	})
	if err != nil {
		log.Fatalf(" Failed to schedule LunchSummaryForEmail: %v", err)
	}

	_, err = c.AddFunc("12 17 * * *", func() {
		log.Println("Sending Snacks Summary at:", time.Now())
		if err := MealActivityService.SnackSummaryForEmail(); err != nil {
			log.Printf(" Error generating snacks summary: %v", err)
		}
	})
	if err != nil {
		log.Fatalf(" Failed to schedule SnackSummaryForEmail: %v", err)
	}

	c.Start()
	log.Println(" Cron jobs started successfully.")

}
