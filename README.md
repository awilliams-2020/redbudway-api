# Redbud Way API

Redbud Way API is a backend service written in Go, designed to power a platform that connects customers with tradespeople for various services. The API provides endpoints for user authentication, account management, service quotes, fixed-price bookings, payments, and more. It is built using the go-swagger toolkit and integrates with Stripe for payment processing.

## Features

- **User Authentication & Account Management**
  - Supports both customer and tradesperson accounts
  - JWT-based authentication and session management
  - Password reset and email verification flows

- **Service Discovery & Booking**
  - Customers can browse and book services (quotes or fixed prices) offered by tradespeople
  - Tradespeople can manage their service offerings, schedules, and branding

- **Payments & Invoicing**
  - Stripe integration for secure payments, subscriptions, and invoicing
  - Tradespeople onboarding and account management via Stripe Connect

- **Reviews & Ratings**
  - Customers can review services and tradespeople
  - Ratings and repeat customer tracking

- **Admin & Security**
  - Admin endpoints for platform management
  - Secure API endpoints with Bearer token authentication

## API Documentation

The API is fully documented using Swagger (OpenAPI 2.0). You can find the API specification in [`swagger.yaml`](./swagger.yaml).

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Docker (optional, for containerized deployment)
- MySQL database
- Stripe account (for payment integration)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/redbudway-api.git
   cd redbudway-api
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Configure environment variables:**
   - Database connection settings
   - Stripe API keys
   - Email service credentials
   - (See `.env.example` if available, or check the `internal/` and `database/` packages for required variables)

4. **Run the server:**
   ```bash
   go run ./cmd/redbud-way-api-server/main.go
   ```

   Or using Docker:
   ```bash
   docker build -t redbudway-api .
   docker run -p 8080:8080 --env-file .env redbudway-api
   ```

### Usage

- The API will be available at `http://localhost:8080/v1` (or as configured).
- Use the Swagger UI or tools like Postman to explore and test endpoints.

## Project Structure

- `cmd/redbud-way-api-server/` - Main entry point for the API server
- `handlers/` - Business logic for API endpoints (accounts, quotes, bookings, etc.)
- `models/` - Data models and API payload definitions
- `database/` - Database connection and queries
- `email/` - Email sending logic
- `security/` - Security and authentication helpers
- `stripe/` - Stripe payment integration
- `restapi/` - Swagger-generated server code and API operations

## Contributing

Contributions are welcome! Please open issues or submit pull requests for improvements or bug fixes.

## License

[MIT](LICENSE) (or specify your license) 