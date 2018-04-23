package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/technosolutionscl/trumail/api"
	"github.com/technosolutionscl/trumail/config"
	"github.com/technosolutionscl/trumail/verifier"
)

func main() {
	// Generate a new logrus logger
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	l := logger.WithField("port", config.Port)

	// Define all required dependencies
	l.Info("Defining all service dependencies")
	e := echo.New()
	v := verifier.NewVerifier(retrievePTR(), config.SourceAddr)
	s := api.NewService(logger, config.HTTPClientTimeout, v)

	// Bind endpoints to router
	l.Info("Binding API endpoints to the router")
	if config.RateLimitHours != 0 && config.RateLimitMax != 0 {
		r := api.NewRateLimiter(config.RateLimitMax,
			time.Hour*time.Duration(config.RateLimitHours), config.RateLimitCIDRCustom)
		e.GET("/v1/:format/:email", s.Lookup, r.RateLimit)
		e.GET("/limit-status", r.LimitStatus)
	} else {
		e.GET("/v1/:format/:email", s.Lookup)
	}
	e.GET("/v1/health", s.Health)

	// Listen and Serve
	l.WithField("port", config.Port).Info("Listening and Serving")
	l.Fatal(e.Start(":" + config.Port))
}

// retrievePTR attempts to retrieve the PTR record for the IP
// address retrieved via an API call on api.ipify.org
func retrievePTR() string {
	transport := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, Proxy: http.ProxyFromEnvironment}
	client := &http.Client{Transport: &transport}

	// Request the IP from ipify
	resp, err := client.Get("https://api.ipify.org/")
	if err != nil {
		log.Fatal("Failed to retrieve IP from api.ipify.org")
	}
	defer resp.Body.Close()

	// Decodes the IP response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read IP response body")
	}

	// Retrieve the PTR record for our IP and return without a trailing dot
	names, err := net.LookupAddr(string(data))
	if err != nil {
		return string(data)
	}
	return strings.TrimSuffix(names[0], ".")
}
