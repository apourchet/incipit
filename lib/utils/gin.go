package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func InternalError(c *gin.Context, err error) bool {
	if err != nil {
		glog.Errorf("500: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return true
	}
	return false
}

func BadRequest(c *gin.Context, err error) bool {
	if err != nil {
		glog.Errorf("400: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return true
	}
	return false
}
