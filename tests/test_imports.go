package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = gin.New()
	_ = jwt.New(jwt.SigningMethodES256)
	_ = assert.Equal                //UNTUK TESTING
	_ = bcrypt.GenerateFromPassword //HASHING
	_ = postgres.Open
	_ = &gorm.Config{}
}
