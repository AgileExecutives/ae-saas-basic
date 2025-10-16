package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system with multi-tenant support
type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username" binding:"required"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	FirstName    string         `json:"first_name" binding:"required"`
	LastName     string         `json:"last_name" binding:"required"`
	Role         string         `gorm:"not null;default:'user'" json:"role"`
	TenantID     uint           `gorm:"not null" json:"tenant_id"`
	// Tenant       Tenant         `gorm:"foreignKey:TenantID" json:"tenant,omitempty"` // Disabled for migration
	Active       bool           `gorm:"default:true" json:"active"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// UserResponse represents the API response structure for User
type UserResponse struct {
	ID        uint           `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Role      string         `json:"role"`
	TenantID  uint           `json:"tenant_id"`
	Tenant    TenantResponse `json:"tenant,omitempty"`
	Active    bool           `json:"active"`
	CreatedAt time.Time      `json:"created_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	response := UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		TenantID:  u.TenantID,
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
	}

	// Temporarily commented out due to migration issues
	// if u.Tenant.ID != 0 {
	// 	response.Tenant = u.Tenant.ToResponse()
	// }

	return response
}

// UserCreateRequest represents the request structure for creating a user
type UserCreateRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      string `json:"role"`
	TenantID  uint   `json:"tenant_id" binding:"required"`
}

// UserUpdateRequest represents the request structure for updating a user
type UserUpdateRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email" binding:"omitempty,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Active    *bool  `json:"active"`
}

// LoginRequest represents the login request structure
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response structure
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
