#!/usr/bin/env sh

set +e
set -x

# --- Configuration ---
API_URL="https://kabbalahmedia.info/feed_api/feed"
API_DATA='{
    "more_items": 20,
    "current_feed": [],
    "namespace": "kmedia-app-1.84",
    "options": {
        "with_posts": false,
        "languages": ["ru", "en", "he"]
    }
}'
TODAY_DATE="$(date '+%Y-%m-%d')"
TIMESTAMP="$(date '+%Y%m%d%H%M%S')"
LOG_FILE="/tmp/feed_sanity_${TIMESTAMP}.log"
OUTPUT_FILE="/tmp/feed_output_${TIMESTAMP}.json"
EMAIL_SUBJECT="ERROR: Feed API Sanity Check Failed, NO MORNING LESSON!"
EMAIL_RECIPIENTS="edoshor@gmail.com kolmanv@gmail.com yurihechter@gmail.com"

# --- Functions ---
cleanup() {
  # Remove temporary files older than 1 day
  find /tmp -name "feed_sanity_*" -type f -mtime +1 -exec rm -f {} \;
  find /tmp -name "feed_output_*" -type f -mtime +1 -exec rm -f {} \;
}

# --- Main Script ---

# Initial log message
echo "Starting Feed API Sanity Check at $(date)" > "${LOG_FILE}"
echo "Checking for DAILY_LESSON with date: ${TODAY_DATE}" >> "${LOG_FILE}"

# Make the API call, save output, log errors
curl --location "${API_URL}" \
     --header 'Accept: application/json' \
     --header 'Content-Type: application/json' \
     --data "${API_DATA}" \
     --silent --show-error \
     -o "${OUTPUT_FILE}" >> "${LOG_FILE}" 2>&1

CURL_EXIT_CODE=$?

# Check if curl command itself failed
if [ "${CURL_EXIT_CODE}" -ne 0 ]; then
  echo "ERROR: curl command failed with exit code ${CURL_EXIT_CODE}. Check log file." >> "${LOG_FILE}"
  # Send failure email
  (echo "Curl command failed. See attached logs."; uuencode "${LOG_FILE}" feed_sanity.log; uuencode "${OUTPUT_FILE}" feed_output.json) | mail -s "${EMAIL_SUBJECT}" "${EMAIL_RECIPIENTS}"
  cleanup
  exit 1
fi

# Check if the output file contains the desired entry
# Grep for the content type, then pipe to grep for today's date within that context
# Use -c to count matches
MATCH_COUNT=$(jq --arg today "${TODAY_DATE}" \
  '.feed | map(select(.content_type == "DAILY_LESSON" and (.created_at | type == "string" and startswith($today)))) | length' \
  "${OUTPUT_FILE}")

# Evaluate the result
if [ "${MATCH_COUNT}" -gt 0 ]; then
  # Success Condition
  echo "Sanity OK: Found ${MATCH_COUNT} DAILY_LESSON entry/entries for ${TODAY_DATE}." | tee -a "${LOG_FILE}"
  cleanup
  exit 0
else
  # Failure Condition
  echo "ERROR: Sanity Check Failed. No DAILY_LESSON entry found for ${TODAY_DATE}." | tee -a "${LOG_FILE}"
  # Send failure email
  (echo "Feed API check failed: DAILY_LESSON for ${TODAY_DATE} not found. See attached logs."; uuencode "${LOG_FILE}" feed_sanity.log; uuencode "${OUTPUT_FILE}" feed_output.json) | mail -s "${EMAIL_SUBJECT}" "${EMAIL_RECIPIENTS}"
  cleanup
  exit 1
fi
