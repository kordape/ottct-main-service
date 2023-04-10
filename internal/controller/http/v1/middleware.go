package v1

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/pkg/token"
	"github.com/sirupsen/logrus"
)

const authHeaderKey = "Authorization"

const ctxLoggerKey = "logger"

func AuthMiddleware(tokenManager *token.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader(authHeaderKey)
		token := strings.Split(bearerToken, "Bearer ")

		if err := tokenManager.VerifyJWT(token[1]); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func Logging(logger *logrus.Entry) gin.HandlerFunc {
	return func(c *gin.Context) {
		if logger == nil {
			logger = defaultLogger()
		}
		// Set path field
		logger = logger.WithFields(logrus.Fields{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		})
		// Save logger to context
		c.Set(ctxLoggerKey, logger)
	}
}

// getLogger returns a logrus entry from the gin.Context
// Always returns a logger.
func getLogger(c *gin.Context) *logrus.Entry {
	value, ok := c.Get(ctxLoggerKey)
	if !ok {
		return defaultLogger()
	}

	logger, ok := value.(*logrus.Entry)
	if !ok {
		return defaultLogger()
	}

	return logger
}

func defaultLogger() *logrus.Entry {
	log := logrus.StandardLogger()
	log.SetFormatter(
		&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		},
	)

	return logrus.NewEntry(log).WithField("foo", "bar")
}
