package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) Ping(c *gin.Context) {
	start := time.Now()
	log.Println("START", start.Format(time.RFC3339Nano))

	stats := s.pool.Stat()

	log.Printf(
		"POOL stats: total=%d idle=%d acquired=%d max=%d "+
			"acquireCount=%d emptyAcquire=%d cancellations=%d",
		stats.TotalConns(),
		stats.IdleConns(),
		stats.AcquiredConns(),
		stats.MaxConns(),
		stats.AcquireCount(),
		stats.EmptyAcquireCount(),
		stats.CanceledAcquireCount(),
	)

	log.Println("END  ", time.Now().Format(time.RFC3339Nano))
	c.String(http.StatusOK, "Salom")
}
