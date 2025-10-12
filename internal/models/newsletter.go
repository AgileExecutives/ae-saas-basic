package models

import (
	"time"

	"gorm.io/gorm"
)

// Newsletter represents a newsletter subscription
type Newsletter struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Status    string         `gorm:"default:'subscribed'" json:"status"` // subscribed, unsubscribed, pending
	Source    string         `gorm:"default:'website'" json:"source"`    // website, api, import
	Tags      string         `json:"tags"`                               // comma-separated tags
	OptInDate *time.Time     `json:"opt_in_date,omitempty"`
	IPAddress string         `json:"ip_address,omitempty"`
	UserAgent string         `json:"user_agent,omitempty"`
	Active    bool           `gorm:"default:true" json:"active"`
}

// TableName specifies the table name for Newsletter
func (Newsletter) TableName() string {
	return "newsletters"
}

// NewsletterSubscribeRequest represents the request for newsletter subscription
type NewsletterSubscribeRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Source    string `json:"source,omitempty"`
	Tags      string `json:"tags,omitempty"`
}

// NewsletterUnsubscribeRequest represents the request for newsletter unsubscription
type NewsletterUnsubscribeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// NewsletterResponse represents the API response structure for Newsletter
type NewsletterResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Status    string    `json:"status"`
	Source    string    `json:"source"`
	Tags      string    `json:"tags"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

// ContactFormRequest represents a public contact form submission
type ContactFormRequest struct {
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Phone        string `json:"phone"`
	Subject      string `json:"subject" binding:"required"`
	Message      string `json:"message" binding:"required"`
	Company      string `json:"company"`
	Newsletter   bool   `json:"newsletter"`           // whether to subscribe to newsletter
	Source       string `json:"source,omitempty"`     // form source identifier
	CaptchaToken string `json:"captcha_token"`        // for spam protection
}

// ContactFormResponse represents the response for contact form submission
type ContactFormResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	ContactID     uint   `json:"contact_id,omitempty"`
	NewsletterID  uint   `json:"newsletter_id,omitempty"`
	SubscribedTo  bool   `json:"subscribed_to_newsletter,omitempty"`
}