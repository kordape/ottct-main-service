package integration_test

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/Eun/go-hit"
)

const (
	// Attempts connection
	host       = "localhost:8080"
	healthPath = "http://" + host + "/health"
	attempts   = 20

	// HTTP REST
	basePath = "http://" + host + "/v1"
)

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

// HTTP GET: v1/echo.
func TestEcho(t *testing.T) {
	Test(t,
		Description("Echo Success"),
		Get(basePath+"/echo"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains("Echo Response"),
	)
}
