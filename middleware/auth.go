package middleware

import (
  "net/http"
  "strings"
  "os"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(getEnv("JWT_SECRET", "-secret-key-change-in-production")) // Claims represents the JWT claims
type Claims struct {
  UserID uint   `json:"user_id"`
  Username string `json:"username"`
  Role     string `json:"role"`
  jwt.RegisteredClaims
}

// getEnv helper function
func getEnv(key, fallback string) string {
  if value := os.Getenv(key); value != "" {
          return value
  }
  return fallback
}

// AuthMiddleware validates JWT token
func AuthMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
          authHeader := c.GetHeader("Authorization")
          if authHeader == "" {
                  c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
                  c.Abort()
                  return
          }

          // Bearer token format
          parts := strings.Split(authHeader, " ")
          if len(parts) != 2 || parts[0] != "Bearer" {
                  c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must start with Bearer"})
                  c.Abort()
                  return
          }

          tokenString := parts[1]
          claims := &Claims{}

          token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
                  return jwtKey, nil
          })

          if err != nil || !token.Valid {
                  c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
                  c.Abort()
                  return
          }

          // Set user info in context
          c.Set("user_id", claims.UserID)
          c.Set("username", claims.Username)
          c.Set("role", claims.Role)
          c.Next()
 }
}

// RoleMiddleware checks if user has required role(s)
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
  return func(c *gin.Context) {
          userRole, exists := c.Get("role")
          if !exists {
                  c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in token"})
                  c.Abort()
                  return
          }
		  
          roleStr := userRole.(string)
          for _, role := range allowedRoles {
                  if roleStr == role {
                          c.Next()
                          return
                  }
        }

        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permission."})
        c.Abort()
        return
}

// GenerateToken creates a JWT token for user
func GenerateToken(userID uint, username, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(expirationTime),
                IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
	}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return token.SignedString(jwtKey)
}