package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HandleAllUsers(c *gin.Context) {

	allUsers, err := s.store.AllUsers(context.Background())

	if err != nil {
		fmt.Println("Something went wrong")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	fmt.Printf("Result: %T %#v, isNil=%v\n", allUsers, allUsers, allUsers == nil)

	c.JSON(http.StatusOK, allUsers)

}
