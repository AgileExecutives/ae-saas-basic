package models

import (
	"time"

	"gorm.io/gorm"
)

// Organization represents a tenant in the multi-tenant system
// Simplified version containing only ID and Name for tenant identification
type Organization struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name      string         `gorm:"not null;unique" json:"name" binding:"required"`
	Slug      string         `gorm:"not null;unique" json:"slug" binding:"required"`
}

// TableName specifies the table name for Organization
func (Organization) TableName() string {
	return "organizations"
}

// OrganizationResponse represents the API response structure for Organization
type OrganizationResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ToResponse converts Organization to OrganizationResponse
func (o *Organization) ToResponse() OrganizationResponse {
	return OrganizationResponse{
		ID:   o.ID,
		Name: o.Name,
		Slug: o.Slug,
	}
}
