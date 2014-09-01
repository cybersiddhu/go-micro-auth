package service

import (
	"database/sql"
	"fmt"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/cybersiddhu/go-micro-auth/api"
	"gopkg.in/dgrijalva/jwt-go.v2"
	"gopkg.in/gin-gonic/gin.v0"
	"gopkg.in/jmoiron/sqlx.v0"
)

const (
	userCheckStmt = `
	SELECT users.* FROM users where email = ?
	`
	userInsertStmt = `
	INSERT INTO users(email, password) VALUES(?, ?)
	`
)

type AuthResource struct {
	Dbh    *sqlx.DB
	PrvKey []byte
}

func (ar *AuthResource) CreateSession(c *gin.Context) {
	var ju api.UserJSON
	// validate input by binding to an expected json structure
	if !c.Bind(&ju) {
		c.JSON(400, gin.H{"message": "malformed json"})
	}

	// check if the email matches
	u := api.UserJSON{}
	err := ar.Dbh.Get(&u, userCheckStmt, ju.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(400, gin.H{"message": fmt.Sprintf("email %s not found", ju.Email)})
		} else {
			c.JSON(401, gin.H{"message": err})
		}
	}

	// check if password matches
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(ju.Password))
	if err != nil {
		c.JSON(401, gin.H{"message": "password do not match"})
	}

	// generate a json web token
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["user_id"] = u.Id
	token.Claims["email"] = u.Email
	tokenString, err := token.SignedString(ar.PrvKey)
	if err != nil {
		c.JSON(400, gin.H{"message": err})
	}
	// successful response
	c.JSON(201, gin.H{"token": tokenString})
}

func (ar *AuthResource) CreateUser(c *gin.Context) {
	var ju api.UserJSON
	// validate input by binding to an expected json structure
	if !c.Bind(&ju) {
		c.JSON(400, gin.H{"message": "malformed json"})
	}
	epass, err := bcrypt.GenerateFromPassword([]byte(ju.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(400, gin.H{"message": err})
	}
	_, err = ar.Dbh.Exec(userInsertStmt, ju.Email, epass)
	if err != nil {
		c.JSON(400, gin.H{"message": err})
	}
	c.JSON(200, gin.H{"message": fmt.Sprintf("user with email %s created", ju.Email)})
}
