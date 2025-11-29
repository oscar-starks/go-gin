# Go Gin REST API with GORM and PostgreSQL

A modern REST API built with Go, Gin framework, GORM ORM, and PostgreSQL database featuring JWT authentication, password hashing, and clean architecture.

## 🚀 Features

- **RESTful API** with Gin framework
- **PostgreSQL** database with GORM ORM
- **JWT Authentication** with middleware protection
- **Password Hashing** using bcrypt
- **Clean Architecture** with organized packages
- **Environment Configuration** with .env support
- **Auto Database Migration** with GORM
- **Soft Deletes** implementation
- **Request Validation** with struct tags

## 📁 Project Structure

```
gin-project/
├── config/                 # Database configuration
│   └── database.go         # DB connection and setup
├── handlers/               # HTTP handlers
│   └── auth.go            # Authentication handlers
├── middleware/             # HTTP middleware
│   └── auth.go            # JWT authentication middleware
├── models/                 # Database models and request/response types
│   ├── auth.go            # Auth request/response structs
│   └── users.go           # User model
├── utils/                  # Utility functions
│   └── auth.go            # JWT and password utilities
├── .env.example           # Environment variables template
├── .gitignore             # Git ignore rules
├── go.mod                 # Go module definition
└── main.go                # Application entry point
```

## 🛠️ Tech Stack

- **Go 1.25+** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library
- **PostgreSQL** - Database
- **JWT** - Authentication tokens
- **Bcrypt** - Password hashing
- **Godotenv** - Environment configuration

## 📋 Prerequisites

- Go 1.25 or higher
- PostgreSQL database
- Git

## 🚀 Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/oscar-starks/go-gin.git
   cd go-gin
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up environment variables**

   ```bash
   cp .env.example .env
   ```

   Edit `.env` with your database credentials:

   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=your_database_name
   DB_SSLMODE=disable
   DB_TIMEZONE=UTC

   PORT=8080
   JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
   ```

4. **Create PostgreSQL database**

   ```sql
   CREATE DATABASE your_database_name;
   ```

5. **Run the application**

   ```bash
   go run main.go
   ```

   Or build and run:

   ```bash
   go build -o main .
   ./main
   ```

The server will start on `http://localhost:8080`

## 📚 API Endpoints

### Public Endpoints

- `GET /` - Welcome message
- `GET /health` - Health check
- `POST /auth/register` - User registration
- `POST /auth/login` - User login

### Protected Endpoints (Requires JWT Token)

- `GET /api/profile` - Get user profile

## 🔐 Authentication

### Register User

```bash
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "age": 30
}
```

### Login User

```bash
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

### Access Protected Routes

Include JWT token in Authorization header:

```bash
GET /api/profile
Authorization: Bearer <your-jwt-token>
```

## 🏗️ Database Schema

### Users Table

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  age INTEGER
);
```

## 🧪 Testing

Build the project to check for errors:

```bash
go build -o main .
```

## 🔧 Development

### Adding New Routes

1. Create handler functions in `handlers/` package
2. Add route definitions in `main.go`
3. Use middleware for protected routes

### Database Models

1. Define models in `models/` package
2. Add to auto-migration in `main.go`
3. Use GORM tags for database constraints

### Environment Variables

All configuration is handled through environment variables. See `.env.example` for available options.

## 📝 Code Features

- **Clean Package Structure** - Organized by responsibility
- **Error Handling** - Consistent error responses
- **Input Validation** - Request validation with struct tags
- **Security** - Password hashing and JWT tokens
- **Database** - Connection pooling and soft deletes
- **Middleware** - JWT authentication protection

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 🙏 Acknowledgments

- [Gin Framework](https://gin-gonic.com/) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Password hashing

---

Built with ❤️ using Go and Gin
