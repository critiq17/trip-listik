#!/usr/bin/env bash
set -euo pipefail

# ── dev-local.sh ─────────────────────────────────────────────────────────────
# Local development setup for TripListik.
# - Starts an ngrok tunnel for the backend (port 8080)
# - Registers the Telegram webhook automatically (if BOT_TOKEN is set)
# - Prints instructions to start backend and frontend
# ─────────────────────────────────────────────────────────────────────────────

BACKEND_PORT=8080
NGROK_API="http://localhost:4040/api/tunnels"

# ── 1. Check dependencies ─────────────────────────────────────────────────────

if ! command -v ngrok &>/dev/null; then
  echo "ERROR: ngrok is not installed."
  echo ""
  echo "Install options:"
  echo "  Homebrew (macOS/Linux):  brew install ngrok/ngrok/ngrok"
  echo "  Direct download:         https://ngrok.com/download"
  echo "  Snap (Linux):            sudo snap install ngrok"
  echo ""
  echo "After installing, authenticate once: ngrok config add-authtoken <YOUR_TOKEN>"
  exit 1
fi

if ! command -v curl &>/dev/null; then
  echo "ERROR: curl is required but not found. Please install curl."
  exit 1
fi

# ── 2. Check BOT_TOKEN ────────────────────────────────────────────────────────

if [[ -z "${BOT_TOKEN:-}" ]]; then
  echo "WARN: BOT_TOKEN env var is not set. Webhook registration will be skipped."
  echo "      Export it before running: export BOT_TOKEN=your_bot_token"
  echo ""
fi

# ── 3. Kill any stale ngrok process ──────────────────────────────────────────

if curl -sf "$NGROK_API" &>/dev/null; then
  echo "INFO: Existing ngrok session detected — reusing it."
else
  echo "INFO: Starting ngrok tunnel on port $BACKEND_PORT..."
  ngrok http "$BACKEND_PORT" --log=stdout &>/tmp/ngrok-triplistik.log &
  NGROK_PID=$!
  echo "INFO: ngrok PID $NGROK_PID — logs at /tmp/ngrok-triplistik.log"
fi

# ── 4. Wait for ngrok to be ready (up to 10 seconds) ─────────────────────────

echo "INFO: Waiting for ngrok to start..."
NGROK_URL=""
for i in $(seq 1 20); do
  sleep 0.5
  if curl -sf "$NGROK_API" &>/dev/null; then
    # Extract HTTPS URL — prefer jq, fall back to python3
    if command -v jq &>/dev/null; then
      NGROK_URL=$(curl -sf "$NGROK_API" | jq -r '.tunnels[] | select(.proto=="https") | .public_url' | head -1)
    else
      NGROK_URL=$(curl -sf "$NGROK_API" | python3 -c "
import json, sys
tunnels = json.load(sys.stdin).get('tunnels', [])
https = [t['public_url'] for t in tunnels if t.get('proto') == 'https']
print(https[0] if https else '')
")
    fi
    if [[ -n "$NGROK_URL" ]]; then
      break
    fi
  fi
done

if [[ -z "$NGROK_URL" ]]; then
  echo "ERROR: ngrok did not start within 10 seconds."
  echo "       Check logs: /tmp/ngrok-triplistik.log"
  echo "       Or visit http://localhost:4040 in your browser."
  exit 1
fi

echo ""
echo "  ngrok tunnel: $NGROK_URL"
echo "  Backend URL:  $NGROK_URL/v1"
echo ""

# ── 5. Register Telegram webhook ──────────────────────────────────────────────

if [[ -n "${BOT_TOKEN:-}" ]]; then
  WEBHOOK_URL="${NGROK_URL}/v1/webhook/telegram"
  echo "INFO: Registering Telegram webhook..."
  echo "  Webhook URL: $WEBHOOK_URL"
  RESULT=$(curl -s -X POST "https://api.telegram.org/bot${BOT_TOKEN}/setWebhook" \
    -H "Content-Type: application/json" \
    -d "{\"url\": \"${WEBHOOK_URL}\"}")
  echo "  Telegram API response: $RESULT"
  echo ""
fi

# ── 6. Print startup instructions ────────────────────────────────────────────

echo "────────────────────────────────────────────────────────"
echo "  TripListik Local Dev Setup"
echo "────────────────────────────────────────────────────────"
echo ""
echo "  Set these env vars in your backend .env:"
echo "    MINI_APP_URL=http://localhost:5173"
echo "    BOT_TOKEN=${BOT_TOKEN:-<your_bot_token>}"
echo ""
echo "  To start the backend (in a new terminal):"
echo "    cd backend && go run ./cmd/api"
echo ""
echo "  To start the frontend (in a new terminal):"
echo "    cd frontend && npm run dev"
echo ""
echo "  ngrok dashboard: http://localhost:4040"
echo "────────────────────────────────────────────────────────"
