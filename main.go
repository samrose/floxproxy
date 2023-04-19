package main

import (
    "crypto/tls"
    "os"
    "net/http"
    "net/http/httputil"
    "net/url"

    log "github.com/sirupsen/logrus"
)

func main() {
    // Create a file for logging.
    file, err := os.OpenFile("server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatalf("Error opening log file: %v", err)
    }
    defer file.Close()

    // Configure the logger to write to the file.
    log.SetOutput(file)

    // Set the log level to Info to only log important messages.
    log.SetLevel(log.InfoLevel)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Log the incoming request.
        dump, err := httputil.DumpRequest(r, true)
        if err != nil {
            log.Errorf("Error dumping request: %v", err)
        } else {
            log.Infof("Incoming request: %v", string(dump))
        }

        // Parse the request header to extract the "target" value.
        target := r.Header.Get("X-Target")

        // Create a transport with the InsecureSkipVerify option enabled.
        transport := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }

        // Create a reverse proxy with the "target" and "transport" options.
        targetURL, _ := url.Parse(target)
        reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
        reverseProxy.Transport = transport

        // Serve the request using the reverse proxy.
        reverseProxy.ServeHTTP(w, r)
    })

    //log.Info("Listening on :8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Error listening: %v", err)
    }
}
