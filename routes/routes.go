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

	// Jenis routes
	jenisController := controllers.NewJenisController(db)
	jeniss := api.Group("/jenis")
	jeniss.POST("", jenisController.Create)
	jeniss.GET("", jenisController.List)
	jeniss.GET("/:id", jenisController.Get)
	jeniss.PUT("/:id", jenisController.Update)
	jeniss.DELETE("/:id", jenisController.Delete)

	// Pembatasan routes
	pembatasanController := controllers.NewPembatasanController(db)
	pembatasans := api.Group("/pembatasan")
	pembatasans.POST("", pembatasanController.Create)
	pembatasans.GET("", pembatasanController.List)
	pembatasans.GET("/:id", pembatasanController.Get)
	pembatasans.PUT("/:id", pembatasanController.Update)
	pembatasans.DELETE("/:id", pembatasanController.Delete)

	// Jurnal routes
	jurnalController := controllers.NewJurnalController(db)
	jurnals := api.Group("/jurnal")
	jurnals.POST("", jurnalController.Create)
	jurnals.GET("", jurnalController.List)
	jurnals.GET("/:id", jurnalController.Get)
	jurnals.PUT("/:id", jurnalController.Update)
	jurnals.DELETE("/:id", jurnalController.Delete)
	jurnals.GET("/:id/detil", jurnalController.ListDetil)
}