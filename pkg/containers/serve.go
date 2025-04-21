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

func InitServe() {
	config.SetConfig()

	db := connection.GetDB()

	EmployeeRepo := repositories.EmployeeDBInstance(db)
	MealPlanRepo := repositories.MealPlanDBInstance(db)
	LoginRepo := repositories.LoginDBInstance(db)
	MealActivityRepo := repositories.MealActivityDBInstance(db)
	DeptRepo := repositories.DeptDBInstance(db)
	ExtraMealRepo := repositories.ExtraMealDBInstance(db)
	PreferenceRepo := repositories.PreferenceDBInstance(db)
	HolidayRepo := repositories.HolidayDBInstance(db)

	EmployeeService := services.EmployeeServiceInstance(EmployeeRepo, MealActivityRepo, HolidayRepo)
	MealPlanService := services.MealPlanServiceInstance(MealPlanRepo)
	LoginService := services.LoginServiceInstance(LoginRepo)
	MealActivityService := services.MealActivityServiceInstance(MealActivityRepo, MealPlanService, EmployeeService)
	DeptService := services.DeptServiceInstance(DeptRepo)
	ExtraMealService := services.ExtraMealServiceInstance(ExtraMealRepo)
	PreferenceService := services.PreferenceServiceInstance(PreferenceRepo)
	HolidayService := services.HolidayServiceInstance(HolidayRepo, EmployeeRepo, MealActivityRepo)

	controllers.SetEmployeeService(EmployeeService)
	controllers.SetMealPlanService(MealPlanService)
	controllers.SetLoginService(LoginService)
	controllers.SetMealActivityService(MealActivityService)
	controllers.SetDeptService(DeptService)
	controllers.SetExtraMealService(ExtraMealService)
	controllers.SetPreferenceService(PreferenceService)
	controllers.SetHolidayService(HolidayService)
}

func Serve(e *echo.Echo) {
	routes.EmployeeRoutes(e)
	routes.MealPlanRoutes(e)
	routes.LoginRoutes(e)
	routes.MealActivityRoutes(e)
	routes.DeptRoutes(e)
	routes.ExtraMealRoutes(e)
	routes.Preference(e)
	routes.HolidayRoutes(e)

	log.Fatal(e.Start(fmt.Sprintf(":%s", config.LocalConfig.ServerPort)))
}
