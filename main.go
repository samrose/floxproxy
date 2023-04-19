package main

import (
    "crypto/tls"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Log the incoming request.
        dump, err := httputil.DumpRequest(r, true)
        if err != nil {
            log.Printf("Error dumping request: %v", err)
        } else {
            log.Printf("Incoming request: %v", string(dump))
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

    log.Println("Listening on :8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Error listening: %v", err)
    }
}
