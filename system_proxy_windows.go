package main

import (
	"log"

	"golang.org/x/sys/windows/registry"
)

const internetSettingsPath = `Software\Microsoft\Windows\CurrentVersion\Internet Settings`

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
	log.Println("SYSTEM PROXY DISABLED")
}
