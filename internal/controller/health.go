package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func (h *Controller) HealthCheckController(c *gin.Context) {
  // REturn a 200 OK response
  c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
