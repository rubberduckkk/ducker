package proxy

import (
	"fmt"
	"io"
	"net/http"
)

func (s *service) ProxyHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.authHTTP(r) {
		w.Header().Set("Content-Type", "application/plain")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, "ducker-proxy: unauthorized")
		return
	}

	// Create a new HTTP request with the same method, URL, and body as the original request
	targetURL := r.URL
	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	// Copy the headers from the original request to the proxy request
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// Send the proxy request using the custom transport
	resp, err := http.DefaultTransport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the headers from the proxy response to the original response
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set the status code of the original response to the status code of the proxy response
	w.WriteHeader(resp.StatusCode)

	// Copy the body of the proxy response to the original response
	_, _ = io.Copy(w, resp.Body)
}
