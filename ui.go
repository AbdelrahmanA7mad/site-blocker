package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func routeRequest(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		http.Redirect(w, r, "/ui", http.StatusTemporaryRedirect)
		return
	case r.URL.Path == "/ui":
		serveUI(w, r)
		return
	case r.URL.Path == "/api/blocked":
		handleBlockedAPI(w, r)
		return
	default:
		handleProxyRequest(w, r)
	}
}

func serveUI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(uiHTML))
}

func handleBlockedAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, map[string]any{"blocked": listBlockedHosts()})
	case http.MethodPost:
		var req struct {
			Host string `json:"host"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		if ok := addBlockedHost(req.Host); !ok {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid or duplicate host"})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"blocked": listBlockedHosts()})
	case http.MethodDelete:
		host := cleanHost(r.URL.Query().Get("host"))
		if host == "" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "missing host"})
			return
		}
		if ok := removeBlockedHost(host); !ok {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "host not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"blocked": listBlockedHosts()})
	default:
		w.Header().Set("Allow", strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodDelete}, ", "))
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

const uiHTML = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>Focus Proxy</title>
<style>
  :root { --bg:#f7f7f5; --card:#ffffff; --line:#ddd; --text:#181818; --muted:#666; --ok:#0a7d3b; --danger:#b42318; }
  * { box-sizing: border-box; }
  body { margin:0; font-family: "Segoe UI", Tahoma, sans-serif; background:var(--bg); color:var(--text); }
  .wrap { max-width:680px; margin:40px auto; padding:0 16px; }
  .card { background:var(--card); border:1px solid var(--line); border-radius:12px; padding:16px; }
  h1 { margin:0 0 10px; font-size:20px; }
  p { margin:0 0 16px; color:var(--muted); font-size:14px; }
  form { display:flex; gap:8px; margin-bottom:12px; }
  input { flex:1; border:1px solid var(--line); border-radius:8px; padding:10px; font-size:14px; }
  button { border:1px solid var(--line); background:#fff; border-radius:8px; padding:10px 12px; cursor:pointer; }
  ul { list-style:none; padding:0; margin:0; border:1px solid var(--line); border-radius:8px; overflow:hidden; }
  li { display:flex; justify-content:space-between; align-items:center; gap:8px; padding:10px 12px; border-bottom:1px solid var(--line); }
  li:last-child { border-bottom:none; }
  .host { font-family: Consolas, monospace; font-size:13px; }
  .remove { color:var(--danger); }
  .status { min-height:20px; margin-top:10px; font-size:13px; color:var(--ok); }
</style>
</head>
<body>
  <div class="wrap">
    <div class="card">
      <h1>Focus Proxy Control</h1>
      <p>Add or remove blocked websites instantly.</p>
      <form id="add-form">
        <input id="host-input" placeholder="example.com" />
        <button type="submit">Block</button>
      </form>
      <ul id="hosts"></ul>
      <div id="status" class="status"></div>
    </div>
  </div>
<script>
const hostsEl = document.getElementById('hosts');
const statusEl = document.getElementById('status');
const formEl = document.getElementById('add-form');
const inputEl = document.getElementById('host-input');

function setStatus(msg, isError=false) {
  statusEl.textContent = msg;
  statusEl.style.color = isError ? '#b42318' : '#0a7d3b';
}

function renderHosts(hosts) {
  hostsEl.innerHTML = '';
  if (!hosts.length) {
    const li = document.createElement('li');
    li.textContent = 'No blocked hosts.';
    hostsEl.appendChild(li);
    return;
  }
  for (const host of hosts) {
    const li = document.createElement('li');
    li.innerHTML = '<span class="host"></span><button class="remove">Unblock</button>';
    li.querySelector('.host').textContent = host;
    li.querySelector('button').onclick = () => removeHost(host);
    hostsEl.appendChild(li);
  }
}

async function loadHosts() {
  const res = await fetch('/api/blocked');
  const data = await res.json();
  renderHosts(data.blocked || []);
}

async function addHost(host) {
  const res = await fetch('/api/blocked', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ host })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Failed');
  renderHosts(data.blocked || []);
  setStatus('Blocked: ' + host);
}

async function removeHost(host) {
  const res = await fetch('/api/blocked?host=' + encodeURIComponent(host), { method: 'DELETE' });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Failed');
  renderHosts(data.blocked || []);
  setStatus('Unblocked: ' + host);
}

formEl.addEventListener('submit', async (e) => {
  e.preventDefault();
  const host = inputEl.value.trim();
  if (!host) return;
  try {
    await addHost(host);
    inputEl.value = '';
    inputEl.focus();
  } catch (err) {
    setStatus(err.message, true);
  }
});

loadHosts().catch(() => setStatus('Failed to load hosts', true));
</script>
</body>
</html>`
