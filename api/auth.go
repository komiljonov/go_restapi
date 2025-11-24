package api

import (
	"context"
	"net/http"
	db "restapi/db/sqlc"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name      string `json:"name" binding:"required"`
	Password  string `json:"password" binding:"required,min=8"`
	Birthdate string `json:"birthdate" binding:"required,datetime=2006-01-02"`
}

func (s *Server) RegisterHandler(c *gin.Context) {

	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := s.pool.Begin(context.Background())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction in pool"})
	}

	// query :=  





}
