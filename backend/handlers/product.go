package handlers

import (
	"net/http"
	"strconv"

	"silverlord/backend/cache"
	"silverlord/backend/database"
	"silverlord/backend/models"
	"silverlord/backend/services"

	"github.com/labstack/echo/v4"
)

type productResponse struct {
	models.Product
	Price int `json:"price"`
}

// GET /api/v1/products
func GetProducts(c echo.Context) error {
	snap := cache.Get()
	if snap.XagUsd == 0 {
		return errServiceUnavailable(c, ErrPriceNotReady, "price not ready, please retry in a moment")
	}

	var products []models.Product
	if err := database.DB.Where("is_active = ?", true).Find(&products); err != nil {
		return errInternal(c, err.Error())
	}

	var cfg models.SystemConfig
	database.DB.ID(1).Get(&cfg)

	result := make([]productResponse, 0, len(products))
	for _, p := range products {
		result = append(result, productResponse{
			Product: p,
			Price:   services.CalcPrice(&p, &cfg),
		})
	}
	return c.JSON(http.StatusOK, result)
}

// GET /api/v1/products/:id
func GetProduct(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var p models.Product
	has, _ := database.DB.ID(id).Where("is_active = ?", true).Get(&p)
	if !has {
		return errNotFound(c, ErrProductNotFound, "product not found")
	}
	var cfg models.SystemConfig
	database.DB.ID(1).Get(&cfg)
	snap := cache.Get()
	price := 0
	if snap.XagUsd > 0 {
		price = services.CalcPrice(&p, &cfg)
	}
	return c.JSON(http.StatusOK, productResponse{Product: p, Price: price})
}

// GET /api/v1/spot-price
func GetSpotPrice(c echo.Context) error {
	return c.JSON(http.StatusOK, cache.Get())
}
