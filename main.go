package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rotemtam/ent-bank-io/ent"
)

const (
	bearerKey = "bearerKey"
)

func main() {
	// create a new ent client with an in memory database.
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Panicf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Panicf("failed creating schema resources: %v", err)
	}
	// Start the server
	if err := newServer(client).Run(); err != nil {
		log.Panicf("failed: %v: ", err)
	}
}

type server struct {
	client *ent.Client
	*gin.Engine
}

// newServer creates a new ent-bank server.
func newServer(client *ent.Client) *server {
	r := gin.Default()
	s := &server{client: client, Engine: r}
	r.Use(extractBearer)
	r.POST("/v1/user", s.createUser)
	r.PATCH("/v1/user/:id/balance", s.updateBalance)
	r.GET("/v1/user/:id/balance/:timestamp", s.balanceAt)
	return s
}

// createUser creates a new user.
func (s *server) createUser(c *gin.Context) {
	var payload ent.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user, err := s.client.User.Create().
		SetName(payload.Name).
		SetEmail(payload.Email).
		SetBalance(payload.Balance).
		Save(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("user created successfully by %s", c.Value(bearerKey))
	c.JSON(http.StatusOK, gin.H{"id": user.ID})
}

// updateBalance updates the balance of a User.
func (s *server) updateBalance(c *gin.Context) {
	var payload struct {
		Delta float64 `json:"delta"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user, err := s.client.User.UpdateOneID(id).
		AddBalance(payload.Delta).
		Save(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": user.Balance})
}

// TODO: Implement
func (s *server) balanceAt(ctx *gin.Context) {
	panic("not implemented")
}

// extractBearer extracts the authorization token from the `Authorization` header and
// places it on the *gin.Context. .
func extractBearer(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if len(token) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		c.Abort()
		return
	}
	c.Set(bearerKey, token)
	c.Next()
}
