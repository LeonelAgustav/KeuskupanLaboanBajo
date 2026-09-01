package controllers

import (
	"net/http"

	"KeuskupanLaboanBajo_BE/middleware"
	"KeuskupanLaboanBajo_BE/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

type RegisterRequest struct {
	Nama     string `json:"nama" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (a *AuthController) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{Nama: req.Nama, Email: req.Email, PasswordHash: string(hash)}
	if err := a.db.Create(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusConflict, "Email sudah terdaftar")
	}
	access, refresh, _ := middleware.GenerateToken(user.ID, user.Email)
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Registrasi berhasil",
		"data": map[string]string{
			"id": user.ID, "email": user.Email, "access_token": access, "refresh_token": refresh,
		},
	})
}

func (a *AuthController) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return err
	}
	var user models.User
	if err := a.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Email atau password salah")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Email atau password salah")
	}
	access, refresh, _ := middleware.GenerateToken(user.ID, user.Email)
	return c.JSON(http.StatusOK, map[string]string{
		"access_token": access, "refresh_token": refresh,
	})
}

func (a *AuthController) Refresh(c echo.Context) error {
	var body struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := c.Bind(&body); err != nil {
		return err
	}
	if err := c.Validate(&body); err != nil {
		return err
	}
	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(body.RefreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return middleware.JWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return echo.NewHTTPError(http.StatusUnauthorized, "Refresh token tidak valid")
	}
	access, refresh, _ := middleware.GenerateToken(claims.UserID, claims.Email)
	return c.JSON(http.StatusOK, map[string]string{"access_token": access, "refresh_token": refresh})
}
