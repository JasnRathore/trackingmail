package emailtracker

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type OpenEvent struct {
	ID            string
	IP            string
	XForwardedFor string
	UserAgent     string
	Referer       string
	AcceptLang    string
	Time          time.Time
}

type Config struct {
	Port   int    // Port to listen on, e.g., 8080
	Domain string // Domain or host, e.g., "localhost:8080" or "tracker.example.com"
	Path   string // Tracking pixel path, e.g., "/pixel"
}

type Tracker struct {
	config   Config
	callback func(OpenEvent)
}

func NewTracker(cfg Config, cb func(OpenEvent)) *Tracker {
	return &Tracker{
		config:   cfg,
		callback: cb,
	}
}

func (t *Tracker) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event := OpenEvent{
			ID:            r.URL.Query().Get("id"),
			IP:            getIP(r),
			XForwardedFor: r.Header.Get("X-Forwarded-For"),
			UserAgent:     r.Header.Get("User-Agent"),
			Referer:       r.Header.Get("Referer"),
			AcceptLang:    r.Header.Get("Accept-Language"),
			Time:          time.Now(),
		}
		if t.callback != nil {
			t.callback(event)
		}
		w.Header().Set("Content-Type", "image/gif")
		w.WriteHeader(http.StatusOK)
		w.Write(pixelData)
	}
}

func (t *Tracker) Start() error {
	http.HandleFunc(t.config.Path, t.Handler())
	addr := fmt.Sprintf(":%d", t.config.Port)
	return http.ListenAndServe(addr, nil)
}

func (t *Tracker) GenerateLink(id string) string {
	protocol := "https"
	// Use http for localhost or 127.0.0.1
	if t.config.Domain == "localhost" ||
		t.config.Domain == "localhost:"+fmt.Sprint(t.config.Port) ||
		t.config.Domain == "127.0.0.1" ||
		t.config.Domain == "127.0.0.1:"+fmt.Sprint(t.config.Port) {
		protocol = "http"
	}
	return fmt.Sprintf("%s://%s%s?id=%s", protocol, t.config.Domain, t.config.Path, id)
}

func getIP(r *http.Request) string {
	// If behind proxy, prefer X-Forwarded-For
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// 1x1 Transparent Gif
var pixelData = []byte{
	0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00,
	0x01, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00,
	0xFF, 0xFF, 0xFF, 0x21, 0xF9, 0x04, 0x01, 0x00,
	0x00, 0x00, 0x00, 0x2C, 0x00, 0x00, 0x00, 0x00,
	0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x02, 0x44,
	0x01, 0x00, 0x3B,
}
