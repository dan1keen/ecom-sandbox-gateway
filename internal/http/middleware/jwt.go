package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token invalid"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		authUser, ok := claims["authUser"].(map[string]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing auth user data"})
			return
		}

		rawID, ok := authUser["user_id"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing user id"})
			return
		}

		userID, err := extractUserID(rawID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user id"})
			return
		}

		c.Set("user_id", userID)

		c.Request.Header.Set("X-User-ID", fmt.Sprintf("%v", userID))

		c.Next()
	}
}

func Unauthenticated(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err == nil && token.Valid {
			c.JSON(http.StatusForbidden, gin.H{"error": "Already authenticated"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func extractUserID(val interface{}) (int, error) {
	switch v := val.(type) {
	case float64:
		return int(v), nil

	case int:
		return v, nil

	case int64:
		return int(v), nil

	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.New("user_id is not a valid number")
		}
		return id, nil

	default:
		return 0, fmt.Errorf("unsupported user_id type: %T", v)
	}
}
