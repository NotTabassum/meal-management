package controllers

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func StartCronJobMealActivity() {
	c := cron.New()
	_, err := c.AddFunc("0 13 * * *", func() {
		log.Println("Running GenerateMealActivities at:", time.Now())
		if err := MealActivityService.GenerateMealActivities(); err != nil {
			log.Printf("Error generating meal activities: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Failed to schedule cron job: %v", err)
	}
	c.Start()
	log.Println("Cron job started for meal activity")
}

func StartCronJobExtraMeal() {
	c := cron.New()
	_, err := c.AddFunc("55 12 * * *", func() {
		log.Println("Running Extra Meal at:", time.Now())
		if err := ExtraMealService.GenerateExtraMeal(); err != nil {
			log.Printf("Error generating extra meal activities: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Failed to schedule cron job for extra meal: %v", err)
	}
	c.Start()
	log.Println("Cron job started for extra meal")
}

func StartCronJobLunchSummary() {
	c := cron.New()
	_, err := c.AddFunc("55 12 * * *", func() {
		log.Println("Sending Lunch Summary at :", time.Now())
		if err := MealActivityService.LunchSummaryForEmail(); err != nil {
			log.Printf("Error generating lunch summary: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Failed to schedule cron job for lunch summary: %v", err)
	}
	c.Start()
	log.Println("Cron job started for sending lunch summary")

}

func CronJob() {
	c := cron.New()
	_, err := c.AddFunc("45 14 * * *", func() {
		log.Println("Running GenerateMealActivities at:", time.Now())
		if err := MealActivityService.GenerateMealActivities(); err != nil {
			log.Printf("Error generating meal activities: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Failed to schedule cron job: %v", err)
	}

	_, err = c.AddFunc("46 14 * * *", func() {
		log.Println("Running Extra Meal at:", time.Now())
		if err := ExtraMealService.GenerateExtraMeal(); err != nil {
			log.Printf("Error generating extra meal activities: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Failed to schedule cron job for extra meal: %v", err)
	}
	_, err = c.AddFunc("47 14 * * *", func() {
		log.Println("Sending Lunch Summary at :", time.Now())
		if err := MealActivityService.LunchSummaryForEmail(); err != nil {
			log.Printf("Error generating lunch summary: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Failed to schedule cron job for lunch summary: %v", err)
	}
	c.Start()
	log.Println("Cron job started")
}

//func CronJob() {
//	c := cron.New()
//
//	var err error
//
//	_, err = c.AddFunc("30 14 * * *", func() {
//		log.Println("🕛 Running GenerateMealActivities at:", time.Now())
//		if err := MealActivityService.GenerateMealActivities(); err != nil {
//			log.Printf("Error generating meal activities: %v", err)
//		}
//	})
//	if err != nil {
//		log.Fatalf("🚨 Failed to schedule GenerateMealActivities: %v", err)
//	}
//
//	_, err = c.AddFunc("31 14 * * *", func() { // Slightly different time to avoid conflicts
//		log.Println("🍽️ Running Extra Meal at:", time.Now())
//		if err := ExtraMealService.GenerateExtraMeal(); err != nil {
//			log.Printf(" Error generating extra meal activities: %v", err)
//		}
//	})
//	if err != nil {
//		log.Fatalf("🚨 Failed to schedule ExtraMeal: %v", err)
//	}
//
//	_, err = c.AddFunc("32 14 * * *", func() { // Slightly different time
//		log.Println("📩 Sending Lunch Summary at:", time.Now())
//		if err := MealActivityService.LunchSummaryForEmail(); err != nil {
//			log.Printf(" Error generating lunch summary: %v", err)
//		}
//	})
//	if err != nil {
//		log.Fatalf(" Failed to schedule LunchSummaryForEmail: %v", err)
//	}
//
//	c.Start()
//	log.Println(" Cron jobs started successfully.")
//
//}
