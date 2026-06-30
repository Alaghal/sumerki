#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
SMOKE_EMAIL="${SMOKE_EMAIL:-northern@example.com}"
SMOKE_PASSWORD="${SMOKE_PASSWORD:-password123}"

if ! command -v jq >/dev/null 2>&1; then
  echo "ERROR jq is required for scripts/smoke-api.sh" >&2
  exit 1
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

TOKEN=""

request() {
  local method="$1"
  local path="$2"
  local body="${3:-}"
  local token="${4:-}"
  local output="$TMP_DIR/response.json"
  local status

  if [[ -n "$body" && -n "$token" ]]; then
    status="$(curl -sS -o "$output" -w "%{http_code}" -X "$method" "$API_BASE_URL$path" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $token" \
      -d "$body")"
  elif [[ -n "$body" ]]; then
    status="$(curl -sS -o "$output" -w "%{http_code}" -X "$method" "$API_BASE_URL$path" \
      -H "Content-Type: application/json" \
      -d "$body")"
  elif [[ -n "$token" ]]; then
    status="$(curl -sS -o "$output" -w "%{http_code}" -X "$method" "$API_BASE_URL$path" \
      -H "Authorization: Bearer $token")"
  else
    status="$(curl -sS -o "$output" -w "%{http_code}" -X "$method" "$API_BASE_URL$path")"
  fi

  RESPONSE_BODY="$(cat "$output")"
  RESPONSE_STATUS="$status"
}

require_success() {
  local label="$1"
  if [[ "$RESPONSE_STATUS" -lt 200 || "$RESPONSE_STATUS" -ge 300 ]]; then
    echo "FAIL $label HTTP $RESPONSE_STATUS"
    echo "$RESPONSE_BODY"
    exit 1
  fi
}

optional_post() {
  local label="$1"
  local path="$2"
  local body="${3:-}"

  request POST "$path" "$body" "$TOKEN"
  if [[ "$RESPONSE_STATUS" -ge 200 && "$RESPONSE_STATUS" -lt 300 ]]; then
    echo "OK $label"
  else
    local code
    code="$(jq -r '.error.code // "unknown"' <<<"$RESPONSE_BODY")"
    echo "NOTE $label skipped: $code"
  fi
}

auth() {
  local credentials
  credentials="$(jq -nc --arg email "$SMOKE_EMAIL" --arg password "$SMOKE_PASSWORD" '{email:$email,password:$password}')"

  request POST "/api/auth/login" "$credentials"
  if [[ "$RESPONSE_STATUS" -ge 200 && "$RESPONSE_STATUS" -lt 300 ]]; then
    TOKEN="$(jq -r '.token' <<<"$RESPONSE_BODY")"
    echo "OK auth login"
    return
  fi

  request POST "/api/auth/register" "$credentials"
  require_success "auth register"
  TOKEN="$(jq -r '.token' <<<"$RESPONSE_BODY")"
  echo "OK auth register"
}

ensure_kingdom() {
  request GET "/api/kingdoms/me" "" "$TOKEN"
  require_success "kingdom fetch"
  if [[ "$(jq -r '.kingdom == null' <<<"$RESPONSE_BODY")" == "true" ]]; then
    request POST "/api/kingdoms" '{"name":"Smoke Hold","culture":"northern_principality"}' "$TOKEN"
    require_success "kingdom create"
    echo "OK kingdom created"
  else
    echo "OK kingdom"
  fi
}

first_event_choice() {
  request GET "/api/events/me" "" "$TOKEN"
  require_success "events"
  local event_id choice_key
  event_id="$(jq -r '.events[0].id // empty' <<<"$RESPONSE_BODY")"
  choice_key="$(jq -r '.events[0].choices[0].key // empty' <<<"$RESPONSE_BODY")"
  if [[ -n "$event_id" && -n "$choice_key" ]]; then
    local body
    body="$(jq -nc --arg choiceKey "$choice_key" '{choiceKey:$choiceKey}')"
    optional_post "event choice" "/api/events/$event_id/choose" "$body"
  else
    echo "NOTE event choice skipped: no active event"
  fi
  echo "OK events"
}

first_raid_target() {
  request GET "/api/neighbors" "" "$TOKEN"
  require_success "neighbors"
  local target
  target="$(jq -r '.neighbors[]? | select(.canRaid == true) | .kingdomId' <<<"$RESPONSE_BODY" | head -n 1)"
  if [[ -n "$target" ]]; then
    local body
    body="$(jq -nc --arg defenderKingdomId "$target" '{defenderKingdomId:$defenderKingdomId,units:[{unitType:"militia",amount:3}]}')"
    optional_post "raid start" "/api/raids/start" "$body"
  else
    echo "WARN no raid target available"
  fi

  request GET "/api/raids/me" "" "$TOKEN"
  require_success "raids"
  echo "OK raids"
}

auth
ensure_kingdom

request GET "/api/ruler/me" "" "$TOKEN"
require_success "ruler"
echo "OK ruler"

request GET "/api/resources/me" "" "$TOKEN"
require_success "resources"
echo "OK resources"

request GET "/api/buildings/me" "" "$TOKEN"
require_success "buildings"
optional_post "building upgrade" "/api/buildings/farm/upgrade"
echo "OK buildings"

request GET "/api/army/me" "" "$TOKEN"
require_success "army"
optional_post "unit training" "/api/army/train" '{"unitType":"militia","amount":3}'
echo "OK army"

request GET "/api/missions/available" "" "$TOKEN"
require_success "available missions"
optional_post "mission start" "/api/missions/start" '{"missionKey":"black_forest_expedition","units":[{"unitType":"militia","amount":7}]}'

request GET "/api/missions/me" "" "$TOKEN"
require_success "missions"
echo "OK missions"

request GET "/api/reports/me" "" "$TOKEN"
require_success "reports"
echo "OK reports"

request GET "/api/patron/options" "" "$TOKEN"
require_success "patron options"
request GET "/api/patron/me" "" "$TOKEN"
require_success "patron status"
if [[ "$(jq -r '.patron == null' <<<"$RESPONSE_BODY")" == "true" ]]; then
  optional_post "patron join" "/api/patron/join" '{"patron":"old_pact"}'
fi
request GET "/api/patron/pressure" "" "$TOKEN"
require_success "patron pressure"
echo "OK patron"

first_event_choice

request GET "/api/reports/me" "" "$TOKEN"
require_success "reports after event"
echo "OK reports after event"

first_raid_target

request GET "/api/reports/me" "" "$TOKEN"
require_success "reports after raids"
echo "OK smoke-api complete"
