package api

import (
	"errors"
	"net/http"
	db "restapi/db/sqlc"
	"restapi/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (s *Server) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			// if the Auth header is missing, just skip authentication
			// Because we have register, login routes too
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		token, err := utils.ParseAndVerifyJWT(parts[1])

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		//conn, err := s.pool.Acquire(c.Request.Context())
		//
		//if err != nil {
		//	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//	return
		//}

		//defer conn.Release()

		//q := db.New(conn)

		userId, err := strconv.ParseInt(token.UserID, 10, 32)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//user, err := q.GetUser(c.Request.Context(), int32(userId))

		user, err := s.store.GetUser(c.Request.Context(), int32(userId))

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return

		}

		c.Set("user", user)

		c.Next()

	}
}

func (s *Server) AuthRequired(handler func(c *gin.Context, user db.User)) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("user")

		if !ok || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		handler(c, user.(db.User))

	}
}
