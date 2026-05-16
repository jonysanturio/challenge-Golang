 # Product & Category Management API

A RESTful API with real-time WebSocket updates for managing products and categories, built with Go (Golang) and PostgreSQL.

## Features
- **RESTful API**: Complete CRUD operations for products and categories
- **Real-time Updates**: WebSocket integration for live product/category updates
- **Authentication & Authorization**: JWT-based auth with role-based access control (admin/client)
- **Advanced Filtering**: Pagination, sorting, and search capabilities
- **Data History**: Track product price and stock changes over time
- **Data Validation**: Input validation and error handling
- **Database Seeding**: Sample data for development and testing

## Technology Stack

- **Language**: Go 1.19+
- **Framework**: Gin Web Framework
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt/jwt)
- **WebSockets**: Gorilla WebSocket
- **Environment**: Godotenv for configuration management

## API Endpoints

### Products

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/api/products` | List products (paginated, filterable, sortable) | Public |
| GET | `/api/products/{id}` | Get product by ID | Public |
| POST | `/api/products` | Create new product | Admin only |
| PUT | `/api/products/{id}` | Update product | Admin only |
| DELETE | `/api/products/{id}` | Delete product | Admin only |
| GET | `/api/products/{id}/history` | Get product price/stock history | Public |
### Categories

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/api/categories` | List categories (paginated) | Public |
| GET | `/api/categories/{id}` | Get category by ID | Public |
| POST | `/api/categories` | Create new category | Admin only |
| PUT | `/api/categories/{id}` | Update category | Admin only |
| DELETE | `/api/categories/{id}` | Delete category | Admin only |
### Search

| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/api/search?q={query}&type={product|category}` | Search products or categories | Public |

### WebSocket

| Endpoint | Description |
|----------|-------------|
| `GET /ws` | WebSocket connection for real-time updates |

## WebSocket Events

The WebSocket server broadcasts the following events:

- `product.created` - When a new product is created
- `product.updated` - When a product is updated
- `product.deleted` - When a product is deleted
- `category.created` - When a new category is created
- `category.updated` - When a category is updated
- `category.deleted` - When a category is deleted

Events are sent as JSON messages with the following structure:
```json
{
  "event": "product.created",
  "data": { /* product or category object */ }
}
```

 ## Database Schema

 ### Tables

 1. **products**
    - id (PK)
    - name
    - description
    - price
    - stock
    - created_at
    - updated_at
2. **categories**
    - id (PK)
    - name
    - description
    - created_at
    - updated_at

3. **product_categories** (join table)
    - product_id (FK)
    - category_id (FK)

4. **product_history**
    - id (PK)
    - product_id (FK)
    - price
    - stock
    - changed_at
### Relationships

- Products ⇄ Categories: Many-to-many (through product_categories)
- Products → Product History: One-to-many

## Installation & Setup

### Prerequisites

- Go 1.19 or higher
- PostgreSQL 12 or higher
- Git

### Steps

1. **Clone the repository**
```bash
git clone <repository-url>
cd go-product-api
```

2. **Install dependencies**
```bash
go mod download
```

3. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. **Create database**
```bash
# Using psql
createdb product_api
# Or with your preferred PostgreSQL admin tool
```

5. **Run migrations and seeders**
   ```bash
   go run cmd/main.go
   # The application will automatically migrate and seed the database on startup
   ```

6. **Start the server**
```bash
go run cmd/main.go
# Server will start on http://localhost:8080
```

## Usage Examples

### Get all products
```bash
curl http://localhost:8080/api/products
```

### Get product with ID 1
```bash
curl http://localhost:8080/api/products/1
```

### Create a new product (requires admin token)
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
"name": "New Product",
"description": "Product description",
"price": 29.99,
"stock": 100,
"category_ids": [1, 2]
}'
```

### Search for products
```bash
curl "http://localhost:8080/api/search?q=phone&type=product"
```

### Connect to WebSocket
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = function(event) {
console.log('Received:', event.data);
// Handle real-time updates here
};

ws.onopen = function() {
console.log('WebSocket connected');
};
```

## Authentication
### Obtaining a JWT Token

For demonstration purposes, you can create a simple login endpoint (not included in base implementation) or generate a token manually:

```go
token, err := middleware.GenerateToken(userID, username, role)
```

Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### Roles

- **admin**: Full access to all endpoints
- **client**: Read-only access to public endpoints (can be extended)

## Project Structure

```
go-product-api/
 ├── cmd/               # Application entry points
 ├── main.go           # Main application
 ├── config/           # Configuration (database, env)
 ├── handlers/         # HTTP request handlers
 │   ├── product.go    # Product handlers
 │   ├── category.go   # Category handlers
 │   └── search.go     # Search handlers
 ├── middleware/         # Custom middleware
 │   └── auth.go       # JWT authentication & authorization
 ├── models/           # Data models and database schemas
 │   └── product.go    # Product, Category ProductHistory models
 ├── seeders/          # Database seeders
 │   └── seeders.go    # Sample data generation
 ├── websocket/        # WebSocket implementation
 │   └── hub.go        # WebSocket hub and client management
 ├── docs/             # Documentation (diagrams, etc.)
 ├── go.mod            # Go dependencies
 ├── .env.example      # Environment variables template
 └── README.md         # This file
```

## Design Decisions

### 1. Modular Architecture
Separated concerns into distinct packages:
- **handlers**: HTTP request/response logic
- **models**: Data structures and database mappings
- **middleware**: Cross-cutting concerns (auth, logging)
- **websocket**: Real-time communication layer
- **seeders**: Initial data population

### 2. Database Relationships
- Used GORM for ORM capabilities
- Many-to-many relationship between products and categories via join table
- Product history table for audit trail of price/stock changes

### 3. WebSocket Implementation
- Hub pattern for managing client connections
- Broadcast mechanism for sending updates to all connected clients
- Gorilla WebSocket library for reliable WebSocket handling

### 4. Security
- JWT-based authentication with HS256 signing
- Role-based access control middleware
- Input validation using Gin's binding capabilities
- Environment-based configuration for secrets

### 5. Performance Considerations
- Pagination for large dataset handling
- Database indexing on frequently queried fields
- Efficient WebSocket broadcasting with goroutine separation
 - Prepared statements through GORM

 ## Deployment

 ### Docker (Example)
 ```dockerfile
 FROM golang:1.19-alpine AS builder

 WORKDIR /app

 COPY go.mod go.sum ./
 RUN go mod download

 COPY . .

 RUN go build -o main cmd/main.go

 FROM alpine:latest
 RUN apk --no-cache add ca-certificates

 WORKDIR /root/
 COPY --from=builder /app/main .
 COPY --from=builder /app/.env.example ./.env

 EXPOSE 8080

 CMD ["./main"]
 ```

### Kubernetes
Refer to the `docs/` directory for sample Kubernetes manifests.

## Testing

### Unit Tests
```bash
go test ./...
```

### Integration Tests
Use tools like Postman or curl to test endpoints:
- Test authentication flows
- Test CRUD operations
- Test WebSocket connections
- Test error handling and validation

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Gin Web Framework team
- GORM ORM library
- Gorilla WebSocket contributors
- JWT library maintainers