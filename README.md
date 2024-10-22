# E-Procurement System

## Overview
This project is an E-Procurement System designed to streamline the procurement process for organizations. It provides a platform for vendors to register, submit products, and for administrators to manage users and approve vendors.

## Features
- User Registration and Authentication
- Vendor Registration and Approval Process
- Product Management
- Role-based Access Control (Admin, Vendor, User)
- RESTful API

## Tech Stack
- Go (Golang)
- Chi Router for HTTP routing
- PostgreSQL for database
- JWT for authentication
- Docker and Docker Compose for containerization

## Getting Started
### Running without Docker

1. Ensure you have Go and PostgreSQL installed on your system.

2. Clone the repository and navigate to the project directory.

3. Copy and update the `config.yaml`.

4. Install dependencies:
   ```
   go mod tidy
   ```

5. Update the `config.yaml` file with your local PostgreSQL connection details.

6. Run the application:
   ```
   go run cmd/main.go
   ```

## Default Admin Credentials
- Email: admin@example.com
- Password: password

## API Endpoints

### Public Endpoints
- POST `/api/v1/login`: User login
- POST `/api/v1/register-vendor`: Vendor registration
- GET `/api/v1/products`: Get all products
- GET `/api/v1/products/{id}`: Get product by ID

### Protected Endpoints (Require Authentication)
- GET `/api/v1/users/{id}`: Get user details
- PUT `/api/v1/users/{id}`: Update user details

### Admin-only Endpoints
- PUT `/api/v1/users/{id}/approve`: Approve vendor
- PUT `/api/v1/users/{id}/reject`: Reject vendor
- DELETE `/api/v1/users/{id}`: Delete user

### Vendor-only Endpoints
- POST `/api/v1/products`: Create a new product
- PUT `/api/v1/products/{id}`: Update a product
- DELETE `/api/v1/products/{id}`: Delete a product
- GET `/api/v1/my-products`: Get all products created by the vendor

## Docker Configuration

This project uses Docker and Docker Compose for easy deployment and development. Key files:

- `Dockerfile`: Defines the container for the Go application.
- `docker-compose.yaml`: Orchestrates the application and database services.

To modify the Docker setup, edit these files and rebuild the Docker images:
