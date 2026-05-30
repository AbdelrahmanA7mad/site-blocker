# Site Blocker

`Site Blocker` is a Windows desktop tool that blocks selected websites by routing traffic through a local proxy and applying that proxy at the system level.

The app runs a local proxy on `127.0.0.1:8080`, then enables or disables the Windows proxy settings so traffic goes through that local filter.

## What It Does

This project is meant to help reduce distraction while working or studying. Instead of editing the `hosts` file manually, the app:

- manages a blocked website list from a JSON file
- turns the system proxy on or off
- blocks requests to matching websites
- provides a terminal UI for quick management

## Requirements

- Windows
- Go 1.26.1 or a compatible newer version
- Permission to update system proxy settings

## Running the App

### Development

```bash
go run .
```

### Build an Executable

```bash
go build -o sitesblocker.exe
```

## How It Works

1. The app starts in interactive terminal mode.
2. It creates or loads `blocked_hosts.json`.
3. It starts a local HTTP proxy on port `8080`.
4. When enabled, Windows proxy settings are updated to point to the local proxy.
5. If a requested host matches the blocked list, the app returns `403 Forbidden`.
6. On exit, the proxy is disabled and the system settings are restored.

## Terminal Menu

The interactive menu lets you:

- start the proxy
- stop the proxy
- list blocked websites
- add a website to the block list
- remove a website from the block list
- reload configuration from the JSON file
- run the app in the background
- reset Windows proxy settings
- quit and clean up

## Configuration File

Blocked websites are stored in:

```text
blocked_hosts.json
```

### Default List

If the file does not exist, the app creates it automatically with these defaults:

- `facebook.com`
- `instagram.com`
- `tiktok.com`
- `x.com`
- `twitter.com`
- `whatsapp.com`

### File Format

The file is a simple JSON array:

```json
[
  "facebook.com",
  "instagram.com",
  "tiktok.com"
]
```

## Main Files

- `main.go`: application entry point, interactive mode, and background mode
- `config.go`: load, save, and manage the blocked website list
- `proxy_handler.go`: request interception and allow/block logic
- `system_proxy_windows.go`: Windows Registry proxy control
- `ui.go`: terminal menu and user interaction

## Important Notes

- This project is Windows-only because it depends on the Windows Registry.
- The local proxy listens on `127.0.0.1:8080`.
- Blocking is based on hostnames only, not page content inspection.
- When the app exits, it disables the system proxy automatically.
