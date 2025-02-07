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
	MealActivityRepo := repositories.MealActivityDBInstance(db)
	DeptRepo := repositories.DeptDBInstance(db)
	ExtraMealRepo := repositories.ExtraMealDBInstance(db)

	EmployeeService := services.EmployeeServiceInstance(EmployeeRepo)
	MealPlanService := services.MealPlanServiceInstance(MealPlanRepo)
	LoginService := services.LoginServiceInstance(LoginRepo)
	MealActivityService := services.MealActivityServiceInstance(MealActivityRepo)
	DeptService := services.DeptServiceInstance(DeptRepo)
	ExtraMealService := services.ExtraMealServiceInstance(ExtraMealRepo)

	controllers.SetEmployeeService(EmployeeService)
	controllers.SetMealPlanService(MealPlanService)
	controllers.SetLoginService(LoginService)
	controllers.SetMealActivityService(MealActivityService)
	controllers.SetDeptService(DeptService)
	controllers.SetExtraMealService(ExtraMealService)

	routes.EmployeeRoutes(e)
	routes.MealPlanRoutes(e)
	routes.LoginRoutes(e)
	routes.MealActivityRoutes(e)
	routes.DeptRoutes(e)
	routes.ExtraMealRoutes(e)

	log.Fatal(e.Start(fmt.Sprintf(":%s", config.LocalConfig.ServerPort)))
}
