#!/system/bin/sh

# Wait for boot to complete
while [ "$(getprop sys.boot_completed)" != "1" ]; do
    sleep 1
done

MODDIR=${0%/*}
CONF_FILE=$MODDIR/config.conf
LOG_FILE=$MODDIR/logs/service.log
ERROR_LOG=$MODDIR/logs/error.log

# Ensure log directory exists
mkdir -p $MODDIR/logs

# Function to log messages
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" >> "$LOG_FILE"
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Function to log errors
log_error() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - ERROR: $1" >> "$ERROR_LOG"
    echo "$(date '+%Y-%m-%d %H:%M:%S') - ERROR: $1" >> "$LOG_FILE"
    echo "$(date '+%Y-%m-%d %H:%M:%S') - ERROR: $1"
}

# Load configuration
if [ -f "$CONF_FILE" ]; then
    . "$CONF_FILE"
    log "Configuration loaded successfully"
else
    log_error "Configuration file not found at $CONF_FILE"
    exit 1
fi

# Validate configuration
if [ -z "$SERVER_URL" ]; then
    log_error "SERVER_URL is not set in config file"
    exit 1
fi

if [ -z "$API_KEY" ]; then
    log_error "API_KEY is not set in config file"
    exit 1
fi

# Function to send heartbeat
send_heartbeat() {
    local response
    local http_code
    
    # Use curl with timeout and retry
    response=$(curl -L -s -w "\n%{http_code}" \
        --connect-timeout 5 \
        --max-time 10 \
        --retry 3 \
        --retry-delay 2 \
        --retry-max-time 30 \
        "$SERVER_URL/heartbeat?key=$API_KEY")
    
    http_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "200" ]; then
        log "Heartbeat sent successfully"
    else
        log_error "Failed to send heartbeat: HTTP $http_code - $response_body"
    fi
}

# Main loop
log "Service started with PID $$"
log "Server URL: $SERVER_URL"
log "Heartbeat interval: 30 seconds"

while true; do
    # Check network connectivity
    if ping -c 1 -W 1 8.8.8.8 >/dev/null 2>&1; then
        send_heartbeat
    else
        log "No network connectivity, skipping heartbeat"
    fi
    
    sleep 30
done
