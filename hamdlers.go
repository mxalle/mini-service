package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// ========== AUTH ==========

func Signup(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "bad request"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := User{Name: req.Name, Email: req.Email, PasswordHash: string(hash)}
	if err := DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"id": user.ID, "email": user.Email})
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.BindJSON(&req)

	var user User
	if err := DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	c.JSON(200, gin.H{"token": signed})
}

// ========== POSTS ==========

func CreatePost(c *gin.Context) {
	var req struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	c.BindJSON(&req)
	uid := c.GetUint("userID")
	post := Post{Title: req.Title, Body: req.Body, UserID: uid}
	DB.Create(&post)
	c.JSON(201, post)
}

func ListPosts(c *gin.Context) {
	var posts []Post
	DB.Preload("User").Find(&posts)
	c.JSON(200, posts)
}

func DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var post Post
	if err := DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	uid := c.GetUint("userID")
	if post.UserID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	DB.Delete(&post)
	c.JSON(200, gin.H{"deleted": true})
}
