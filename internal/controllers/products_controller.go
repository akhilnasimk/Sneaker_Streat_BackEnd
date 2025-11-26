package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/akhilnasimk/SS_backend/internal/dto"
	"github.com/akhilnasimk/SS_backend/internal/helpers"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/akhilnasimk/SS_backend/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	// New filters
	categoryID := ctx.Query("category_id")
	search := ctx.Query("search")
	minPriceStr := ctx.Query("min_price")
	maxPriceStr := ctx.Query("max_price")

	var minPrice, maxPrice int64
	var _ = minPriceStr
	var _ = maxPriceStr

	if minPriceStr != "" {
		minPrice, _ = strconv.ParseInt(minPriceStr, 10, 64)
	}

	if maxPriceStr != "" {
		maxPrice, _ = strconv.ParseInt(maxPriceStr, 10, 64)
	}

	products, total, err := R.PService.GetAllProducts(page, limit, categoryID, search, minPrice, maxPrice)
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

func (R *ProductController) GetProductById(ctx *gin.Context) {
	id := ctx.Param("id")

	product, err := R.PService.GetProductById(id)

	if err != nil {
		ctx.JSON(400, response.Failure("failed to get the product", err.Error()))
		return
	}

	ctx.JSON(200, response.Success("product has fetched", product))
}

// Uploading product withn cloudinery
func (c *ProductController) UploadProduct(ctx *gin.Context) {
	// fmt.Println("Content-Type:", ctx.ContentType())

	// Use Gin's built-in method - DON'T call ParseMultipartForm manually
	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	priceStr := ctx.PostForm("price")
	stockStr := ctx.PostForm("stock_count")
	categoryIDStr := ctx.PostForm("category_id")

	// fmt.Printf("Parsed values - Name: '%s', Description: '%s', Price: '%s', Stock: '%s', Category: '%s'\n",
	// 	name, description, priceStr, stockStr, categoryIDStr)

	// Validation
	if name == "" || description == "" || priceStr == "" || stockStr == "" || categoryIDStr == "" {
		ctx.JSON(http.StatusBadRequest, response.Failure("all fields are required", nil))
		return
	}

	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("invalid price", nil))
		return
	}

	stockCount, err := strconv.Atoi(stockStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("invalid stock count", nil))
		return
	}

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("invalid category UUID", nil))
		return
	}

	// Get files using Gin's method
	form, err := ctx.MultipartForm()
	if err != nil {
		fmt.Printf("MultipartForm error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, response.Failure(fmt.Sprintf("failed to get files: %v", err), nil))
		return
	}

	files, exists := form.File["images"]
	if !exists || len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, response.Failure("at least one image file is required", nil))
		return
	}

	fmt.Printf("Received %d files\n", len(files))

	// Call service
	product, err := c.PService.CreateProduct(name, description, price, stockCount, categoryID, files)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Failure(fmt.Sprintf("failed to create product: %v", err), nil))
		return
	}

	resp := dto.ToProductResponse(product)
	ctx.JSON(http.StatusOK, response.Success("product uploaded successfully", resp))
}

func (c *ProductController) GetAllCategory(ctx *gin.Context) {

	categories, err := c.PService.GetAllCategory()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Failure("failed to fetch categories", err))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("categories fetched successfully", categories))
}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	// Parse product ID from URL
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("Invalid product ID", err.Error()))
		return
	}

	// Parse multipart form (32 MB)
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("Failed to parse form data", err.Error()))
		return
	}

	// Bind other form fields into DTO
	var req dto.UpdateProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("Invalid request data", err.Error()))
		return
	}
	// fmt.Println("update req (after bind):", req)

	//  Get multipart form (for files and raw values)
	form := ctx.Request.MultipartForm
	if form == nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("Multipart form not found", nil))
		return
	}

	// fmt.Printf("multipart form keys: %+v\n", form.Value)

	// Extract new_images files (if any)
	// Explicitly guard accesses to form.File to avoid staticcheck nil warnings
	if form.File != nil {
		if files, exists := form.File["new_images"]; exists {
			req.NewImages = files
			// fmt.Printf("new_images count: %d\n", len(files))
		}
	}

	var keepRaw []string
	hasPlain := false
	hasBracket := false

	// Guard form.Value as well (map might be nil)
	if form.Value != nil {
		if v, ok := form.Value["keep_images"]; ok {
			keepRaw = append(keepRaw, v...)
			hasPlain = true
		}
		if v, ok := form.Value["keep_images[]"]; ok {
			keepRaw = append(keepRaw, v...)
			hasBracket = true
		}
	}

	keepPresent := hasPlain || hasBracket

	// filter blank entries and trim spaces
	filtered := make([]string, 0, len(keepRaw))
	for _, s := range keepRaw {
		s = strings.TrimSpace(s)
		if s != "" {
			filtered = append(filtered, s)
		}
	}

	if !keepPresent {
		// key not present at all -> don't touch images
		req.KeepImages = nil
		fmt.Println("keep_images key not present -> will not modify images")
	} else {
		req.KeepImages = filtered
	}

	// Validate request  - check all the essenaila need - in helpers
	if err := helpers.ValidateUpdateProductRequest(req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("Validation failed", err.Error()))
		return
	}
	fmt.Println("final req before service call:", req)

	// Call service
	if err := c.PService.UpdateProduct(id, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Failure("Failed to update product", err.Error()))
		return
	}

	//  Success response
	ctx.JSON(http.StatusOK, response.Success("Product updated successfully", nil))
}
