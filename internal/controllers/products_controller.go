package controllers

import (
	"net/http"
	"strconv"

	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	PService services.ProductsService
}

func NewProductController(service services.ProductsService) ProductController {
	return ProductController{
		PService: service,
	}
}

func (R *ProductController) GetAllProducts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	products, total, err := R.PService.GetAllProducts(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}
