package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/apourchet/incipit/lib/auth"
	"github.com/apourchet/incipit/lib/healthz"
	"github.com/apourchet/incipit/lib/utils"
	"github.com/gin-gonic/gin"
	glog "github.com/golang/glog"
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
	engine.Handle("POST", "/api/v1/auth/register", s.Register)
	engine.Handle("GET", "/api/v1/auth/login", s.Login)
	engine.Handle("POST", "/api/v1/auth/logout", s.Logout)
	engine.Handle("POST", "/api/v1/auth/deregister", s.Deregister)
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
	var form map[string]string
	if utils.InternalError(c, c.Bind(&form)) {
		return
	}

	key := form["key"]
	pass := form["pass"]

	found, err := s.client.UserExists(key)
	if utils.InternalError(c, err) {
		return
	}

	if found {
		glog.Errorf("400: %s", key)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is taken"})
		return
	}

	err = s.client.Register(key, pass)
	if utils.InternalError(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

func (s *AuthService) Login(c *gin.Context) {
	key := c.Query("key")
	pass := c.Query("pass")

	token, err := s.client.Login(key, pass)
	if utils.InternalError(c, err) {
		return
	}

	auth.SetToken(c, token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *AuthService) Logout(c *gin.Context) {
	token, err := auth.GetToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	err = s.client.Logout(token)
	if utils.InternalError(c, err) {
		return
	}
	c.Status(http.StatusOK)
}

func (s *AuthService) Deregister(c *gin.Context) {
	token, err := auth.GetToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	_, err = s.client.Validate(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = s.client.Deregister(token)
	if utils.InternalError(c, err) {
		return
	}
	c.Status(http.StatusOK)
}

func main() {
	flag.Parse()
	healthz.SpawnHealthCheck(healthz.DefaultPort)

	authClient, err := auth.NewK8sAuthClient()
	if err != nil {
		glog.Fatalf("Failed to get k8s auth client %v", err)
	}
	err = NewAuthService(authClient).Run()
	log.Fatal(err)
}
