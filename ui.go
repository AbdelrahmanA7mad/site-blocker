package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func runTerminalUI() {
	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		fmt.Print("\nEnter your choice (1-9): ")
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
			enableSystemProxy()
			fmt.Println("\n\033[32;1m>>> Focus Proxy has been STARTED! <<<\033[0m")
		case "2":
			disableSystemProxy()
			fmt.Println("\n\033[31;1m>>> Focus Proxy has been STOPPED! <<<\033[0m")
		case "3":
			listBlockedSitesUI()
		case "4":
			addBlockedSiteUI(reader)
		case "5":
			removeBlockedSiteUI(reader)
		case "6":
			reloadConfigUI()
		case "7":
			fmt.Println("\n\033[33;1m>>> Closing this window and running in background... <<<\033[0m")
			fmt.Println("Focus Proxy will continue running silently in the background.")
			fmt.Println("To restore this window at any time, simply run sitesblocker.exe again!")
			time.Sleep(2500 * time.Millisecond)

			err := spawnBackgroundClone(isProxyActive())
			if err != nil {
				fmt.Printf("\nError running in background: %v\n", err)
				fmt.Print("Press Enter to continue...")
				_, _ = reader.ReadString('\n')
				continue
			}
			os.Exit(0)
		case "8":
			resetSystemProxy()
			fmt.Println("\n\033[33;1m>>> Windows System Proxy has been fully RESET to defaults! <<<\033[0m")
		case "9":
			fmt.Println("\nExiting and cleaning up system proxy settings...")
			return
		default:
			fmt.Println("\nInvalid choice! Please enter a number between 1 and 9.")
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
	
	statusStr := "\033[31;1m[DISABLED]\033[0m"
	if isProxyActive() {
		statusStr = "\033[32;1m[ENABLED]\033[0m"
	}
	fmt.Printf("  Status      : System Proxy %s\n", statusStr)
	fmt.Printf("  Proxy Addr  : %s\n", proxyAddr)
	fmt.Printf("  Config File : %s (%d sites loaded)\n", hostsFilePath, len(listBlockedHosts()))
	fmt.Println("==========================================================")
	fmt.Println("  [1] Start Focus Proxy")
	fmt.Println("  [2] Stop Focus Proxy")
	fmt.Println("  [3] List currently blocked sites")
	fmt.Println("  [4] Add a new website to block list")
	fmt.Println("  [5] Remove a website from block list (interactive)")
	fmt.Println("  [6] Reload configuration from JSON file")
	fmt.Println("  [7] Hide terminal and run in background")
	fmt.Println("  [8] Reset system proxy settings to Windows defaults")
	fmt.Println("  [9] Quit & Disable system proxy")
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
