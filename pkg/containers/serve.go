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

	EmployeeRepo := repositories.EmployeeDBInstance(db)
	MealPlanRepo := repositories.MealPlanDBInstance(db)
	LoginRepo := repositories.LoginDBInstance(db)

	EmployeeService := services.EmployeeServiceInstance(EmployeeRepo)
	MealPlanService := services.MealPlanServiceInstance(MealPlanRepo)
	LoginService := services.LoginServiceInstance(LoginRepo)

	controllers.SetEmployeeService(EmployeeService)
	controllers.SetMealPlanService(MealPlanService)
	controllers.SetLoginService(LoginService)

	routes.EmployeeRoutes(e)
	routes.MealPlanRoutes(e)
	routes.LoginRoutes(e)
	routes.AuthRoutes(e)

	log.Fatal(e.Start(fmt.Sprintf(":%s", config.LocalConfig.ServerPort)))
}
