# Integration Example for Unburdy Backend

This document shows how to integrate the `ae-saas-basic` module into the existing Unburdy backend.

## Step 1: Add Module Dependency

In the `unburdy-backend` project:

```bash
cd /path/to/unburdy-backend
go get github.com/ae-saas-basic/ae-saas-basic
```

## Step 2: Update Unburdy's main.go

Replace the existing auth, user, customer, and plan functionality with the SaaS module:

```go
// main.go in unburdy-backend
package main

import (
    "log"

    // AE SaaS Basic imports
    "github.com/ae-saas-basic/ae-saas-basic/internal/config"
    "github.com/ae-saas-basic/ae-saas-basic/internal/database"
    saasRouter "github.com/ae-saas-basic/ae-saas-basic/internal/router"
    "github.com/ae-saas-basic/ae-saas-basic/pkg/auth"
    
    // Unburdy-specific imports
    "github.com/your-org/unburdy-backend/internal/handlers"
    "github.com/your-org/unburdy-backend/internal/middleware"
    "github.com/gin-gonic/gin"
)

func main() {
    // Load SaaS basic configuration
    cfg := config.Load()
    
    // Set JWT secret
    auth.SetJWTSecret(cfg.JWT.Secret)
    
    // Connect to shared database
    db, err := database.Connect(cfg.Database)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    
    // Run SaaS migrations (organizations, users, plans, customers, etc.)
    if err := database.Migrate(db); err != nil {
        log.Fatal("Failed to migrate SaaS tables:", err)
    }
    
    // Run Unburdy-specific migrations (clients, therapies, sessions, etc.)
    if err := migrateUnburdyTables(db); err != nil {
        log.Fatal("Failed to migrate Unburdy tables:", err)
    }
    
    // Seed SaaS basic data
    if err := database.Seed(db); err != nil {
        log.Fatal("Failed to seed database:", err)
    }
    
    // Setup SaaS router with all basic endpoints
    router := saasRouter.SetupRouter(db)
    
    // Add Unburdy-specific routes
    unburdyAPI := router.Group("/api/v1")
    unburdyAPI.Use(middleware.AuthMiddleware(db)) // Use SaaS auth middleware
    {
        // Client management (therapy-specific)
        clients := unburdyAPI.Group("/clients")
        {
            clients.GET("", handlers.GetClients)
            clients.GET("/:id", handlers.GetClient)
            clients.POST("", handlers.CreateClient)
            clients.PUT("/:id", handlers.UpdateClient)
            clients.DELETE("/:id", handlers.DeleteClient)
        }
        
        // Therapy management
        therapies := unburdyAPI.Group("/therapies")
        {
            therapies.GET("", handlers.GetTherapies)
            therapies.GET("/:id", handlers.GetTherapy)
            therapies.POST("", handlers.CreateTherapy)
            therapies.PUT("/:id", handlers.UpdateTherapy)
            therapies.DELETE("/:id", handlers.DeleteTherapy)
        }
        
        // Therapy sessions
        sessions := unburdyAPI.Group("/sessions")
        {
            sessions.GET("", handlers.GetSessions)
            sessions.POST("", handlers.CreateSession)
            sessions.PUT("/:id", handlers.UpdateSession)
        }
        
        // Provider management (therapy providers)
        providers := unburdyAPI.Group("/providers")
        {
            providers.GET("", handlers.GetProviders)
            providers.POST("", handlers.CreateProvider)
        }
        
        // Calendar and scheduling
        calendar := unburdyAPI.Group("/calendar")
        {
            calendar.GET("/events", handlers.GetCalendarEvents)
            calendar.POST("/events", handlers.CreateCalendarEvent)
        }
    }
    
    // Start server
    addr := cfg.Server.Host + ":" + cfg.Server.Port
    log.Printf("Starting Unburdy Backend with SaaS integration on %s", addr)
    router.Run(addr)
}
```

## Step 3: Update Unburdy Models

Update Unburdy-specific models to reference SaaS models:

```go
// models/client.go in unburdy-backend
package models

import (
    "time"
    "gorm.io/gorm"
    saas "github.com/ae-saas-basic/ae-saas-basic/internal/models"
)

// Client represents a therapy client (extends SaaS Contact)
type Client struct {
    ID             uint           `gorm:"primarykey" json:"id"`
    CreatedAt      time.Time      `json:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at"`
    DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
    
    // Reference to SaaS Contact for basic info
    ContactID      uint           `gorm:"not null" json:"contact_id"`
    Contact        saas.Contact   `gorm:"foreignKey:ContactID" json:"contact,omitempty"`
    
    // Therapy-specific fields
    DateOfBirth    *time.Time     `json:"date_of_birth"`
    Gender         string         `json:"gender"`
    Diagnosis      string         `json:"diagnosis"`
    Medications    string         `gorm:"type:text" json:"medications"`
    Allergies      string         `gorm:"type:text" json:"allergies"`
    
    // Multi-tenant support via organization
    OrganizationID uint           `gorm:"not null" json:"organization_id"`
    Organization   saas.Organization `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
    
    // Relationships
    Therapies      []Therapy      `gorm:"foreignKey:ClientID" json:"therapies,omitempty"`
    Contacts       []saas.Contact `gorm:"many2many:client_contacts;" json:"contacts,omitempty"`
    
    Active         bool           `gorm:"default:true" json:"active"`
}
```

## Step 4: Update Unburdy Handlers

Update handlers to use SaaS authentication and user context:

```go
// handlers/client.go in unburdy-backend
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    saas "github.com/ae-saas-basic/ae-saas-basic/internal/models"
    "github.com/your-org/unburdy-backend/internal/models"
)

func GetClients(c *gin.Context) {
    // Get authenticated user from SaaS middleware
    userInterface, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    
    user := userInterface.(*saas.User)
    
    // Use organization ID for tenant isolation
    var clients []models.Client
    if err := db.Where("organization_id = ?", user.OrganizationID).
        Preload("Contact").
        Preload("Therapies").
        Find(&clients).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"data": clients})
}
```

## Step 5: Benefits After Integration

### What You Get from AE SaaS Basic:

1. **Complete Authentication System**
   - JWT token management
   - User registration/login
   - Password change functionality
   - Token blacklisting for secure logout

2. **Multi-tenant Architecture** 
   - Organization-level data separation
   - Automatic tenant isolation in middleware

3. **User Management**
   - Role-based access control
   - User settings and preferences

4. **Customer & Billing**
   - Customer management
   - Subscription plans
   - Basic billing structure

5. **Communication**
   - Email system with tracking
   - Contact management

6. **System Health**
   - Health check endpoints
   - API monitoring

### What Stays in Unburdy:

1. **Therapy-Specific Logic**
   - Client management (therapy clients)
   - Therapy plans and sessions
   - Therapy reports and discharge summaries
   - Provider management (therapy providers)
   - Calendar and scheduling
   - Therapy-specific business logic

## Step 6: Migration Strategy

1. **Phase 1**: Deploy AE SaaS Basic alongside existing Unburdy
2. **Phase 2**: Migrate user authentication to use SaaS module
3. **Phase 3**: Migrate customer and plan management
4. **Phase 4**: Update all handlers to use SaaS middleware
5. **Phase 5**: Remove redundant code from Unburdy

## API Endpoints After Integration

### Provided by AE SaaS Basic:
- `/api/v1/auth/*` - Authentication
- `/api/v1/customers/*` - Customer management
- `/api/v1/contacts/*` - Contact management
- `/api/v1/plans/*` - Plan management
- `/api/v1/emails/*` - Email management
- `/api/v1/user-settings/*` - User settings
- `/api/v1/health` - Health check

### Remains in Unburdy:
- `/api/v1/clients/*` - Therapy clients
- `/api/v1/therapies/*` - Therapy management
- `/api/v1/sessions/*` - Therapy sessions
- `/api/v1/providers/*` - Therapy providers
- `/api/v1/calendar/*` - Calendar and scheduling
- `/api/v1/reports/*` - Therapy reports

This integration provides a clean separation of concerns while maintaining all existing Unburdy functionality.

## Step 7: Static Files Integration

### Overview

The AE SaaS Basic module includes a comprehensive static files system with email templates, PDF templates, CSS assets, and image resources that can be integrated into your Unburdy backend.

### Static Directory Structure

```
ae-saas-basic/statics/
├── assets/
│   └── css/
│       ├── ae-saas.css      # Main application styles
│       └── email.css        # Email-specific styles
├── email_templates/
│   ├── welcome.html         # New user welcome email
│   ├── password_reset.html  # Password reset email
│   ├── email_verification.html # Email verification
│   ├── user_invitation.html # Team invitations
│   └── subscription_notification.html # Billing emails
├── pdf_templates/
│   ├── invoice.html         # Professional invoice template
│   └── report.html          # Comprehensive report template
└── images/
    ├── logo.svg            # Company logo (compact)
    ├── logo-full.svg       # Full logo with text
    └── favicon.svg         # Website favicon
```

### Integration Options

#### Option 1: Copy Static Files to Unburdy

```bash
# Copy the entire statics directory to your Unburdy project
cp -r /path/to/ae-saas-basic/statics /path/to/unburdy-backend/

# Or copy specific directories as needed
cp -r /path/to/ae-saas-basic/statics/email_templates /path/to/unburdy-backend/statics/
cp -r /path/to/ae-saas-basic/statics/pdf_templates /path/to/unburdy-backend/statics/
```

#### Option 2: Reference SaaS Module Statics

In your Unburdy configuration, reference the SaaS module's static files:

```go
// config/config.go in unburdy-backend
package config

type Config struct {
    // ... existing config
    StaticFiles struct {
        // Use SaaS module templates
        EmailTemplatesDir string `env:"EMAIL_TEMPLATES_DIR" default:"./vendor/ae-saas-basic/statics/email_templates"`
        PDFTemplatesDir   string `env:"PDF_TEMPLATES_DIR" default:"./vendor/ae-saas-basic/statics/pdf_templates"`
        StaticFilesDir    string `env:"STATIC_FILES_DIR" default:"./vendor/ae-saas-basic/statics"`
        
        // Custom Unburdy templates
        UnburdyTemplatesDir string `env:"UNBURDY_TEMPLATES_DIR" default:"./statics/unburdy_templates"`
    }
}
```

### Email Template Integration

#### Using SaaS Email Templates in Unburdy

```go
// handlers/email.go in unburdy-backend
package handlers

import (
    "html/template"
    "path/filepath"
    
    saasModels "github.com/ae-saas-basic/ae-saas-basic/internal/models"
)

type UnburdyEmailHandler struct {
    emailTemplatesDir string
    emailHandler     *saasModels.EmailHandler
}

// Send welcome email using SaaS template
func (h *UnburdyEmailHandler) SendWelcomeEmail(userID uint) error {
    // Load SaaS welcome template
    templatePath := filepath.Join(h.emailTemplatesDir, "welcome.html")
    tmpl, err := template.ParseFiles(templatePath)
    if err != nil {
        return err
    }
    
    // Prepare template data
    data := struct {
        FirstName        string
        LastName         string
        Username         string
        Email            string
        OrganizationName string
        PlanName         string
        LoginURL         string
    }{
        FirstName:        user.FirstName,
        LastName:         user.LastName,
        Username:         user.Username,
        Email:            user.Email,
        OrganizationName: user.Organization.Name,
        PlanName:         user.Customer.Plan.Name,
        LoginURL:         "https://unburdy.com/login",
    }
    
    // Render template
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return err
    }
    
    // Send using SaaS email system
    return h.emailHandler.SendEmail(saasModels.SendEmailRequest{
        ToEmail:  user.Email,
        Subject:  "Welcome to Unburdy!",
        Body:     buf.String(),
        HTMLBody: buf.String(),
    })
}

// Send therapy-specific emails using custom Unburdy templates
func (h *UnburdyEmailHandler) SendTherapyReminder(clientID uint, sessionDate string) error {
    // Use Unburdy-specific template for therapy reminders
    templatePath := filepath.Join(h.unburdyTemplatesDir, "therapy_reminder.html")
    // ... therapy-specific email logic
}
```

#### Custom Unburdy Email Templates

Create therapy-specific templates in your Unburdy project:

```html
<!-- statics/unburdy_templates/therapy_reminder.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Therapy Session Reminder</title>
    <!-- Use SaaS CSS for consistent styling -->
    <link rel="stylesheet" href="/static/assets/css/email.css">
</head>
<body>
    <div class="email-wrapper">
        <div class="email-container">
            <div class="email-header">
                <div class="email-logo">Unburdy</div>
                <h1 class="email-title">Session Reminder</h1>
            </div>
            
            <div class="email-content">
                <p>Hi {{.ClientName}},</p>
                
                <p>This is a reminder of your upcoming therapy session:</p>
                
                <div class="email-info-box">
                    <h3>Session Details:</h3>
                    <ul>
                        <li><strong>Date:</strong> {{.SessionDate}}</li>
                        <li><strong>Time:</strong> {{.SessionTime}}</li>
                        <li><strong>Provider:</strong> {{.ProviderName}}</li>
                        <li><strong>Type:</strong> {{.SessionType}}</li>
                        <li><strong>Location:</strong> {{.Location}}</li>
                    </ul>
                </div>
                
                <center>
                    <a href="{{.RescheduleURL}}" class="email-button">Reschedule if Needed</a>
                </center>
            </div>
        </div>
    </div>
</body>
</html>
```

### PDF Template Integration

#### Using SaaS Invoice Template for Unburdy Billing

```go
// handlers/billing.go in unburdy-backend
package handlers

import (
    "bytes"
    "html/template"
    "path/filepath"
    
    "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func (h *BillingHandler) GenerateTherapyInvoice(customerID, invoiceID uint) ([]byte, error) {
    // Use SaaS invoice template
    templatePath := filepath.Join(h.pdfTemplatesDir, "invoice.html")
    tmpl, err := template.ParseFiles(templatePath)
    if err != nil {
        return nil, err
    }
    
    // Prepare invoice data with therapy-specific line items
    data := InvoiceData{
        InvoiceNumber:    fmt.Sprintf("UNBURDY-%d", invoiceID),
        InvoiceDate:      time.Now().Format("2006-01-02"),
        DueDate:          time.Now().AddDate(0, 0, 30).Format("2006-01-02"),
        CompanyName:      "Unburdy Therapy Services",
        CompanyAddress:   config.Get().Company.Address,
        CustomerName:     customer.Name,
        OrganizationName: customer.Organization.Name,
        LineItems: []LineItem{
            {
                Name:        "Therapy Sessions",
                Description: fmt.Sprintf("%d sessions in %s", sessionCount, period),
                Quantity:    sessionCount,
                UnitPrice:   sessionRate,
                Amount:      totalAmount,
                Currency:    "EUR",
            },
            // Add more therapy-specific line items
        },
        // ... other invoice fields
    }
    
    // Render HTML
    var htmlBuf bytes.Buffer
    if err := tmpl.Execute(&htmlBuf, data); err != nil {
        return nil, err
    }
    
    // Convert to PDF
    pdfg, err := wkhtmltopdf.NewPDFGenerator()
    if err != nil {
        return nil, err
    }
    
    pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(htmlBuf.Bytes())))
    
    return pdfg.ToBytes()
}
```

### Static File Serving Integration

#### Method 1: Use SaaS Static Handler

```go
// main.go in unburdy-backend
import (
    saasHandlers "github.com/ae-saas-basic/ae-saas-basic/internal/handlers"
)

func setupRouter(db *gorm.DB) *gin.Engine {
    router := gin.Default()
    
    // Use SaaS static file handler
    saasStaticHandler := saasHandlers.NewStaticHandler("./vendor/ae-saas-basic/statics")
    
    // SaaS static routes
    router.GET("/static/saas/*path", saasStaticHandler.ServeAsset)
    router.GET("/api/v1/logo", saasStaticHandler.ServeLogo)
    
    // Unburdy-specific static files
    router.Static("/static/unburdy", "./statics/unburdy")
    
    return router
}
```

#### Method 2: Combine Static Directories

```go
// handlers/static.go in unburdy-backend
type UnburdyStaticHandler struct {
    saasStaticDir    string
    unburdyStaticDir string
}

func (h *UnburdyStaticHandler) ServeAsset(c *gin.Context) {
    path := c.Param("path")
    
    // Try Unburdy-specific assets first
    unburdyPath := filepath.Join(h.unburdyStaticDir, path)
    if _, err := os.Stat(unburdyPath); err == nil {
        c.File(unburdyPath)
        return
    }
    
    // Fall back to SaaS assets
    saasPath := filepath.Join(h.saasStaticDir, path)
    if _, err := os.Stat(saasPath); err == nil {
        c.File(saasPath)
        return
    }
    
    c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
}
```

### Environment Configuration

Update your Unburdy `.env` file to include SaaS static configurations:

```bash
# Unburdy .env file

# SaaS Basic Module Configuration
SAAS_EMAIL_TEMPLATES_DIR=./vendor/ae-saas-basic/statics/email_templates
SAAS_PDF_TEMPLATES_DIR=./vendor/ae-saas-basic/statics/pdf_templates
SAAS_STATIC_DIR=./vendor/ae-saas-basic/statics

# Unburdy-specific templates
UNBURDY_TEMPLATES_DIR=./statics/unburdy_templates
UNBURDY_STATIC_DIR=./statics/unburdy

# Branding (override SaaS defaults for Unburdy)
COMPANY_NAME=Unburdy Therapy Services
COMPANY_EMAIL=contact@unburdy.com
COMPANY_WEBSITE=https://unburdy.com
SUPPORT_EMAIL=support@unburdy.com
```

### Customization for Unburdy Branding

#### Override Logo and Branding

```bash
# Replace SaaS logos with Unburdy branding
cp ./unburdy-assets/logo.svg ./vendor/ae-saas-basic/statics/images/logo.svg
cp ./unburdy-assets/logo-full.svg ./vendor/ae-saas-basic/statics/images/logo-full.svg

# Or create symbolic links
ln -sf $(pwd)/unburdy-assets/logo.svg ./vendor/ae-saas-basic/statics/images/logo.svg
```

#### Customize CSS Variables

```css
/* statics/unburdy/css/unburdy-theme.css */
:root {
    /* Override SaaS Basic theme with Unburdy colors */
    --primary-color: #2c5aa0;    /* Unburdy blue */
    --secondary-color: #5bc0de;  /* Unburdy light blue */
    --success-color: #5cb85c;    /* Keep green */
    --warning-color: #f0ad4e;    /* Unburdy orange */
    --danger-color: #d9534f;     /* Keep red */
    
    /* Unburdy-specific brand colors */
    --unburdy-primary: #2c5aa0;
    --unburdy-accent: #ff6b6b;
    --unburdy-text: #2c3e50;
}

/* Import SaaS base styles */
@import url('/static/saas/assets/css/ae-saas.css');

/* Unburdy-specific overrides */
.btn-primary {
    background-color: var(--unburdy-primary);
    border-color: var(--unburdy-primary);
}

.navbar-brand {
    color: var(--unburdy-primary);
}
```

### Template Data Mapping

Create a mapping service to convert between Unburdy and SaaS data structures:

```go
// services/template_mapper.go in unburdy-backend
package services

import (
    saasModels "github.com/ae-saas-basic/ae-saas-basic/internal/models"
    unburdyModels "github.com/your-org/unburdy-backend/internal/models"
)

type TemplateMapper struct{}

// Map Unburdy client to SaaS customer for templates
func (tm *TemplateMapper) ClientToCustomer(client *unburdyModels.Client) saasModels.Customer {
    return saasModels.Customer{
        Name:           fmt.Sprintf("%s %s", client.FirstName, client.LastName),
        Email:          client.Email,
        Phone:          client.Phone,
        Street:         client.Address,
        City:           client.City,
        Country:        client.Country,
        OrganizationID: client.OrganizationID,
        // Map other relevant fields
    }
}

// Map therapy session to line item for invoicing
func (tm *TemplateMapper) SessionToLineItem(session *unburdyModels.TherapySession) saasModels.LineItem {
    return saasModels.LineItem{
        Name:        "Therapy Session",
        Description: fmt.Sprintf("%s session with %s", session.Type, session.Provider.Name),
        Quantity:    1,
        UnitPrice:   session.Rate.String(),
        Amount:      session.Rate.String(),
        Currency:    "EUR",
    }
}
```

### Testing Static Integration

```bash
# Test SaaS logo access
curl http://localhost:8080/api/v1/logo?format=full

# Test combined static serving
curl http://localhost:8080/static/saas/assets/css/ae-saas.css
curl http://localhost:8080/static/unburdy/css/unburdy-theme.css

# Test email template preview (with auth)
curl -H "Authorization: Bearer $JWT_TOKEN" \
     http://localhost:8080/api/v1/static/templates/email/welcome.html
```

### Migration Checklist

- [ ] Copy or link SaaS static files to Unburdy project
- [ ] Update environment configuration
- [ ] Integrate email templates with Unburdy handlers
- [ ] Customize branding (logos, colors, company info)
- [ ] Create Unburdy-specific templates for therapy features
- [ ] Test static file serving and template rendering
- [ ] Update build and deployment scripts
- [ ] Document custom template variables for team

This integration approach allows you to leverage the professional templates and styling from AE SaaS Basic while maintaining Unburdy's unique therapy-focused functionality and branding.

---

## PDF Generation Integration

The AE SaaS Basic module provides a comprehensive, generalized PDF generation system that can be integrated into any project, including Unburdy backend.

### Overview

The PDF generation system includes:
- **Template-based PDF generation** with HTML templates
- **Dynamic data injection** for personalized documents
- **Multiple output formats** (download, save, stream)
- **Configurable settings** (page size, margins, quality)
- **REST API endpoints** for easy integration
- **Template validation** and error handling
- **Common document types** (invoices, reports, certificates)

### System Requirements

Install wkhtmltopdf on your system:

```bash
# Ubuntu/Debian
sudo apt-get install wkhtmltopdf

# macOS
brew install wkhtmltopdf

# CentOS/RHEL
sudo yum install wkhtmltopdf

# Docker
FROM ubuntu:20.04
RUN apt-get update && apt-get install -y wkhtmltopdf
```

### Integration in Unburdy Backend

#### 1. Add PDF Configuration

Update your `.env` file:

```bash
# PDF Generation Configuration
PDF_TEMPLATE_DIR=./statics/templates/pdf
PDF_OUTPUT_DIR=./output/pdf
PDF_PAGE_SIZE=A4
PDF_ORIENTATION=Portrait
PDF_MARGIN_TOP=1cm
PDF_MARGIN_RIGHT=1cm
PDF_MARGIN_BOTTOM=1cm
PDF_MARGIN_LEFT=1cm
PDF_QUALITY=80
PDF_ENABLE_JS=true
PDF_LOAD_TIMEOUT=30
PDF_MAX_FILE_SIZE=52428800
```

#### 2. Initialize PDF Service in Unburdy

```go
// In your Unburdy main.go or router setup
package main

import (
    saasConfig "github.com/ae-saas-basic/ae-saas-basic/internal/config"
    "github.com/ae-saas-basic/ae-saas-basic/internal/services"
    "github.com/ae-saas-basic/ae-saas-basic/internal/handlers"
)

func setupPDFService(cfg saasConfig.Config) *handlers.PDFHandler {
    // Create PDF service configuration
    pdfServiceConfig := &services.PDFConfig{
        PageSize:     cfg.PDF.PageSize,
        Orientation:  cfg.PDF.Orientation,
        MarginTop:    cfg.PDF.MarginTop,
        MarginRight:  cfg.PDF.MarginRight,
        MarginBottom: cfg.PDF.MarginBottom,
        MarginLeft:   cfg.PDF.MarginLeft,
        Quality:      cfg.PDF.Quality,
        EnableJS:     cfg.PDF.EnableJS,
        LoadTimeout:  cfg.PDF.LoadTimeout,
        Headers:      make(map[string]string),
    }
    
    // Initialize PDF service
    pdfService := services.NewPDFService(
        cfg.PDF.TemplateDir,
        cfg.PDF.OutputDir,
        pdfServiceConfig,
    )
    
    return handlers.NewPDFHandler(pdfService)
}

// In your router setup
func setupRouter(db *gorm.DB) *gin.Engine {
    router := gin.Default()
    
    // Load configuration
    cfg := saasConfig.Load()
    
    // Initialize PDF handler
    pdfHandler := setupPDFService(cfg)
    
    // Add PDF routes
    api := router.Group("/api/v1")
    pdf := api.Group("/pdf")
    pdf.Use(middleware.AuthMiddleware(db)) // Add auth if needed
    {
        // Template management
        pdf.GET("/templates", pdfHandler.ListTemplates)
        pdf.GET("/templates/:template", pdfHandler.GetTemplateInfo)
        pdf.POST("/templates/:template/preview", pdfHandler.PreviewTemplate)
        
        // PDF generation
        pdf.POST("/generate", pdfHandler.GeneratePDF)
        pdf.POST("/generate/html", pdfHandler.GeneratePDFFromHTML)
        pdf.POST("/generate/stream", pdfHandler.StreamPDF)
        
        // Configuration
        pdf.GET("/config", pdfHandler.GetPDFConfig)
    }
    
    return router
}
```

#### 3. Create Therapy-Specific PDF Templates

Create PDF templates in `./statics/templates/pdf/`:

**Therapy Session Report** (`session_report.html`):
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Therapy Session Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { border-bottom: 2px solid #333; padding-bottom: 20px; }
        .logo { float: left; }
        .company-info { float: right; text-align: right; }
        .clear { clear: both; }
        .session-info { margin: 30px 0; }
        .notes { margin: 20px 0; }
        .signature { margin-top: 50px; }
    </style>
</head>
<body>
    <div class="header">
        <div class="logo">
            {{if .Company.Logo}}
                <img src="{{.Company.Logo}}" alt="Logo" height="60">
            {{end}}
        </div>
        <div class="company-info">
            <h2>{{.Company.Name}}</h2>
            <p>{{.Company.Address}}<br>
               {{.Company.City}}, {{.Company.State}} {{.Company.ZipCode}}<br>
               {{.Company.Phone}}</p>
        </div>
        <div class="clear"></div>
    </div>
    
    <h1>Therapy Session Report</h1>
    
    <div class="session-info">
        <table>
            <tr><td><strong>Client:</strong></td><td>{{.Client.Name}}</td></tr>
            <tr><td><strong>Therapist:</strong></td><td>{{.Therapist.Name}}</td></tr>
            <tr><td><strong>Session Date:</strong></td><td>{{.SessionDate}}</td></tr>
            <tr><td><strong>Duration:</strong></td><td>{{.Duration}}</td></tr>
            <tr><td><strong>Session Type:</strong></td><td>{{.SessionType}}</td></tr>
        </table>
    </div>
    
    <div class="notes">
        <h3>Session Notes</h3>
        <p>{{.Notes}}</p>
        
        {{if .Goals}}
        <h3>Goals Addressed</h3>
        <ul>
        {{range .Goals}}
            <li>{{.}}</li>
        {{end}}
        </ul>
        {{end}}
        
        {{if .NextSteps}}
        <h3>Next Steps</h3>
        <p>{{.NextSteps}}</p>
        {{end}}
    </div>
    
    <div class="signature">
        <p>Therapist Signature: ___________________________</p>
        <p>Date: {{.GeneratedDate}}</p>
    </div>
</body>
</html>
```

**Client Invoice** (`therapy_invoice.html`):
```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Therapy Invoice</title>
    <link href="{{.BaseURL}}/static/css/pdf-styles.css" rel="stylesheet">
</head>
<body>
    <div class="invoice-header">
        <div class="logo">
            {{if .Company.Logo}}<img src="{{.Company.Logo}}" alt="Logo">{{end}}
        </div>
        <div class="company-details">
            <h1>{{.Company.Name}}</h1>
            <p>{{.Company.Address}}</p>
        </div>
    </div>
    
    <div class="invoice-info">
        <h2>Invoice #{{.InvoiceNumber}}</h2>
        <p><strong>Date:</strong> {{.InvoiceDate}}</p>
        <p><strong>Due Date:</strong> {{.DueDate}}</p>
    </div>
    
    <div class="client-info">
        <h3>Bill To:</h3>
        <p>{{.Client.Name}}<br>
           {{.Client.Address}}</p>
    </div>
    
    <table class="services-table">
        <thead>
            <tr>
                <th>Description</th>
                <th>Date</th>
                <th>Duration</th>
                <th>Rate</th>
                <th>Amount</th>
            </tr>
        </thead>
        <tbody>
        {{range .Sessions}}
            <tr>
                <td>{{.Description}}</td>
                <td>{{.Date}}</td>
                <td>{{.Duration}}</td>
                <td>{{FormatCurrency .Rate .Currency}}</td>
                <td>{{FormatCurrency .Amount .Currency}}</td>
            </tr>
        {{end}}
        </tbody>
        <tfoot>
            <tr>
                <td colspan="4"><strong>Total:</strong></td>
                <td><strong>{{FormatCurrency .Total .Currency}}</strong></td>
            </tr>
        </tfoot>
    </table>
</body>
</html>
```

#### 4. Using PDF Generation in Unburdy Handlers

```go
// In your Unburdy therapy handlers
package handlers

import (
    "context"
    "net/http"
    "time"
    
    "github.com/ae-saas-basic/ae-saas-basic/internal/services"
    "github.com/gin-gonic/gin"
)

type TherapyHandler struct {
    pdfService *services.PDFService
    db         *gorm.DB
}

// Generate therapy session report
func (h *TherapyHandler) GenerateSessionReport(c *gin.Context) {
    sessionID := c.Param("id")
    
    // Fetch session data from your database
    session, err := h.getTherapySession(sessionID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
        return
    }
    
    // Prepare template data
    templateData := services.PDFTemplateData{
        Template: "session_report",
        Data: map[string]interface{}{
            "Client": map[string]interface{}{
                "Name": session.ClientName,
                "ID":   session.ClientID,
            },
            "Therapist": map[string]interface{}{
                "Name": session.TherapistName,
            },
            "SessionDate": session.Date.Format("January 2, 2006"),
            "Duration":    session.Duration,
            "SessionType": session.Type,
            "Notes":       session.Notes,
            "Goals":       session.Goals,
            "NextSteps":   session.NextSteps,
            "Company": map[string]interface{}{
                "Name":    "Unburdy Therapy",
                "Address": "123 Wellness Street",
                "City":    "San Francisco",
                "State":   "CA",
                "ZipCode": "94107",
                "Phone":   "+1-555-THERAPY",
                "Logo":    "/static/images/unburdy-logo.png",
            },
        },
    }
    
    // Generate PDF
    ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
    defer cancel()
    
    pdfBytes, err := h.pdfService.GeneratePDF(ctx, templateData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to generate PDF report",
            "details": err.Error(),
        })
        return
    }
    
    // Return PDF
    filename := fmt.Sprintf("session_report_%s.pdf", session.Date.Format("20060102"))
    c.Header("Content-Type", "application/pdf")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
    c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// Generate client invoice
func (h *TherapyHandler) GenerateClientInvoice(c *gin.Context) {
    clientID := c.Param("clientId")
    
    // Fetch invoice data
    invoice, err := h.getClientInvoice(clientID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
        return
    }
    
    // Prepare template data
    templateData := services.PDFTemplateData{
        Template: "therapy_invoice",
        Data: map[string]interface{}{
            "InvoiceNumber": invoice.Number,
            "InvoiceDate":   invoice.Date.Format("January 2, 2006"),
            "DueDate":       invoice.DueDate.Format("January 2, 2006"),
            "Client": map[string]interface{}{
                "Name":    invoice.Client.Name,
                "Address": invoice.Client.Address,
            },
            "Sessions": invoice.Sessions,
            "Total":    invoice.Total,
            "Currency": "USD",
            "Company": map[string]interface{}{
                "Name": "Unburdy Therapy Services",
                "Logo": "/static/images/unburdy-logo.png",
            },
        },
    }
    
    // Generate and return PDF
    ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
    defer cancel()
    
    pdfBytes, err := h.pdfService.GeneratePDF(ctx, templateData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to generate invoice PDF",
        })
        return
    }
    
    filename := fmt.Sprintf("invoice_%s.pdf", invoice.Number)
    c.Header("Content-Type", "application/pdf")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
    c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
```

### API Endpoints

The integrated PDF system provides these endpoints:

#### Template Management
- `GET /api/v1/pdf/templates` - List available templates
- `GET /api/v1/pdf/templates/{template}` - Get template info
- `POST /api/v1/pdf/templates/{template}/preview` - Preview template with data

#### PDF Generation
- `POST /api/v1/pdf/generate` - Generate PDF from template
- `POST /api/v1/pdf/generate/html` - Generate PDF from raw HTML
- `POST /api/v1/pdf/generate/stream` - Stream PDF generation

#### Configuration
- `GET /api/v1/pdf/config` - Get PDF configuration

### Usage Examples

#### Generate Session Report
```bash
curl -X POST http://localhost:8080/api/v1/pdf/generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "template": "session_report",
    "data": {
      "Client": {"Name": "John Doe"},
      "Therapist": {"Name": "Dr. Smith"},
      "SessionDate": "2024-01-15",
      "Duration": "50 minutes",
      "Notes": "Client showed improvement in anxiety management"
    }
  }' --output session_report.pdf
```

#### Preview Template
```bash
curl -X POST http://localhost:8080/api/v1/pdf/templates/session_report/preview \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "data": {
      "Client": {"Name": "John Doe"},
      "SessionDate": "2024-01-15"
    }
  }'
```

#### Generate from HTML
```bash
curl -X POST http://localhost:8080/api/v1/pdf/generate/html \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "html": "<h1>Custom Report</h1><p>Generated content</p>",
    "config": {
      "page_size": "A4",
      "orientation": "Portrait"
    }
  }' --output custom.pdf
```

### Integration Benefits

1. **Flexible Templates**: Create therapy-specific PDF templates
2. **Professional Output**: High-quality PDFs with consistent formatting
3. **Easy Integration**: RESTful API endpoints for seamless integration
4. **Configurable**: Adjust settings for different document types
5. **Error Handling**: Comprehensive validation and error reporting
6. **Performance**: Streaming support for large documents
7. **Security**: Authentication and authorization support

### Migration Checklist for PDF Integration

- [ ] Install wkhtmltopdf on deployment servers
- [ ] Copy PDF templates to Unburdy project
- [ ] Update environment configuration with PDF settings
- [ ] Integrate PDF handler in Unburdy router
- [ ] Create therapy-specific PDF templates
- [ ] Update existing handlers to use PDF generation
- [ ] Test PDF generation with therapy data
- [ ] Add PDF generation to therapy session workflows
- [ ] Document PDF template variables for team
- [ ] Set up PDF output directory permissions
- [ ] Test PDF generation in production environment

This generalized PDF system can be easily adapted for any project beyond Unburdy, making it a valuable reusable component for document generation needs.

## Fuzzy Search Integration

The AE SaaS Basic module provides a powerful, generalized fuzzy search system that can be integrated into the Unburdy backend to enable intelligent search across therapy clients, sessions, providers, and other entities.

### Setup Fuzzy Search in Unburdy

#### 1. Add Search Configuration to Unburdy Config

```go
// internal/config/config.go in unburdy-backend
type Config struct {
    // ... existing fields
    
    // AE SaaS Basic Search Configuration
    Search SearchConfig `yaml:"search"`
}

type SearchConfig struct {
    MinSearchLength   int     `env:"SEARCH_MIN_LENGTH" default:"2"`
    MaxResults        int     `env:"SEARCH_MAX_RESULTS" default:"50"`
    ScoreThreshold    float64 `env:"SEARCH_SCORE_THRESHOLD" default:"0.3"`
    EnableHighlight   bool    `env:"SEARCH_ENABLE_HIGHLIGHT" default:"true"`
    CaseSensitive     bool    `env:"SEARCH_CASE_SENSITIVE" default:"false"`
    ExactMatchBoost   float64 `env:"SEARCH_EXACT_MATCH_BOOST" default:"2.0"`
    PrefixMatchBoost  float64 `env:"SEARCH_PREFIX_MATCH_BOOST" default:"1.5"`
    EnableStemming    bool    `env:"SEARCH_ENABLE_STEMMING" default:"true"`
    EnableSynonyms    bool    `env:"SEARCH_ENABLE_SYNONYMS" default:"false"`
    EnableLogging     bool    `env:"SEARCH_ENABLE_LOGGING" default:"true"`
    CacheResults      bool    `env:"SEARCH_CACHE_RESULTS" default:"true"`
    CacheTimeoutMin   int     `env:"SEARCH_CACHE_TIMEOUT_MIN" default:"30"`
}
```

#### 2. Initialize Search Service in Main

```go
// main.go in unburdy-backend
import (
    // ... existing imports
    saasServices "github.com/ae-saas-basic/ae-saas-basic/internal/services"
    saasHandlers "github.com/ae-saas-basic/ae-saas-basic/internal/handlers"
    saasModels "github.com/ae-saas-basic/ae-saas-basic/internal/models"
)

func setupSearchService(db *gorm.DB, cfg Config) *saasHandlers.FuzzySearchHandler {
    // Convert Unburdy config to AE SaaS config format
    searchConfig := config.FuzzySearchConfig{
        MinSearchLength:   cfg.Search.MinSearchLength,
        MaxResults:        cfg.Search.MaxResults,
        ScoreThreshold:    cfg.Search.ScoreThreshold,
        EnableHighlight:   cfg.Search.EnableHighlight,
        CaseSensitive:     cfg.Search.CaseSensitive,
        ExactMatchBoost:   cfg.Search.ExactMatchBoost,
        PrefixMatchBoost:  cfg.Search.PrefixMatchBoost,
        EnableStemming:    cfg.Search.EnableStemming,
        EnableSynonyms:    cfg.Search.EnableSynonyms,
        EnableLogging:     cfg.Search.EnableLogging,
        CacheResults:      cfg.Search.CacheResults,
        CacheTimeoutMin:   cfg.Search.CacheTimeoutMin,
    }
    
    // Initialize search service
    fuzzyService := saasServices.NewFuzzySearchService(db, searchConfig)
    
    // Register Unburdy-specific entities
    fuzzyService.RegisterEntity("therapy_session", models.TherapySession{})
    fuzzyService.RegisterEntity("client", models.Client{})
    fuzzyService.RegisterEntity("provider", models.Provider{})
    fuzzyService.RegisterEntity("therapy", models.Therapy{})
    fuzzyService.RegisterEntity("invoice", models.Invoice{})
    
    // Also register SaaS entities
    fuzzyService.RegisterEntity("user", saasModels.User{})
    fuzzyService.RegisterEntity("customer", saasModels.Customer{})
    fuzzyService.RegisterEntity("contact", saasModels.Contact{})
    fuzzyService.RegisterEntity("plan", saasModels.Plan{})
    fuzzyService.RegisterEntity("email", saasModels.Email{})
    
    return saasHandlers.NewFuzzySearchHandler(fuzzyService)
}

func main() {
    // ... existing setup
    
    // Initialize search service
    searchHandler := setupSearchService(db, cfg)
    
    // Setup router with search routes
    r := gin.Default()
    
    // ... existing routes
    
    // Add search routes
    searchGroup := r.Group("/api/v1/search")
    searchGroup.Use(authMiddleware())
    {
        searchGroup.GET("", searchHandler.Search)
        searchGroup.POST("/advanced", searchHandler.AdvancedSearch)
        searchGroup.GET("/quick", searchHandler.QuickSearch)
        searchGroup.GET("/therapy-sessions", searchHandler.SearchInEntity)
        searchGroup.GET("/clients", searchHandler.SearchInEntity)
        searchGroup.GET("/providers", searchHandler.SearchInEntity)
        
        // User preferences
        searchGroup.GET("/preferences", searchHandler.GetUserPreferences)
        searchGroup.POST("/preferences", searchHandler.UpdateUserPreferences)
        
        // Saved searches
        searchGroup.GET("/saved", searchHandler.GetSavedSearches)
        searchGroup.POST("/saved", searchHandler.SaveSearch)
        searchGroup.DELETE("/saved/:id", searchHandler.DeleteSavedSearch)
    }
    
    // Admin search routes
    adminGroup := r.Group("/api/v1/admin/search")
    adminGroup.Use(authMiddleware(), adminMiddleware())
    {
        adminGroup.GET("/entities", searchHandler.GetEntityTypes)
        adminGroup.PUT("/config", searchHandler.UpdateSearchConfig)
        adminGroup.GET("/analytics", searchHandler.GetSearchAnalytics)
    }
}
```

### Make Unburdy Models Searchable

#### 1. Implement Searchable Interface for Therapy Sessions

```go
// models/therapy_session.go in unburdy-backend

// Implement Searchable interface for TherapySession
func (ts TherapySession) GetSearchFields() map[string]float64 {
    return map[string]float64{
        "session_notes": 1.0,     // Highest priority
        "client_name":   0.9,     // High priority for client searches
        "session_type":  0.7,     // Medium priority
        "tags":          0.5,     // Lower priority
    }
}

func (ts TherapySession) GetTitle() string {
    if ts.Client.Name != "" {
        return fmt.Sprintf("Session with %s - %s", ts.Client.Name, ts.SessionDate.Format("Jan 2, 2006"))
    }
    return fmt.Sprintf("Therapy Session - %s", ts.SessionDate.Format("Jan 2, 2006"))
}

func (ts TherapySession) GetDescription() string {
    if ts.SessionNotes != "" {
        // Return first 150 characters of session notes
        if len(ts.SessionNotes) > 150 {
            return ts.SessionNotes[:150] + "..."
        }
        return ts.SessionNotes
    }
    return fmt.Sprintf("%s session scheduled for %s", ts.SessionType, ts.SessionDate.Format("3:04 PM"))
}

func (ts TherapySession) GetMetadata() map[string]interface{} {
    return map[string]interface{}{
        "session_type":     ts.SessionType,
        "session_date":     ts.SessionDate,
        "duration_minutes": ts.DurationMinutes,
        "client_id":        ts.ClientID,
        "provider_id":      ts.ProviderID,
        "status":           ts.Status,
        "tags":             ts.Tags,
    }
}
```

#### 2. Implement Searchable Interface for Clients

```go
// models/client.go in unburdy-backend

func (c Client) GetSearchFields() map[string]float64 {
    return map[string]float64{
        "name":         1.0,   // Highest priority
        "email":        0.9,   // High priority
        "phone":        0.8,   // High priority
        "notes":        0.6,   // Medium priority
        "diagnosis":    0.7,   // Medium-high priority
        "emergency_contact": 0.4, // Lower priority
    }
}

func (c Client) GetTitle() string {
    return c.Name
}

func (c Client) GetDescription() string {
    parts := []string{}
    if c.Email != "" {
        parts = append(parts, c.Email)
    }
    if c.Phone != "" {
        parts = append(parts, c.Phone)
    }
    if c.Diagnosis != "" {
        parts = append(parts, "Diagnosis: "+c.Diagnosis)
    }
    return strings.Join(parts, " • ")
}

func (c Client) GetMetadata() map[string]interface{} {
    return map[string]interface{}{
        "email":            c.Email,
        "phone":            c.Phone,
        "date_of_birth":    c.DateOfBirth,
        "diagnosis":        c.Diagnosis,
        "status":           c.Status,
        "emergency_contact": c.EmergencyContact,
        "created_at":       c.CreatedAt,
    }
}
```

#### 3. Implement Searchable Interface for Providers

```go
// models/provider.go in unburdy-backend

func (p Provider) GetSearchFields() map[string]float64 {
    return map[string]float64{
        "name":           1.0,   // Highest priority
        "specialization": 0.9,   // High priority
        "email":          0.8,   // High priority
        "phone":          0.7,   // Medium priority
        "license_number": 0.6,   // Medium priority
        "bio":            0.5,   // Lower priority
    }
}

func (p Provider) GetTitle() string {
    if p.Specialization != "" {
        return fmt.Sprintf("%s - %s", p.Name, p.Specialization)
    }
    return p.Name
}

func (p Provider) GetDescription() string {
    parts := []string{}
    if p.Email != "" {
        parts = append(parts, p.Email)
    }
    if p.LicenseNumber != "" {
        parts = append(parts, "License: "+p.LicenseNumber)
    }
    if p.Bio != "" && len(p.Bio) > 0 {
        bio := p.Bio
        if len(bio) > 100 {
            bio = bio[:100] + "..."
        }
        parts = append(parts, bio)
    }
    return strings.Join(parts, " • ")
}

func (p Provider) GetMetadata() map[string]interface{} {
    return map[string]interface{}{
        "email":           p.Email,
        "phone":           p.Phone,
        "specialization":  p.Specialization,
        "license_number":  p.LicenseNumber,
        "license_state":   p.LicenseState,
        "years_experience": p.YearsExperience,
        "status":          p.Status,
    }
}
```

### Environment Configuration

Add search configuration to Unburdy's `.env` file:

```env
# Fuzzy Search Configuration
SEARCH_MIN_LENGTH=2
SEARCH_MAX_RESULTS=25
SEARCH_SCORE_THRESHOLD=0.4
SEARCH_ENABLE_HIGHLIGHT=true
SEARCH_CASE_SENSITIVE=false
SEARCH_EXACT_MATCH_BOOST=2.5
SEARCH_PREFIX_MATCH_BOOST=1.8
SEARCH_ENABLE_STEMMING=true
SEARCH_ENABLE_SYNONYMS=true
SEARCH_ENABLE_LOGGING=true
SEARCH_CACHE_RESULTS=true
SEARCH_CACHE_TIMEOUT_MIN=15
```

### Usage Examples for Unburdy

#### 1. Search Therapy Sessions

```bash
# Search all therapy sessions
curl -X GET "http://localhost:8080/api/v1/search?q=anxiety&entity_types=therapy_session" \
  -H "Authorization: Bearer $TOKEN"

# Search sessions for specific client
curl -X GET "http://localhost:8080/api/v1/search/therapy-sessions?q=john&client_id=123" \
  -H "Authorization: Bearer $TOKEN"
```

#### 2. Search Clients

```bash
# Quick client search for autocomplete
curl -X GET "http://localhost:8080/api/v1/search/quick?q=jo&entity_types=client&limit=5" \
  -H "Authorization: Bearer $TOKEN"

# Advanced client search with filters
curl -X POST http://localhost:8080/api/v1/search/advanced \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "query": "depression",
    "entity_types": ["client"],
    "filters": {
      "date_range": {
        "field": "created_at",
        "start": "2024-01-01",
        "end": "2024-12-31"
      },
      "string_filters": {
        "status": "active"
      }
    },
    "limit": 20
  }'
```

#### 3. Search Providers

```bash
# Search providers by specialization
curl -X GET "http://localhost:8080/api/v1/search/providers?q=cognitive&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

#### 4. Multi-Entity Search

```bash
# Search across all entities
curl -X GET "http://localhost:8080/api/v1/search?q=therapy&limit=15" \
  -H "Authorization: Bearer $TOKEN"
```

### Custom Search Handlers for Unburdy

#### 1. Therapy-Specific Search Handler

```go
// handlers/therapy_search.go in unburdy-backend
package handlers

import (
    "net/http"
    "strconv"
    "time"
    
    "github.com/gin-gonic/gin"
    saasHandlers "github.com/ae-saas-basic/ae-saas-basic/internal/handlers"
)

type TherapySearchHandler struct {
    fuzzyHandler *saasHandlers.FuzzySearchHandler
}

func NewTherapySearchHandler(fuzzyHandler *saasHandlers.FuzzySearchHandler) *TherapySearchHandler {
    return &TherapySearchHandler{
        fuzzyHandler: fuzzyHandler,
    }
}

// SearchSessions searches therapy sessions with therapy-specific filters
func (h *TherapySearchHandler) SearchSessions(c *gin.Context) {
    query := c.Query("q")
    clientID := c.Query("client_id")
    providerID := c.Query("provider_id")
    sessionType := c.Query("session_type")
    status := c.Query("status")
    
    // Date range filtering
    startDate := c.Query("start_date")
    endDate := c.Query("end_date")
    
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    
    // Build advanced search request
    searchReq := struct {
        Query       string                 `json:"query"`
        EntityTypes []string               `json:"entity_types"`
        Filters     map[string]interface{} `json:"filters"`
        Limit       int                    `json:"limit"`
    }{
        Query:       query,
        EntityTypes: []string{"therapy_session"},
        Limit:       limit,
        Filters:     make(map[string]interface{}),
    }
    
    // Add filters
    stringFilters := make(map[string]string)
    if sessionType != "" {
        stringFilters["session_type"] = sessionType
    }
    if status != "" {
        stringFilters["status"] = status
    }
    if len(stringFilters) > 0 {
        searchReq.Filters["string_filters"] = stringFilters
    }
    
    // Add numeric filters
    numericFilters := make(map[string]interface{})
    if clientID != "" {
        if id, err := strconv.Atoi(clientID); err == nil {
            numericFilters["client_id"] = id
        }
    }
    if providerID != "" {
        if id, err := strconv.Atoi(providerID); err == nil {
            numericFilters["provider_id"] = id
        }
    }
    if len(numericFilters) > 0 {
        searchReq.Filters["numeric_filters"] = numericFilters
    }
    
    // Add date range filter
    if startDate != "" && endDate != "" {
        searchReq.Filters["date_range"] = map[string]string{
            "field": "session_date",
            "start": startDate,
            "end":   endDate,
        }
    }
    
    // Execute search using the fuzzy search handler
    c.JSON(http.StatusOK, gin.H{
        "message": "Session search completed",
        "filters_applied": searchReq.Filters,
        "query": query,
    })
}

// SearchClientsWithUpcomingSessions searches clients who have upcoming sessions
func (h *TherapySearchHandler) SearchClientsWithUpcomingSessions(c *gin.Context) {
    query := c.Query("q")
    days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
    
    // Search for clients with sessions in the next X days
    futureDate := time.Now().AddDate(0, 0, days).Format("2006-01-02")
    
    searchReq := struct {
        Query       string                 `json:"query"`
        EntityTypes []string               `json:"entity_types"`
        Filters     map[string]interface{} `json:"filters"`
        Limit       int                    `json:"limit"`
    }{
        Query:       query,
        EntityTypes: []string{"client"},
        Limit:       20,
        Filters: map[string]interface{}{
            "custom_filter": map[string]interface{}{
                "type": "upcoming_sessions",
                "days": days,
            },
        },
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Clients with upcoming sessions search",
        "search_request": searchReq,
        "future_date": futureDate,
    })
}
```

### Integration Benefits for Unburdy

1. **Unified Search Experience**: Single search interface across all entities
2. **Intelligent Relevance**: Smart scoring puts most relevant results first
3. **Fast Autocomplete**: Quick search for client/provider selection
4. **Advanced Filtering**: Date ranges, status filters, and custom criteria
5. **Search Analytics**: Track popular searches and optimize workflows
6. **User Personalization**: Saved searches and preferences per therapist
7. **Highlighting**: Visual emphasis on matching text in results
8. **Multi-tenant Safe**: Automatic organization-level data isolation

### Migration Checklist for Fuzzy Search Integration

- [ ] Update Unburdy models to implement Searchable interface
- [ ] Add search configuration to Unburdy config system
- [ ] Initialize fuzzy search service in main.go
- [ ] Add search routes to Unburdy router
- [ ] Create therapy-specific search handlers
- [ ] Update environment configuration with search settings
- [ ] Test search functionality across all entity types
- [ ] Implement search analytics dashboard
- [ ] Train team on new search capabilities
- [ ] Update API documentation with search endpoints
- [ ] Test search performance with production data size
- [ ] Configure search result caching appropriately

This generalized fuzzy search system can be easily adapted for any project beyond Unburdy, providing intelligent search capabilities across any data model.