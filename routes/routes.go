package routes

import (
	"KeuskupanLaboanBajo/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gorm.io/gorm"
)

func UserRoutes(e *echo.Echo, db *gorm.DB) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api/v1")

	userController := controllers.NewUserController(db)
	user := api.Group("/user")
	user.POST("", userController.Create)
	user.GET("", userController.List)
	user.GET("/:id", userController.Get)
	user.PUT("/id", userController.Update)
	user.DELETE("/:id", userController.Delete)
}
