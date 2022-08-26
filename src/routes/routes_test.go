package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRoutes(t *testing.T) {
	// initialize gin router
	router := gin.Default()
	Routes(router)
}
