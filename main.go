package main

import "github.com/labstack/echo/v4"

func main() {
	e := echo.New()
	containers.Serve(e)
}
