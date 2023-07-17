clear;

function get() {
  RES=$(curl -s http://localhost:8083/api/foo)
  OUT=$(echo "$RES" | jq -r '.message')
  if [[ -z "$OUT" || "null" == "$OUT" ]]; then
    OUT=$(echo "$RES" | jq -r '.error')
  fi
  echo "$OUT"
}

while true; do

  get &
  sleep 0.1

done