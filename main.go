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
    file, err := os.OpenFile("server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatalf("Error opening log file: %v", err)
    }
    defer file.Close()

    log.SetOutput(file)

    log.SetLevel(log.InfoLevel)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Log the incoming request.
        dump, err := httputil.DumpRequest(r, true)
        if err != nil {
            log.Errorf("Error dumping request: %v", err)
        } else {
            log.Infof("Incoming request: %v", string(dump))
        }

        target := r.Header.Get("X-Target")

        transport := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }

        targetURL, _ := url.Parse(target)
        reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
        reverseProxy.Transport = transport

        reverseProxy.ServeHTTP(w, r)
    })

    //log.Info("Listening on :8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Error listening: %v", err)
    }
}
