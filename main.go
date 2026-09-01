package main

import (
	"fmt"

	"KeuskupanLaboanBajo_BE/config"
	"KeuskupanLaboanBajo_BE/middleware"
	"KeuskupanLaboanBajo_BE/routes"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	db := config.ConnectDB()

	e := echo.New()

	e.HTTPErrorHandler = middleware.ErrorHandler
	e.Validator = &middleware.CustomValidator{Validator: validator.New()}

	routes.RegisterRoutes(e, db)

	fmt.Println("starting web server at http://localhost:8080/")
	e.Logger.Fatal(e.Start(":8080"))
}
