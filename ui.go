package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func runTerminalUI() {
	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		fmt.Print("\nEnter your choice (1-5): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		switch input {
		case "1":
			listBlockedSitesUI()
		case "2":
			addBlockedSiteUI(reader)
		case "3":
			removeBlockedSiteUI(reader)
		case "4":
			reloadConfigUI()
		case "5":
			fmt.Println("\nExiting and cleaning up system proxy settings...")
			return
		default:
			fmt.Println("\nInvalid choice! Please enter a number between 1 and 5.")
		}

		fmt.Print("\nPress Enter to continue...")
		_, _ = reader.ReadString('\n')
	}
}

func printMenu() {
	// ANSI sequence to clear screen and move cursor to top-left
	fmt.Print("\033[H\033[2J")
	fmt.Println("==========================================================")
	fmt.Println("   🚫  FOCUS PROXY - INTERACTIVE TERMINAL CONTROL  🚫   ")
	fmt.Println("==========================================================")
	fmt.Printf("  Status      : System Proxy [ENABLED]\n")
	fmt.Printf("  Proxy Addr  : %s\n", proxyAddr)
	fmt.Printf("  Config File : %s (%d sites loaded)\n", hostsFilePath, len(listBlockedHosts()))
	fmt.Println("==========================================================")
	fmt.Println("  [1] List currently blocked sites")
	fmt.Println("  [2] Add a new website to block list")
	fmt.Println("  [3] Remove a website from block list (interactive)")
	fmt.Println("  [4] Reload configuration from JSON file")
	fmt.Println("  [5] Quit & Disable system proxy")
	fmt.Println("==========================================================")
}

func listBlockedSitesUI() {
	hosts := listBlockedHosts()
	fmt.Println("\n--- BLOCKED WEBSITES ---")
	if len(hosts) == 0 {
		fmt.Println("   No sites blocked currently.")
		return
	}
	for i, host := range hosts {
		fmt.Printf("   [%d] %s\n", i+1, host)
	}
}

func addBlockedSiteUI(reader *bufio.Reader) {
	fmt.Println("\n --- BLOCK A NEW WEBSITE ---")
	fmt.Print("Enter website domain (e.g., facebook.com): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input.")
		return
	}
	host := strings.TrimSpace(input)
	if host == "" {
		fmt.Println("Domain cannot be empty.")
		return
	}

	if addBlockedHost(host) {
		fmt.Printf("Successfully added and saved: %s\n", host)
	} else {
		fmt.Println("Failed to add. Domain might be invalid or already blocked.")
	}
}

func removeBlockedSiteUI(reader *bufio.Reader) {
	hosts := listBlockedHosts()
	fmt.Println("\n --- UNBLOCK A WEBSITE ---")
	if len(hosts) == 0 {
		fmt.Println("   No blocked sites to remove.")
		return
	}

	for i, host := range hosts {
		fmt.Printf("   [%d] %s\n", i+1, host)
	}

	fmt.Print("\nEnter the number of the site to unblock (or press Enter to cancel): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input.")
		return
	}
	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("Action cancelled.")
		return
	}

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(hosts) {
		fmt.Println("Invalid selection.")
		return
	}

	targetHost := hosts[idx-1]
	if removeBlockedHost(targetHost) {
		fmt.Printf("Successfully unblocked and saved: %s\n", targetHost)
	} else {
		fmt.Println("Failed to remove website.")
	}
}

func reloadConfigUI() {
	fmt.Println("\nReloading configuration from JSON...")
	if err := initHostsFile(); err != nil {
		fmt.Printf("Error reloading config: %v\n", err)
	} else {
		fmt.Printf("Config reloaded successfully. %d sites loaded.\n", len(listBlockedHosts()))
	}
}
