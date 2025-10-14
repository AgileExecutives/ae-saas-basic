package handlers

import (
	"net/http"

	"github.com/ae-saas-basic/ae-saas-basic/internal/models"
	"github.com/ae-saas-basic/ae-saas-basic/pkg/auth"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Login authenticates a user and returns a JWT token
// @Summary Login user
// @Description Authenticate user with username/email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.APIResponse{data=models.LoginResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	var user models.User
	// Find user by username or email
	if err := h.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Invalid credentials", "User not found"))
		return
	}

	// Check if user is active
	if !user.Active {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Account disabled", "User account is not active"))
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("Invalid credentials", "Password mismatch"))
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.TenantID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to generate token", err.Error()))
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Login successful", response))
}

// Logout blacklists the current JWT token
// @Summary Logout user
// @Description Blacklist the current JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get token from context (set by auth middleware)
	tokenString, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("No token found", "Token not provided"))
		return
	}

	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}

	user := userInterface.(*models.User)

	// Parse token to get JTI and expiration
	tokenID, expiresAt, err := auth.ParseTokenClaims(tokenString.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid token", err.Error()))
		return
	}

	// Add token to blacklist
	blacklistEntry := models.TokenBlacklist{
		TokenID:   tokenID,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
		Reason:    "User logout",
	}

	if err := h.db.Create(&blacklistEntry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to logout", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Logout successful", nil))
}

// Register creates a new user account
// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.UserCreateRequest true "User registration data"
// @Success 201 {object} models.APIResponse{data=models.UserResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	// Check if username or email already exists
	var existingUser models.User
	if err := h.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, models.ErrorResponseFunc("User already exists", "Username or email already taken"))
		return
	}

	// Check if tenant exists
	var tenant models.Tenant
	if err := h.db.First(&tenant, req.TenantID).Error; err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Tenant not found", "Invalid tenant ID"))
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to hash password", err.Error()))
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "user"
	}

	// Create user
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
		TenantID:     req.TenantID,
		Active:       true,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to create user", err.Error()))
		return
	}

	// Load tenant for response
	h.db.Preload("Tenant").First(&user, user.ID)

	c.JSON(http.StatusCreated, models.SuccessResponse("User created successfully", user.ToResponse()))
}

// ChangePassword changes the user's password
// @Summary Change password
// @Description Change user password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{current_password=string,new_password=string} true "Password change data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid request", err.Error()))
		return
	}

	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}

	user := userInterface.(*models.User)

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseFunc("Invalid current password", "Current password is incorrect"))
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to hash password", err.Error()))
		return
	}

	// Update password
	if err := h.db.Model(user).Update("password_hash", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseFunc("Failed to update password", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Password changed successfully", nil))
}

// Me returns the current user information
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.UserResponse}
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseFunc("User not found", "User not authenticated"))
		return
	}

	user := userInterface.(*models.User)

	// Preload tenant
	h.db.Preload("Tenant").First(user, user.ID)

	c.JSON(http.StatusOK, models.SuccessResponse("User retrieved successfully", user.ToResponse()))
}
