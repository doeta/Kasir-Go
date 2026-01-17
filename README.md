# Kasir-Go API

Kasir-Go is a robust Point of Sale (POS) Backend API built with Go (Golang) and Gin Gonic. It follows a strict role-based access control system (Admin vs Cashier) and provides comprehensive endpoints for managing retail operations.

## Features

- **Strict Role-Based Auth**: Distinct access levels for 'admin' and 'kasir'.
- **Atomic Transactions**: Ensures data integrity for stock reductions during transactions.
- **Product & Stock Management**: Admin-controlled inventory system.
- **Dynamic Payment Methods**: Configurable payment options.
- **Operational Security**: Secured endpoints using JWT with BCrypt hashing.
- **Documentation**: Integrated Swagger UI.

## Tech Stack

- **Language**: Go 1.25+
- **Framework**: Gin Gonic
- **Database**: PostgreSQL
- **ORM**: GORM
- **Auth**: JWT & Bcrypt

## Installation & Run

1. **Clone Repository**

   ```bash
   git clone https://github.com/doeta/Kasir-Go.git
   cd Kasir-Go
   ```

2. **Setup Environment**
   Create `.env`:

   ```env
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=yourpassword
   DB_NAME=kasir_db
   DB_PORT=5432
   JWT_SECRET=supersecretkey
   PORT=8080
   ```

3. **Run Application**
   ```bash
   go mod tidy
   go run main.go
   ```
   _The application will automatically migrate database tables and seed a default admin user on the first run._

## Roles & Permissions

- **Admin**: Full access to manage users, products, and payment methods. Cannot create transactions.
- **Kasir**: Dedicated role for processing transactions (sales). Read-only access to products and payments.

## API Endpoints

### Public

- `GET /` : Health check & API status.
- `POST /api/login` : Login to receive JWT Token.
- `POST /api/register` : Public registration (automatically assigned 'kasir' role).

### Protected (Requires Bearer Token)

#### Transactions

- `POST /api/transactions` : Create new transaction. **(Kasir Only)**
- `GET /api/transactions` : View transaction history.

#### Products

- `GET /api/products` : List all products.
- `POST /api/products` : Add product. **(Admin Only)**
- `PUT /api/products/:id` : Update product. **(Admin Only)**
- `DELETE /api/products/:id` : Delete product. **(Admin Only)**

#### Payments

- `GET /api/payments` : List payment methods.
- `POST /api/payments` : Add payment method. **(Admin Only)**
- `PUT /api/payments/:id` : Update payment method. **(Admin Only)**
- `DELETE /api/payments/:id` : Delete payment method. **(Admin Only)**

#### User Management (Admin Only)

- `GET /api/admin/users` : List all users.
- `POST /api/admin/users` : Create specific user (Admin/Kasir).
- `PUT /api/admin/users/:id` : Update user data.
- `DELETE /api/admin/users/:id` : Delete user.

## API Documentation

Swagger UI is available at:
`http://localhost:8080/swagger/index.html`

## Default Credentials

If the database is empty, the system generates:

- **Username**: `admin`
- **Password**: `admin123`
