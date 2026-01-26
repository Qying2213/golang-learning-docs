package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qinyang/taskmanager/config"
	"github.com/qinyang/taskmanager/database"
	"github.com/qinyang/taskmanager/models"
	"github.com/qinyang/taskmanager/utils"
)

type AuthHandler struct {
	Config *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{Config: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.GetDB().Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "User with this email already exists")
		return
	}

	// Create new user
	user := models.User{
		Email: req.Email,
		Name:  req.Name,
	}

	if err := user.HashPassword(req.Password); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, h.Config.JWT.Secret, h.Config.JWT.ExpirationHours)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  user,
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Find user by email
	var user models.User
	if err := database.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, h.Config.JWT.Secret, h.Config.JWT.ExpirationHours)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  user,
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var user models.User
	if err := database.GetDB().First(&user, "id = ?", userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}
