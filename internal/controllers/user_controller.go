package controllers

import (
	"net/http"

	"github.com/akhilnasimk/SS_backend/internal/dto"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(service services.UserService) UserController {
	return UserController{
		UserService: service,
	}
}

func (C *UserController) GetProfile(ctx *gin.Context) {

	// Extracting id
	userIDRaw, exists := ctx.Get("UserID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	//paresing the string to uuid
	userUUID, err := uuid.Parse(userIDRaw.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id inside token",
		})
		return
	}

	// using the service
	profile, err := C.UserService.GetProfile(userUUID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	//send the response if everythin goes right
	ctx.JSON(http.StatusOK, gin.H{
		"message": "profile fetched successfully",
		"data":    profile,
	})
}

func (C *UserController) UpdateProfile(ctx *gin.Context) {
	// Extracting id
	userIDRaw, exists := ctx.Get("UserID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	//paresing the string to uuid
	userUUID, err := uuid.Parse(userIDRaw.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id inside token",
		})
		return
	}

	var Req dto.UpdateProfileRequest

	if err := ctx.ShouldBindBodyWithJSON(&Req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := C.UserService.UpdateProfile(userUUID, &Req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "profile Updated ",
	})

}
