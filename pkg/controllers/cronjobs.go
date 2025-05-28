package controllers

import (
	"github.com/robfig/cron/v3"
	"log"
	"meal-management/pkg/config"
	"time"
)

func CronJob() {
	SERVER := config.LocalConfig.SERVER
	if SERVER == "STAGING" || SERVER == "LOCAL" {
		return
	}
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

	_, err = c.AddFunc("45 09 * * *", func() {
		log.Println("Sending Lunch Reminder at :", time.Now())
		if err := MealActivityService.MealUpdateNotification(1); err != nil {
			log.Printf(" Error sending lunch reminder: %v", err)
		}
	})
	if err != nil {
		log.Fatalf(" Failed to schedule ExtraMeal: %v", err)
	}

	_, err = c.AddFunc("45 13 * * *", func() {
		log.Println("Sending Snacks Reminder at :", time.Now())
		if err := MealActivityService.MealUpdateNotification(2); err != nil {
			log.Printf(" Error sending snacks reminder: %v", err)
		}
	})
	if err != nil {
		log.Fatalf(" Failed to schedule ExtraMeal: %v", err)
	}

	_, err = c.AddFunc("05 10 * * *", func() {
		log.Println("Sending Lunch Summary at:", time.Now())
		if err := MealActivityService.LunchSummaryForEmail(); err != nil {
			log.Printf(" Error generating lunch summary: %v", err)
		}
	})
	if err != nil {
		log.Fatalf(" Failed to schedule LunchSummaryForEmail: %v", err)
	}

	_, err = c.AddFunc("05 14 * * *", func() {
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
