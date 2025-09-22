package monitor

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/dhamith93/Skopius/internal/models"
)

type Service struct {
	Name     string        `yaml:"name"`
	URL      string        `yaml:"url"`
	Interval time.Duration `yaml:"interval"`
	MaxRead  int64         `yaml:"max_read"` // bytes to read, 0 = default
}

func (s *Service) Check() models.CheckResult {
	result := s.checkHTTP(s.URL)
	return result
}

func (s *Service) checkHTTP(url string) models.CheckResult {
	var result models.CheckResult
	result.Service = s.Name
	result.URL = s.URL
	result.Timestamp = time.Now().UTC()

	var start, dnsStart, dnsDone, connStart, connDone, tlsStart, tlsDone, firstByte time.Time

	req, _ := http.NewRequest("GET", s.URL, nil)
	trace := &httptrace.ClientTrace{
		DNSStart:             func(info httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone:              func(info httptrace.DNSDoneInfo) { dnsDone = time.Now() },
		ConnectStart:         func(network, addr string) { connStart = time.Now() },
		ConnectDone:          func(network, addr string, err error) { connDone = time.Now() },
		TLSHandshakeStart:    func() { tlsStart = time.Now() },
		TLSHandshakeDone:     func(cs tls.ConnectionState, err error) { tlsDone = time.Now() },
		GotFirstResponseByte: func() { firstByte = time.Now() },
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	start = time.Now()
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		result.Status = "down"
		result.Error = err.Error()
		return result
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	// default maxRead if not set
	maxRead := s.MaxRead
	if maxRead <= 0 {
		maxRead = 64 * 1024 // default 64KB
	}

	if resp != nil && resp.Body != nil {
		_, _ = io.CopyN(io.Discard, resp.Body, maxRead)
	}

	total := time.Since(start)

	if resp != nil {
		result.Code = resp.StatusCode
		result.Status = "up"
	}

	// durations
	if !dnsStart.IsZero() && !dnsDone.IsZero() {
		result.DNS = dnsDone.Sub(dnsStart).Milliseconds()
	}
	if !connStart.IsZero() && !connDone.IsZero() {
		result.Connect = connDone.Sub(connStart).Milliseconds()
	}
	if !tlsStart.IsZero() && !tlsDone.IsZero() {
		result.TLS = tlsDone.Sub(tlsStart).Milliseconds()
	}
	if !firstByte.IsZero() {
		result.TTFB = firstByte.Sub(start).Milliseconds()
		result.Server = total.Milliseconds() - result.TTFB
	}
	result.Total = total.Milliseconds()

	result.Received = time.Now().UTC()
	return result
}

func (s *Service) Probe(ctx context.Context, results chan<- models.CheckResult) {
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			res := s.Check()
			results <- res
		case <-ctx.Done():
			log.Printf("Probe stopped for %s\n", s.Name)
			return
		}
	}
}
