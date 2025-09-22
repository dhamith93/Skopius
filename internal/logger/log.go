package logger

import (
	"log"
	"os"

	"github.com/dhamith93/Skopius/internal/models"
)

var debug = os.Getenv("SKOPIUS_DEBUG") == "1"

func LogResult(res models.CheckResult) {
	if debug {
		// Detailed breakdown
		log.Printf("[%s] %s (code=%d, total=%dms, dns=%dms, connect=%dms, tls=%dms, ttfb=%dms, server=%dms, err=%s)",
			res.Service, res.Status, res.Code,
			res.Total, res.DNS, res.Connect, res.TLS, res.TTFB, res.Server, res.Error)
	} else {
		// Compact summary
		log.Printf("[%s] %s (code=%d, total=%dms, err=%s)",
			res.Service, res.Status, res.Code, res.Total, res.Error)
	}
}
