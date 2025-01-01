package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"meal-management/pkg/containers"
	"meal-management/pkg/controllers"
)

func main() {
	fmt.Println("Hello World")
	//cronjobs.StartCronJob()
	controllers.StartCronJob()
	e := echo.New()
	containers.Serve(e)
}
