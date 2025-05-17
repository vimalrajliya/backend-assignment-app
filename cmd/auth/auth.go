package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vimalrajliya/backend-assignment-app/database"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type AuthClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(email string, userId uint) (string, error) {
	claims := AuthClaims{
		UserID: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

/*
* authentication of session
 1. Check token is present
 2. Check decoded token has proper format
 3. Check if token is not blacklisted
 4. Check token is not expired
 5. Check User exists
*/
func AuthenticateToken() gin.HandlerFunc {

	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing token header",
			})
			c.Abort()
			return
		}

		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString = parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("invalid token")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(*AuthClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "expired token"})
			c.Abort()
			return
		}
		// check if token is blacklisted
		blacklistKey := "blacklist:token:" + tokenString
		_, err = database.Client.Get(c, blacklistKey).Result()
		if err == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()

	}
}

func BlacklistToken(c *gin.Context) error {
	tokenString := c.GetHeader("Authorization")
	parts := strings.Split(tokenString, " ")
	tokenString = parts[1]
	print(tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}

	expiry := time.Until(claims.ExpiresAt.Time)
	if expiry <= 0 {
		return fmt.Errorf("token already expired")
	}
	println("expiry", c)

	blacklistKey := "blacklist:token:" + tokenString
	if err := database.Client.Set(c, blacklistKey, "1", expiry).Err(); err != nil {
		return err
	}

	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
