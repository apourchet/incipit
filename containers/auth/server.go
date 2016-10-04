package main

import (
	"fmt"
	"net/http"

	"github.com/apourchet/incipit/lib/auth"
	"github.com/apourchet/incipit/lib/healthz"
	"github.com/apourchet/incipit/lib/logging"
	"github.com/gin-gonic/gin"
)

type AuthService struct {
	client auth.AuthClient
}

const (
	ApiPort = 8080
)

func NewAuthService(client auth.AuthClient) *AuthService {
	return &AuthService{client}
}

func (s *AuthService) Run() error {
	engine := gin.New()
	engine.Handle("GET", "/api/v1/auth/userexists", s.UserExists)
	engine.Handle("GET", "/api/v1/auth/register", s.Register)
	engine.Handle("GET", "/api/v1/auth/login", s.Login)
	engine.Handle("GET", "/api/v1/auth/logout", s.Logout)
	engine.Handle("GET", "/api/v1/auth/deregister", s.Deregister)
	return engine.Run(fmt.Sprintf(":%d", ApiPort))
}

func (s *AuthService) UserExists(c *gin.Context) {
	key := c.Query("key")
	found, err := s.client.UserExists(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"found": found})
}

func (s *AuthService) Register(c *gin.Context) {
	key := c.Query("key")
	pass := c.Query("pass")
	err := s.client.Register(key, pass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *AuthService) Login(c *gin.Context) {
	key := c.Query("key")
	pass := c.Query("pass")
	token, err := s.client.Login(key, pass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	auth.SetToken(c, token)
	c.Status(http.StatusOK)
}

func (s *AuthService) Logout(c *gin.Context) {
	token, err := auth.GetToken(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	err = s.client.Logout(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (s *AuthService) Deregister(c *gin.Context) {
	token, err := auth.GetToken(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	err = s.client.Deregister(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func main() {
	healthz.SpawnHealthCheck(healthz.DefaultPort)

	err := NewAuthService(auth.NewMockAuthClient()).Run()
	logging.Fatal(err)
}
