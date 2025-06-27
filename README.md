
# TrackingMail 📧

Ever wondered if your carefully crafted emails are actually being opened? TrackingMail is a lightweight Go library that makes email open tracking simple and straightforward. By embedding an invisible 1x1 pixel image in your emails, you can capture valuable insights about recipient engagement without compromising the user experience.

## Project Requirements

- **Go Version**: 1.24.4 or higher
- **Network Access**: Ability to serve HTTP requests on your chosen port
- **Domain/Host**: A publicly accessible domain or host for the tracking pixel

## Dependencies

This library uses only Go's standard library packages:
- `net/http` - HTTP server functionality
- `net` - Network utilities for IP extraction
- `time` - Timestamp generation
- `fmt` - String formatting

No external dependencies required! 🎉

## Getting Started

### Basic Setup

The core concept is simple: when someone opens an email containing your tracking pixel, their email client makes an HTTP request to fetch the image. This request contains valuable metadata that gets captured and processed through your callback function.

### Configuration

Create a `Config` struct to define your tracking server settings:

```go
config := emailtracker.Config{
    Port:   8080,
    Domain: "tracker.example.com",
    Path:   "/pixel",
}
```

### Initialize the Tracker

Set up your tracker with a callback function to handle open events:

```go
tracker := emailtracker.NewTracker(config, func(event emailtracker.OpenEvent) {
    fmt.Printf("Email opened! ID: %s, IP: %s, Time: %s\n", 
        event.ID, event.IP, event.Time.Format(time.RFC3339))
})
```

## How to Run the Application

### Simple Server Example

Here's a complete example to get your tracking server up and running:

```go
package main

import (
    "fmt"
    "log"
    "github.com/jasnrathore/trackingmail"
)

func main() {
    config := emailtracker.Config{
        Port:   8080,
        Domain: "localhost:8080",
        Path:   "/track",
    }
    
    tracker := emailtracker.NewTracker(config, handleEmailOpen)
    
    fmt.Printf("Tracking server starting on port %d\n", config.Port)
    fmt.Printf("Tracking URL: %s\n", tracker.GenerateLink("test-email-001"))
    
    if err := tracker.Start(); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}

func handleEmailOpen(event emailtracker.OpenEvent) {
    fmt.Printf("📧 Email Opened!\n")
    fmt.Printf("   ID: %s\n", event.ID)
    fmt.Printf("   IP: %s\n", event.IP)
    fmt.Printf("   User Agent: %s\n", event.UserAgent)
    fmt.Printf("   Time: %s\n", event.Time.Format("2006-01-02 15:04:05"))
    fmt.Printf("   Referer: %s\n", event.Referer)
    fmt.Println("---")
}
```

### Production Deployment

For production environments, ensure your domain is properly configured:

```go
config := emailtracker.Config{
    Port:   80,  // or 443 for HTTPS
    Domain: "tracking.yourcompany.com",
    Path:   "/pixel",
}
```

## Relevant Examples

### Database Integration

Store tracking events in a database for analytics:

```go
func saveToDatabase(event emailtracker.OpenEvent) {
    // Insert into your preferred database
    query := `INSERT INTO email_opens 
              (email_id, ip_address, user_agent, opened_at) 
              VALUES (?, ?, ?, ?)`
    
    db.Exec(query, event.ID, event.IP, event.UserAgent, event.Time)
}

tracker := emailtracker.NewTracker(config, saveToDatabase)
```

### Email Template Integration

Embed the tracking pixel in your HTML emails:

```go
emailID := "campaign-123-user-456"
trackingURL := tracker.GenerateLink(emailID)

htmlTemplate := fmt.Sprintf(`
<html>
<body>
    <h1>Welcome to our newsletter!</h1>
    <p>Thanks for subscribing to our updates.</p>
    
    <!-- Invisible tracking pixel -->
    <img src="%s" width="1" height="1" style="display:none;" />
</body>
</html>
`, trackingURL)
```

### Advanced Event Processing

Handle different types of tracking scenarios:

```go
func advancedEventHandler(event emailtracker.OpenEvent) {
    // Parse campaign and user info from ID
    parts := strings.Split(event.ID, "-")
    if len(parts) >= 3 {
        campaign := parts[1]
        userID := parts[2]
        
        // Update campaign metrics
        updateCampaignStats(campaign)
        
        // Record user engagement
        recordUserActivity(userID, event.Time)
    }
    
    // Detect mobile vs desktop
    if strings.Contains(event.UserAgent, "Mobile") {
        recordMobileOpen(event.ID)
    }
    
    // Geographic tracking (if you have IP geolocation)
    location := getLocationFromIP(event.IP)
    recordGeoData(event.ID, location)
}
```

### Custom Handler with Middleware

Add your own HTTP middleware for enhanced functionality:

```go
func customHandler(tracker *emailtracker.Tracker) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Add CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        
        // Rate limiting check
        if !checkRateLimit(r.RemoteAddr) {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        
        // Call the original handler
        tracker.Handler()(w, r)
    }
}

// Register with custom handler
http.HandleFunc(config.Path, customHandler(tracker))
```

### Batch Processing

Process events in batches for better performance:

```go
type BatchProcessor struct {
    events   []emailtracker.OpenEvent
    ticker   *time.Ticker
    mu       sync.Mutex
}

func (bp *BatchProcessor) addEvent(event emailtracker.OpenEvent) {
    bp.mu.Lock()
    defer bp.mu.Unlock()
    
    bp.events = append(bp.events, event)
    
    if len(bp.events) >= 100 {
        bp.processBatch()
    }
}

func (bp *BatchProcessor) processBatch() {
    // Process all events in batch
    for _, event := range bp.events {
        // Save to database, send to analytics, etc.
    }
    bp.events = bp.events[:0] // Clear slice
}
```

## Key Features

**Comprehensive Event Data**: Capture IP addresses, user agents, referrers, timestamps, and custom identifiers for detailed analytics.

**Flexible Configuration**: Easy setup for both development and production environments with automatic protocol detection.

**Lightweight Implementation**: Built with Go's standard library, ensuring minimal overhead and maximum performance.

**Custom Callback Support**: Process tracking events exactly how your application needs them.


Ready to start tracking? Set up your configuration, define your callback function, and begin gaining valuable insights into your email performance today!