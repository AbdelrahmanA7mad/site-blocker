package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func handleProxyRequest(w http.ResponseWriter, r *http.Request) {
	// Intercept local restore requests from a new instance
	if (r.Host == "127.0.0.1:8080" || r.Host == "localhost:8080") && r.URL.Path == "/restore" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Restoring control to new terminal instance..."))
		go func() {
			time.Sleep(100 * time.Millisecond)
			close(shutdownChan)
		}()
		return
	}

	host := cleanHost(r.Host)

	if isBlocked(host) {
		log.Println("[BLOCKED]", host)
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Focus Proxy: BLOCKED"))
		return
	}

	if r.Method == http.MethodConnect {
		handleHTTPSConnect(w, r)
		return
	}

	log.Println("[ALLOW HTTP]", host)

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func handleHTTPSConnect(w http.ResponseWriter, r *http.Request) {
	host := cleanHost(r.Host)
	if isBlocked(host) {
		log.Println("[BLOCKED HTTPS]", host)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	dest, err := net.Dial("tcp", r.Host)
	if err != nil {
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		_ = dest.Close()
		return
	}

	client, _, err := hijacker.Hijack()
	if err != nil {
		_ = dest.Close()
		return
	}

	_, _ = client.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	go transfer(dest, client)
	go transfer(client, dest)
}

func transfer(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()
	_, _ = io.Copy(dst, src)
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
