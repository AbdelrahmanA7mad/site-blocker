package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

const (
	proxyAddr     = "127.0.0.1:8080"
	hostsFilePath = "blocked_hosts.json"
)

var defaultBlockedHosts = []string{
	"facebook.com",
	"instagram.com",
	"tiktok.com",
	"x.com",
	"twitter.com",
}

var (
	blockedHosts   []string
	blockedHostsMu sync.RWMutex
)

// initHostsFile checks if the JSON config file exists, loads it if it does,
// or creates it with default hosts if it doesn't.
func initHostsFile() error {
	blockedHostsMu.Lock()
	defer blockedHostsMu.Unlock()

	// Check if file exists
	if _, err := os.Stat(hostsFilePath); os.IsNotExist(err) {
		// File does not exist, create it with defaults
		blockedHosts = append([]string(nil), defaultBlockedHosts...)
		data, err := json.MarshalIndent(blockedHosts, "", "  ")
		if err != nil {
			return err
		}
		err = os.WriteFile(hostsFilePath, data, 0644)
		if err != nil {
			return err
		}
		log.Println("Created default blocked_hosts.json")
		return nil
	}

	// File exists, read it
	data, err := os.ReadFile(hostsFilePath)
	if err != nil {
		return err
	}

	var loadedHosts []string
	if err := json.Unmarshal(data, &loadedHosts); err != nil {
		// If JSON is invalid, backup and reset to defaults
		log.Println("Warning: invalid JSON in blocked_hosts.json. Resetting to defaults.")
		blockedHosts = append([]string(nil), defaultBlockedHosts...)
		_ = saveHostsToFileLocked()
		return err
	}

	// Clean and normalize hosts loaded from file
	cleanedHosts := []string{}
	seen := make(map[string]bool)
	for _, h := range loadedHosts {
		ch := cleanHost(h)
		if ch != "" && !seen[ch] {
			seen[ch] = true
			cleanedHosts = append(cleanedHosts, ch)
		}
	}

	blockedHosts = cleanedHosts
	return nil
}

// saveHostsToFileLocked saves current blocked hosts to the JSON file.
// Callers must hold blockedHostsMu Lock.
func saveHostsToFileLocked() error {
	data, err := json.MarshalIndent(blockedHosts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(hostsFilePath, data, 0644)
}

func cleanHost(host string) string {
	host = strings.TrimSpace(strings.ToLower(host))
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
	_ = saveHostsToFileLocked()
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
			_ = saveHostsToFileLocked()
			return true
		}
	}
	return false
}
