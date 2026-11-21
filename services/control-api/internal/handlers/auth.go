package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/brlia/control-api/internal/config"
	"github.com/brlia/control-api/internal/models"
)

// GetTokenHandler handles JWT token generation
func GetTokenHandler(cfg config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// In production, validate username and password against database
		// For now, use a simple check (replace with proper authentication)
		if req.Username != "admin" || req.Password != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		// Create JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":   req.Username,
			"iat":   time.Now().Unix(),
			"exp":   time.Now().Add(cfg.JWTExpiry).Unix(),
			"iss":   cfg.Issuer,
			"roles": []string{"admin"},
		})

		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
				"details": err.Error(),
			})
			return
		}

		// Create refresh token
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": req.Username,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(cfg.RefreshExpiry).Unix(),
			"iss": cfg.Issuer,
		})

		refreshTokenString, err := refreshToken.SignedString([]byte(cfg.RefreshSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate refresh token",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.LoginResponse{
			AccessToken:  tokenString,
			RefreshToken: refreshTokenString,
			TokenType:    "Bearer",
			ExpiresIn:    int64(cfg.JWTExpiry.Seconds()),
			UserID:       0, // Replace with actual user ID
			Username:     req.Username,
			DisplayName:  req.Username,
			Email:        "", // Replace with actual email
			MFARequired:  false,
		})
	}
}

// RefreshTokenHandler handles JWT token refresh
func RefreshTokenHandler(cfg config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.TokenRefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		// Parse and validate refresh token
		token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.RefreshSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid refresh token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			return
		}

		username, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token subject",
			})
			return
		}

		// Create new JWT token
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":   username,
			"iat":   time.Now().Unix(),
			"exp":   time.Now().Add(cfg.JWTExpiry).Unix(),
			"iss":   cfg.Issuer,
			"roles": []string{"admin"},
		})

		newTokenString, err := newToken.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
				"details": err.Error(),
			})
			return
		}

		// Create new refresh token
		newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": username,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(cfg.RefreshExpiry).Unix(),
			"iss": cfg.Issuer,
		})

		newRefreshTokenString, err := newRefreshToken.SignedString([]byte(cfg.RefreshSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate refresh token",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.TokenRefreshResponse{
			AccessToken:  newTokenString,
			RefreshToken: newRefreshTokenString,
			TokenType:    "Bearer",
			ExpiresIn:    int64(cfg.JWTExpiry.Seconds()),
		})
	}
}
