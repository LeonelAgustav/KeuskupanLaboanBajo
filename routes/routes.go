package routes

import (
	"KeuskupanLaboanBajo/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api/v1")

	// User routes
	userController := controllers.NewUserController(db)
	users := api.Group("/users")
	users.POST("", userController.Create)
	users.GET("", userController.List)
	users.GET("/:id", userController.Get)
	users.PUT("/:id", userController.Update)
	users.DELETE("/:id", userController.Delete)

	// Paroki routes
	parokiController := controllers.NewParokiController(db)
	parokis := api.Group("/paroki")
	parokis.POST("", parokiController.Create)
	parokis.GET("", parokiController.List)
	parokis.GET("/:id", parokiController.Get)
	parokis.PUT("/:id", parokiController.Update)
	parokis.DELETE("/:id", parokiController.Delete)

	// Keuskupan routes
	keuskupanController := controllers.NewKeuskupanController(db)
	keuskupans := api.Group("/keuskupan")
	keuskupans.POST("", keuskupanController.Create)
	keuskupans.GET("", keuskupanController.List)
	keuskupans.GET("/:id", keuskupanController.Get)
	keuskupans.PUT("/:id", keuskupanController.Update)
	keuskupans.DELETE("/:id", keuskupanController.Delete)

	// Akun routes
	akunController := controllers.NewAkunController(db)
	akuns := api.Group("/akun")
	akuns.POST("", akunController.Create)
	akuns.GET("", akunController.List)
	akuns.GET("/:id", akunController.Get)
	akuns.PUT("/:id", akunController.Update)
	akuns.DELETE("/:id", akunController.Delete)
}