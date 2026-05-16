# Deployment Guide

## Overview
This document provides instructions for deploying the Product & Category Management API in various environments.

 ## Prerequisites
 - Go 1.19+ installed
 - PostgreSQL 12+ accessible
 - Git for cloning the repository

 ## Local Development

 ### Step 1: Clone Repository
 ```bash
 git clone <repository-url>
 cd go-product-api
 ```

 ### Step 2: Install Dependencies
 ```bash
 go mod download
 ```

 ### Step 3: Configure Environment
 ```bash
 cp .env.example .env
 # Edit .env file with your local database credentials
 ```

 ### Step 4: Initialize Database
 ```bash
 # Create database (if not exists)
 createdb product_api

 # Or using psql:
 # CREATE DATABASE product_api;
 ```

 ### Step 5: Run Application
 ```bash
 go run cmd/main.go
 ```

 The application will:
 1. Automatically run database migrations
 2. Seed initial data
 3. Start the server on port 8080 (or PORT from .env)

 ## Docker Deployment

 ### Building the Docker Image
 ```bash
 docker build -t product-api:latest .
 ```
   
 ### Running with Docker Compose
 Create a `docker-compose.yml` file:

 ```yaml
 version: '3.8'

 services:
   postgres:
     image: postgres:15-alpine
     container_name: product-api-postgres
     environment:
       POSTGRES_DB: product_api
       POSTGRES_USER: postgres
       POSTGRES_PASSWORD: postgres
     volumes:
       - postgres_data:/var/lib/postgresql/data
     ports:
       - "5432:5432"
     healthcheck:
       test: ["CMD", "pg_isready", "-U", "postgres"]
       interval: 10s
       timeout: 5s
       retries: 5

   api:
     build: .
     container_name: product-api
     ports:
       - "8080:8080"
     environment:
       - DB_HOST=postgres
       - DB_PORT=5432
       - DB_USER=postgres
       - DB_PASSWORD=postgres
       - DB_NAME=product_api
       - JWT_SECRET=your-super-secret-key-change-in-production
     depends_on:
       postgres:
         condition: service_healthy
     volumes:
       - ./:/app
       - ~/.ssh:/root/.ssh:ro

 volumes:
   postgres_data:
 ```

 Then run:
 ```bash
 docker-compose up -d
 ```

 ### Running Docker Container Directly
 ```bash
 docker run -d \
   --name product-api \
   -p 8080:8080 \
   -e DB_HOST=host.docker.internal \
   -e DB_PORT=5432 \
   -e DB_USER=postgres \
   -e DB_PASSWORD=postgres \
   -e DB_NAME=product_api \
   -e JWT_SECRET=your-super-secret-key-change-in-production \
  product-api:latest
```

## Kubernetes Deployment

### Prerequisites
- Kubernetes cluster (1.20+)
- kubectl configured
- PostgreSQL available (can use Helm chart or external service)

### Deployment Manifests
Create the following files:

#### 1. Namespace (optional)
 ```yaml
 apiVersion: v1
 kind: Namespace
 metadata:
   name: product-api
 ```

 #### 2. PostgreSQL Deployment (if managing internally)
 ```yaml
 apiVersion: apps/v1
 kind: Deployment
 metadata:
   name: postgres
   labels:
     app: postgres
 spec:
   replicas: 1
   selector:
     matchLabels:
       app: postgres
   template:
     metadata:
       labels:
         app: postgres
     spec:
       containers:
       - name: postgres
         image: postgres:15-alpine
         env:
         - name: POSTGRES_DB
           value: product_api
         - name: POSTGRES_USER
           value: postgres
         - name: POSTGRES_PASSWORD
           valueFrom:
             secretKeyRef:
               name: postgres-secret
               key: postgres-password
         ports:
         - containerPort: 5432
         volumeMounts:
         - name: postgres-storage
           mountPath: /var/lib/postgresql/data
       volumes:
       - name: postgres-storage
         persistentVolumeClaim:
           claimName: postgres-pvc
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: postgres
 spec:
   selector:
     app: postgres
   ports:
     - protocol: TCP
       port: 5432
       targetPort: 5432
   type: ClusterIP
 ```

 #### 3. API Deployment
 ```yaml
 apiVersion: apps/v1
 kind: Deployment
 metadata:
   name: product-api
   labels:
     app: product-api
 spec:
   replicas: 2
   selector:
     matchLabels:
       app: product-api
   template:
     metadata:
       labels:
         app: product-api
     spec:
       containers:
       - name: api
         image: your-registry/product-api:latest
         ports:
         - containerPort: 8080
         env:
         - name: DB_HOST
  value: postgres  # or your external postgres host
         - name: DB_PORT
           value: "5432"
         - name: DB_USER
           value: postgres
         - name: DB_PASSWORD
           valueFrom:
             secretKeyRef:
               name: product-api-secrets
               key: db-password
         - name: DB_NAME
           value: product_api
         - name: JWT_SECRET
           valueFrom:
             secretKeyRef:
               name: product-api-secrets
               key: jwt-secret
         - name: PORT
           value: "8080"
         readinessProbe:
           httpGet:
             path: /api/products
             port: 8080
           initialDelaySeconds: 10
           periodSeconds: 10
         livenessProbe:
           httpGet:
             path: /api/products
             port: 8080
           initialDelaySeconds: 30
           periodSeconds: 30
       imagePullPolicy: IfNotPresent
 ---
 apiVersion: v1
 kind: Service
 metadata:
   name: product-api-service
 spec:
   selector:
     app: product-api
    ports:
      - protocol: TCP
        port: 80
        targetPort: 8080
    type: LoadBalancer  # or NodePort/ClusterIP based on your setup
```

#### 4. Secrets
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: product-api-secrets
type: Opaque
data:
  db-password: cG9zdGdyZXM=  # base64 encoded "postgres"
  jwt-secret: eW91ci1zdXBlci1zZWNyZXQta2V5LWNoYW5nZS1pbi1wcm9kdWN0aW9u  # base64 encoded
```

### Deploying
```bash
# Apply namespace (if using)
kubectl apply -f namespace.yaml

# Apply secrets
kubectl apply -f secrets.yaml

# Apply PostgreSQL (if managing internally)
kubectl apply -f postgres-deployment.yaml

# Apply API
kubectl apply -f api-deployment.yaml
```

## Configuration Options

### Environment Variables
| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | PostgreSQL host | localhost |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | PostgreSQL username | postgres |
| DB_PASSWORD | PostgreSQL password | postgres |
| DB_NAME | PostgreSQL database name | product_api |
| DB_SSLMODE | PostgreSQL SSL mode | disable |
| DATABASE_URL | Full PostgreSQL connection string (overrides individual DB vars) | "" |
| JWT_SECRET | Secret key for JWT signing | your-secret-key-change-in-production |
| PORT | Server port | 8080 |

### Database Connection
The application supports two ways to configure database connection:
1. Individual parameters (DB_HOST, DB_PORT, etc.)
2. Full DATABASE_URL connection string

If DATABASE_URL is set, it takes precedence over individual parameters.

Example DATABASE_URL:
```
postgresql://username:password@host:port/database?sslmode=disable
```

## Health Checks
The application provides implicit health checks through:
- Root endpoint: GET / (returns 404 but server is responding)
- Products endpoint: GET /api/products (returns 200 when database is connected)

For explicit health checks, you can add a middleware or use the existing endpoints.

## Logging
The application uses Go's standard logger:
- Startup messages
- Database connection status
- WebSocket connection events
- Error conditions

For production, consider redirecting output to a logging system or using a structured logging library.

## Scaling Considerations
### Horizontal Scaling
- The API is stateless (except for WebSocket connections)
- Multiple instances can run behind a load balancer
- WebSocket connections are sticky to instances (consider Redis adapter for pub/sub if needed)

### Database Considerations
- Ensure proper indexing on frequently queried columns
- Monitor connection pool size (GORM defaults to 10 connections)
- Consider read replicas for heavy read workloads

### Caching
For high-traffic scenarios, consider adding:
- Redis for caching frequent queries
- CDN for static assets (if any)
- API gateway for rate limiting and SSL termination

## Monitoring
Key metrics to monitor:
- Database connection usage
- API response times and error rates
- WebSocket connection count
- Memory and CPU usage
- Request throughput

Consider integrating with:
- Prometheus + Grafana
- ELK stack for log aggregation
- Application Performance Monitoring (APM) tools

## Backup and Recovery
### Database Backup
```bash
# Using pg_dump
pg_dump -U postgres -h localhost product_api > backup_$(date +%Y%m%d).sql
  
  # Restore
  psql -U postgres -h localhost product_api < backup_$(date +%Y%m%d).sql
  ```
  
### Application Recovery
Since the application is stateless, recovery involves:
1. Ensuring database is available
2. Restarting application containers/instances
3. Verifying health checks pass
  
## Troubleshooting

### Common Issues

#### Database Connection Failures
- Verify PostgreSQL service is running
- Check network connectivity between app and database
- Validate credentials in .env or environment variables
- Confirm database exists and user has proper permissions
  
#### "Address already in use" Error
- Another process is using the configured port
- Change PORT environment variable or stop conflicting process

#### WebSocket Connection Issues
- Check if proxy/load balancer supports WebSocket upgrades
- Verify CORS settings in WebSocket upgrader
- Check firewall/WebSocket timeout settings

#### Performance Problems
- Enable slow query logging in PostgreSQL
- Check for missing indexes on queried columns
- Monitor garbage collection pauses in Go
- Consider connection pool tuning

### Log Locations
- Standard output/stderr (captured by container orchestrator or systemd)
- Can be redirected to file: `go run cmd/main.go > app.log 2>&1`

## Version Updates
When upgrading:
1. Backup database
2. Review changelog for breaking changes
3. Update Go dependencies: `go get -u ./...`
4. Run migrations (handled automatically on startup)
5. Test in staging environment before production

## Support
For issues and questions:
- Check the README.md for basic usage
- Review API examples in docs/api_examples.md
- Examine source code for implementation details