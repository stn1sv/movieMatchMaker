package middleware

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

func Logging(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/getAllParticipants" {
			c.Next()
			return
		}
		log = log.With(
			slog.String("method:", c.Request.Method),
			slog.String("path:", c.Request.URL.Path),
			slog.String("remote addr:", c.Request.RemoteAddr),
		)

		start := time.Now()
		defer func() {
			duration := time.Since(start).String()
			if len(c.Errors) > 0 {
				for _, err := range c.Errors {
					log.Error("request failed",
						slog.Int("status", c.Writer.Status()),
						slog.String("duration", duration),
						slog.String("error", err.Error()),
					)
				}
			} else {
				log.Info("request completed",
					slog.Int("status", c.Writer.Status()),
					slog.String("duration", duration),
				)
			}
		}()
		c.Next()
	}
}
