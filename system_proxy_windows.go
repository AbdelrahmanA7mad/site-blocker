package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

const internetSettingsPath = `Software\Microsoft\Windows\CurrentVersion\Internet Settings`

func spawnBackgroundClone(enabled bool) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command(exePath, "-background", fmt.Sprintf("-enabled=%t", enabled))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
	return cmd.Start()
}


var (
	proxyActive   bool
	proxyActiveMu sync.RWMutex
)

func init() {
	proxyActiveMu.Lock()
	proxyActive = checkSystemProxyState()
	proxyActiveMu.Unlock()
}

func checkSystemProxyState() bool {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		internetSettingsPath,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return false
	}
	defer key.Close()

	val, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil {
		return false
	}
	return val == 1
}

func isProxyActive() bool {
	proxyActiveMu.RLock()
	defer proxyActiveMu.RUnlock()
	return proxyActive
}

func enableSystemProxy() {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		internetSettingsPath,
		registry.ALL_ACCESS,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer key.Close()

	_ = key.SetDWordValue("ProxyEnable", 1)
	_ = key.SetStringValue("ProxyServer", proxyAddr)

	proxyActiveMu.Lock()
	proxyActive = true
	proxyActiveMu.Unlock()

	log.Println("SYSTEM PROXY ENABLED")
}

func disableSystemProxy() {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		internetSettingsPath,
		registry.ALL_ACCESS,
	)
	if err != nil {
		return
	}
	defer key.Close()

	_ = key.SetDWordValue("ProxyEnable", 0)

	proxyActiveMu.Lock()
	proxyActive = false
	proxyActiveMu.Unlock()

	log.Println("SYSTEM PROXY DISABLED")
}

func resetSystemProxy() {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		internetSettingsPath,
		registry.ALL_ACCESS,
	)
	if err != nil {
		log.Println("Failed to open registry key:", err)
		return
	}
	defer key.Close()

	_ = key.SetDWordValue("ProxyEnable", 0)
	_ = key.DeleteValue("ProxyServer")
	_ = key.DeleteValue("ProxyOverride")

	proxyActiveMu.Lock()
	proxyActive = false
	proxyActiveMu.Unlock()

	log.Println("SYSTEM PROXY RESET TO DEFAULT")
}


