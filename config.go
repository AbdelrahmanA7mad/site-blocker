package main

import (
	"sort"
	"strings"
	"sync"
)

const proxyAddr = "127.0.0.1:8080"

var defaultBlockedHosts = []string{
	"facebook.com",
	"instagram.com",
	"tiktok.com",
	"x.com",
	"twitter.com",
}

var (
	blockedHosts   = append([]string(nil), defaultBlockedHosts...)
	blockedHostsMu sync.RWMutex
)

func cleanHost(host string) string {
	host = strings.ToLower(host)
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	return host
}

func isBlocked(host string) bool {
	host = cleanHost(host)
	blockedHostsMu.RLock()
	defer blockedHostsMu.RUnlock()

	for _, blocked := range blockedHosts {
		if host == blocked || strings.HasSuffix(host, "."+blocked) {
			return true
		}
	}
	return false
}

func listBlockedHosts() []string {
	blockedHostsMu.RLock()
	defer blockedHostsMu.RUnlock()

	out := append([]string(nil), blockedHosts...)
	sort.Strings(out)
	return out
}

func addBlockedHost(host string) bool {
	host = cleanHost(host)
	if host == "" {
		return false
	}

	blockedHostsMu.Lock()
	defer blockedHostsMu.Unlock()

	for _, blocked := range blockedHosts {
		if blocked == host {
			return false
		}
	}
	blockedHosts = append(blockedHosts, host)
	return true
}

func removeBlockedHost(host string) bool {
	host = cleanHost(host)
	if host == "" {
		return false
	}

	blockedHostsMu.Lock()
	defer blockedHostsMu.Unlock()

	for i, blocked := range blockedHosts {
		if blocked == host {
			blockedHosts = append(blockedHosts[:i], blockedHosts[i+1:]...)
			return true
		}
	}
	return false
}
