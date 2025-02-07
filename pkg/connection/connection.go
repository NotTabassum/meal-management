package connection

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"meal-management/pkg/config"
	"meal-management/pkg/models"
)

var db *gorm.DB

func Connect() {
	dbConfig := config.LocalConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPass, dbConfig.DBIP, dbConfig.DBPort, dbConfig.DBName)
	fmt.Println(dsn)
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database")
		panic(err)
	}
	fmt.Println("Database Connected")
	db = d
}

func migrate() {
	err := db.Migrator().AutoMigrate(&models.Employee{})
	if err != nil {
		return
	}
	err = db.Migrator().AutoMigrate(&models.MealPlan{})
	if err != nil {
		return
	}
	err = db.Migrator().AutoMigrate(&models.MealActivity{})
	if err != nil {
		return
	}
	err = db.Migrator().AutoMigrate(&models.Department{})
	if err != nil {
		return
	}
	err = db.Migrator().AutoMigrate(&models.ExtraMeal{})
	if err != nil {
		return
	}
}

func GetDB() *gorm.DB {
	if db == nil {
		Connect()
	}
	//jhaamela
	migrate()
	return db
}
