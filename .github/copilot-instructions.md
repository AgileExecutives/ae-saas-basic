<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->
- [x] Verify that the copilot-instructions.md file in the .github directory is created. ✅ COMPLETED

- [x] Clarify Project Requirements ✅ COMPLETED - Go SaaS module with Gin, GORM, PostgreSQL, JWT

- [x] Scaffold the Project ✅ COMPLETED
	Created complete Go module structure with models, handlers, middleware, database migrations, router setup, and configuration management.

- [x] Customize the Project ✅ COMPLETED
	Extracted and adapted essential SaaS functionality from unburdy-backend including User, Customer, Plan, Email, Organization, UserSettings models and corresponding handlers for auth, customer, contact, email, health, plans, and user-settings endpoints.

- [x] Install Required Extensions ✅ COMPLETED - No VS Code extensions required for Go module

- [x] Compile the Project ✅ COMPLETED
	Project builds successfully with `go build .` - all dependencies resolved and compilation errors fixed.

- [x] Create and Run Task ✅ COMPLETED - Go project uses standard commands
	Created Makefile with build, run, test targets. Use `make run` or `go run main.go` to start server.

- [x] Launch the Project ✅ COMPLETED
	Server starts on localhost:8080 with health check at /api/v1/health and full REST API for SaaS functionality.

- [x] Ensure Documentation is Complete ✅ COMPLETED
	README.md, INTEGRATION_EXAMPLE.md, Makefile, .env.example, and basic test suite completed.