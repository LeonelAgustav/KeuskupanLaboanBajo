package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"KeuskupanLaboanBajo/dto"
	"KeuskupanLaboanBajo/models"
)

var validate = validator.New()

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

func (c *UserController) List(ctx echo.Context) error {
	var query dto.ListQuery
	if err := ctx.Bind(&query); err != nil {
		return err
	}
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	db := c.db.Model(&models.User{})

	if query.Search != "" {
		search := "%" + query.Search + "%"
		db = db.Where("nama ILIKE ? OR email ILIKE ?", search, search)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return err
	}

	offset := (query.Page - 1) * query.Limit
	orderBy := "id DESC"
	if query.Sort != "" {
		orderBy = query.Sort
	}

	var user []models.User
	if err := db.Order(orderBy).Offset(offset).Limit(query.Limit).Find(&user).Error; err != nil {
		return err
	}

	responses := make([]dto.UserResponse, len(user))
	for i, u := range user {
		responses[i] = dto.UserResponse{
			ID:        u.ID,
			Nama:      u.Nama,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}

	return ctx.JSON(http.StatusOK, dto.ListResponse{
		Data:       responses,
		Page:       query.Page,
		Limit:      query.Limit,
		Total:      total,
		TotalPages: int((total + int64(query.Limit) - 1) / int64(query.Limit)),
	})
}

func (c *UserController) Get(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID tidak valid")
	}

	var user models.User
	if err := c.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "User tidak ditemukan")
		}
		return err
	}

	return ctx.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID,
		Nama:      user.Nama,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (c *UserController) Create(ctx echo.Context) error {
	var req dto.UserCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	if err := validate.Struct(&req); err != nil {
		return err
	}

	user := models.User{
		Nama:  req.Nama,
		Email: req.Email,
	}

	if err := c.db.Create(&user).Error; err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, dto.CreateResponse{
		Message: "Data Berhasil di Buat",
		Data: dto.UserResponse{
			ID:        user.ID,
			Nama:      user.Nama,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

func (c *UserController) Update(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID tidak valid")
	}

	var req dto.UserUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}
	if err := validate.Struct(&req); err != nil {
		return err
	}

	var user models.User
	if err := c.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "User tidak ditemukan")
		}
		return err
	}

	if req.Nama != nil {
		user.Nama = *req.Nama
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	if err := c.db.Save(&user).Error; err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, dto.CreateResponse{
		Message: "Data Berhasil di Perbarui",
		Data: dto.UserResponse{
			ID:        user.ID,
			Nama:      user.Nama,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

func (c *UserController) Delete(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID tidak valid")
	}

	result := c.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "User tidak ditemukan")
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Data Berhasil di Hapus",
	})
}
