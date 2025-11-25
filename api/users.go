package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HandleAllUsers(c *gin.Context) {

	allUsers, err := s.store.AllUsers(c.Request.Context())

	if err != nil {
		fmt.Println("Something went wrong")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Result: %T %#v, isNil=%v\n", allUsers, allUsers, allUsers == nil)

	c.JSON(http.StatusOK, allUsers)

}
