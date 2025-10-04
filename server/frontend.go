package server

import (
	"fmt"
	"net/http"
)

// HandleFrontend serves the HTML frontend
func HandleFrontend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><title>Touchpad Visualizer</title></head>
<body>
<h1>Multitouch WebSocket Server</h1>
<p>WebSocket endpoint: <code>ws://localhost:8080/ws</code></p>
<p>Status: <span id="status">Connecting...</span></p>
<pre id="output"></pre>
<script>
const ws = new WebSocket('ws://localhost:8080/ws');
const status = document.getElementById('status');
const output = document.getElementById('output');

ws.onopen = () => {
    status.textContent = 'Connected';
    status.style.color = 'green';
};

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    output.textContent = JSON.stringify(data, null, 2);
};

ws.onerror = (error) => {
    status.textContent = 'Error';
    status.style.color = 'red';
};

ws.onclose = () => {
    status.textContent = 'Disconnected';
    status.style.color = 'red';
};
</script>
</body>
</html>`)
}
