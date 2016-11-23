package main

import (
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"github.com/apourchet/incipit/lib/auth"
	"github.com/apourchet/incipit/lib/healthz"
	"github.com/apourchet/incipit/lib/logger"
	"github.com/apourchet/incipit/lib/utils"
	protos "github.com/apourchet/incipit/protos/go"
	"github.com/gin-gonic/gin"
	glog "github.com/golang/glog"
)

type AuthService struct {
	client auth.AuthClient
	logger protos.LoggerClient
}

const (
	ApiPort = 10000
)

func NewAuthService(client auth.AuthClient, logger protos.LoggerClient) *AuthService {
	return &AuthService{client, logger}
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
	glog.Infof("Service: UserExists")
	key := c.Query("key")
	found, err := s.client.UserExists(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"found": found})
}

func (s *AuthService) Register(c *gin.Context) {
	glog.Infof("Service: Register")
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

	ok, err := s.client.Register(key, pass)
	if utils.InternalError(c, err) {
		return
	}

	if !ok {
		c.JSON(http.StatusConflict, gin.H{"error": "Username taken"})
		return
	}

	c.Status(http.StatusOK)
}

func (s *AuthService) Login(c *gin.Context) {
	glog.Infof("Service: Login")
	key := c.Query("key")
	pass := c.Query("pass")

	token, valid, err := s.client.Login(key, pass)
	if utils.InternalError(c, err) {
		return
	}
	if !valid {
		glog.Infof("Returning InvalidCredentialsError")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "InvalidCredentialsError"})
		return
	}

	go s.logger.LogLogin(context.Background(), &protos.LogLoginReq{key})

	auth.SetToken(c, token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *AuthService) Logout(c *gin.Context) {
	glog.Infof("Service: Logout")
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
	glog.Infof("Service: Deregister")
	token, err := auth.GetToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	_, valid, err := s.client.Validate(token)
	if utils.InternalError(c, err) {
		return
	}
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "InvalidTokenError"})
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

	authClient := auth.GetDefaultClient()
	loggerClient, err := logger.GetDefaultClient()
	if err != nil {
		glog.Fatalf("Fatal error: %v", err)
	}
	err = NewAuthService(authClient, loggerClient).Run()
	glog.Fatal(err)
}
