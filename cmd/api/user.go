package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vimalrajliya/backend-assignment-app/cmd/auth"
	"github.com/vimalrajliya/backend-assignment-app/database"
	"github.com/vimalrajliya/backend-assignment-app/models"
	"gorm.io/gorm"
)

func PostUser(c *gin.Context) {
	var newUser models.User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// add created at and updated at in user details
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	newUser.Password = hashedPassword
	result := database.DB.Db.Create(&newUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) ||
			strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.email") {
			c.JSON(http.StatusConflict, gin.H{"error": "email is already registered"})
			return
		}
		return
	}
	c.IndentedJSON(http.StatusCreated, "User has been created successfully")
}

func SignInUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	result := database.DB.Db.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect Email"})
		return
	}

	if !auth.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenString, err := auth.GenerateToken(user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func GetUserDetails(c *gin.Context) {
	var userOutput struct {
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	var user models.User
	userId := c.GetUint("user_id")
	userEmail := c.GetString("email")
	result := database.DB.Db.Where("id = ?", userId).Where("email = ?", userEmail).Find(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	userOutput.Name = user.Name
	userOutput.Email = user.Email
	userOutput.CreatedAt = user.CreatedAt
	userOutput.UpdatedAt = user.UpdatedAt
	c.IndentedJSON(http.StatusOK, userOutput)
}

func RefreshToken(c *gin.Context) {
	var user models.User
	userId := c.GetUint("user_id")
	userEmail := c.GetString("email")
	result := database.DB.Db.Where("id = ?", userId).Where("email = ?", userEmail).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	tokenString, err := auth.GenerateToken(user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func LogOutUser(c *gin.Context) {
	auth.BlacklistToken(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged out successfully",
	})
}
