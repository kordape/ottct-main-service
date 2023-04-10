package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/pkg/token"
	"github.com/sirupsen/logrus"
)

const authHeaderKey = "Authorization"

const ctxLoggerKey = "logger"

func AuthMiddleware(tokenManager *token.Manager, log *logrus.Entry) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader(authHeaderKey)
		token := strings.Split(bearerToken, "Bearer ")

		if err := tokenManager.VerifyJWT(token[1]); err != nil {
			log.WithError(err).Warn("Failed to verify token")
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
		logger.Debug("Request received")
		// Save logger to context
		c.Set(ctxLoggerKey, logger)

		// run next middleware in chain
		c.Next()

		// log response code
		responseCode := c.Writer.Status()
		if responseCode != http.StatusOK {
			logger.WithField("code", responseCode).Error("Request failed")
		} else {
			logger.Debug("Request successful")
		}
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
	logrus.SetReportCaller(true)
	logrus.SetFormatter(
		&logrus.TextFormatter{
			ForceColors: true,
		},
	)

	return logrus.NewEntry(log)
}
