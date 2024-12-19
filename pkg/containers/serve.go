package containers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"meal-management/pkg/config"
	"meal-management/pkg/connection"
	"meal-management/pkg/controllers"
	"meal-management/pkg/repositories"
	"meal-management/pkg/routes"
	"meal-management/pkg/services"
)

func Serve(e *echo.Echo) {
	config.SetConfig()

	db := connection.GetDB()

	EmployeeRepo := repositories.EmloyeeDBInstance(db)
	MealPlanRepo := repositories.MealPlanDBInstance(db)

	EmployeeService := services.EmployeeServiceInstance(EmployeeRepo)
	MealPlanService := services.MealPlanServiceInstance(MealPlanRepo)

	controllers.SetEmployeeService(EmployeeService)
	controllers.SetMealPlanService(MealPlanService)

	routes.EmployeeRoutes(e)
	routes.MealPlanRoutes(e)

	log.Fatal(e.Start(fmt.Sprintf(":%s", config.LocalConfig.ServerPort)))
}
