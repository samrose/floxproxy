package main

import (
    "crypto/tls"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        target := r.Header.Get("X-Target")

        transport := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }

        targetURL, _ := url.Parse(target)
        reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
        reverseProxy.Transport = transport

        // Serve the request using the reverse proxy.
        reverseProxy.ServeHTTP(w, r)
    })

    http.ListenAndServe(":8080", nil)
}
