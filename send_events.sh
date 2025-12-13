#!/bin/zsh

# Traffic simulation for oldavista.com
# Game loop with continuous session churn

API_URL="http://localhost:8099/api/event"
DOMAIN="oldavista.com"

PAGES=(
    "/" "/about" "/contact" "/blog" "/products" "/home" "/index.html"
    "/pricing" "/features" "/docs" "/api" "/login" "/signup" "/faq"
    "/support" "/news" "/archive" "/search" "/help"
)

# Browser UAs
BROWSER_UAS=(
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15"
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0"
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"
    "Mozilla/5.0 (iPad; CPU OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"
    "Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/106.0.0.0"
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0"
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:120.0) Gecko/20100101 Firefox/120.0"
    "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0"
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15"
    "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1"
    "Mozilla/5.0 (iPad; CPU OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1"
    "Mozilla/5.0 (Linux; Android 12; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"
    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1"
    "Mozilla/5.0 (Linux; Android 13; SM-A536B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
    "Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"
    "Mozilla/5.0 (Linux; Android 10; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"
)

# Bot UAs (for testing bot filtering)
BOT_UAS=(
    "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
    "Mozilla/5.0 (compatible; Bingbot/2.0; +http://www.bing.com/bingbot.htm)"
    "Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)"
    "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)"
    "Twitterbot/1.0"
    "LinkedInBot/1.0 (compatible; Mozilla/5.0; +http://www.linkedin.com)"
    "Mozilla/5.0 (compatible; AhrefsBot/7.0; +http://ahrefs.com/robot/)"
    "Mozilla/5.0 (compatible; SemrushBot/7~bl; +http://www.semrush.com/bot.html)"
)

BROWSER_UA_COUNT=${#BROWSER_UAS[@]}
BOT_UA_COUNT=${#BOT_UAS[@]}
PAGE_COUNT=${#PAGES[@]}

# Session state - associative arrays keyed by session ID
typeset -A SESSION_IP
typeset -A SESSION_UA
typeset -A SESSION_SW
typeset -A SESSION_PAGES_LEFT
typeset -A SESSION_NEXT_TIME

# Active session IDs
typeset -a ACTIVE_IDS
ACTIVE_IDS=()

# Counters
SESSION_COUNTER=0
TOTAL_COMPLETED=0
TOTAL_PAGEVIEWS=0

# Config
MIN_ACTIVE=100
MAX_ACTIVE=300
BOT_PERCENT=10  # 10% of sessions are bots
TICK_MS=20

# Get time in ms
now_ms() {
    local s=$(date +%s)
    local ns=$(date +%N 2>/dev/null || echo "000000000")
    echo $(( s * 1000 + ${ns:0:3} ))
}

# Generate unique IP from counter
make_ip() {
    local n=$1
    echo "$(( (n % 200) + 20 )).$(( (n / 200) % 256 )).$(( (n / 51200) % 256 )).$(( (RANDOM % 254) + 1 ))"
}

# Random page count 1-50 weighted
rand_pages() {
    local r=$((RANDOM % 100))
    if   (( r < 35 )); then echo 1
    elif (( r < 55 )); then echo $((2 + RANDOM % 2))
    elif (( r < 75 )); then echo $((4 + RANDOM % 4))
    elif (( r < 90 )); then echo $((8 + RANDOM % 13))
    else                    echo $((21 + RANDOM % 30))
    fi
}

# Pick a user agent (mostly browsers, some bots)
pick_ua() {
    if (( RANDOM % 100 < BOT_PERCENT )); then
        echo "${BOT_UAS[$((RANDOM % BOT_UA_COUNT + 1))]}"
    else
        echo "${BROWSER_UAS[$((RANDOM % BROWSER_UA_COUNT + 1))]}"
    fi
}

# Create one new session
create_session() {
    ((SESSION_COUNTER++))
    local id=$SESSION_COUNTER
    local now=$1
    
    SESSION_IP[$id]=$(make_ip $id)
    SESSION_UA[$id]=$(pick_ua)
    SESSION_SW[$id]=$((320 + RANDOM % 1600))
    SESSION_PAGES_LEFT[$id]=$(rand_pages)
    SESSION_NEXT_TIME[$id]=$((now + RANDOM % 300))
    
    ACTIVE_IDS+=($id)
}

# Send pageview (backgrounded so it doesn't block)
send_pv() {
    local id=$1
    local page="${PAGES[$((RANDOM % PAGE_COUNT + 1))]}"
    
    curl -s -X POST "$API_URL" \
        -H "Content-Type: application/json" \
        -H "User-Agent: ${SESSION_UA[$id]}" \
        -H "X-Forwarded-For: ${SESSION_IP[$id]}" \
        -d "{\"name\":\"pageview\",\"domain\":\"$DOMAIN\",\"page\":\"$page\",\"screenWidth\":${SESSION_SW[$id]}}" \
        > /dev/null 2>&1 &
    
    ((TOTAL_PAGEVIEWS++))
}

# Remove session data
remove_session() {
    local id=$1
    unset "SESSION_IP[$id]"
    unset "SESSION_UA[$id]"
    unset "SESSION_SW[$id]"
    unset "SESSION_PAGES_LEFT[$id]"
    unset "SESSION_NEXT_TIME[$id]"
}

# Cleanup
cleanup() {
    echo ""
    echo "Stopping..."
    wait 2>/dev/null
    echo "Sessions completed: $TOTAL_COMPLETED"
    echo "Total pageviews: $TOTAL_PAGEVIEWS"
    exit 0
}
trap cleanup SIGINT SIGTERM

echo "═══════════════════════════════════════════════════════"
echo "  Traffic Simulation for $DOMAIN"
echo "═══════════════════════════════════════════════════════"
echo "  Target: $API_URL"
echo "  Active sessions target: $MIN_ACTIVE - $MAX_ACTIVE"
echo "  Browser UAs: $BROWSER_UA_COUNT | Bot UAs: $BOT_UA_COUNT"
echo "  Bot traffic: $BOT_PERCENT%"
echo "  Press Ctrl+C to stop"
echo "═══════════════════════════════════════════════════════"
echo ""

LAST_STATS=$(now_ms)
NEXT_SPAWN=$(now_ms)

# Main loop
while true; do
    NOW=$(now_ms)
    ACTIVE_COUNT=${#ACTIVE_IDS[@]}
    
    # --- SPAWN NEW SESSIONS ---
    if (( NOW >= NEXT_SPAWN )); then
        if (( ACTIVE_COUNT < MIN_ACTIVE )); then
            create_session $NOW
            NEXT_SPAWN=$((NOW + 10 + RANDOM % 40))
        elif (( ACTIVE_COUNT < MAX_ACTIVE )); then
            if (( RANDOM % 100 < 20 )); then
                create_session $NOW
            fi
            NEXT_SPAWN=$((NOW + 50 + RANDOM % 150))
        else
            if (( RANDOM % 100 < 5 )); then
                create_session $NOW
            fi
            NEXT_SPAWN=$((NOW + 200 + RANDOM % 300))
        fi
    fi
    
    # --- PROCESS ACTIVE SESSIONS ---
    typeset -a NEW_ACTIVE
    NEW_ACTIVE=()
    
    for id in "${ACTIVE_IDS[@]}"; do
        if (( NOW < SESSION_NEXT_TIME[$id] )); then
            NEW_ACTIVE+=($id)
            continue
        fi
        
        send_pv $id
        ((SESSION_PAGES_LEFT[$id]--))
        
        if (( SESSION_PAGES_LEFT[$id] <= 0 )); then
            ((TOTAL_COMPLETED++))
            remove_session $id
        else
            SESSION_NEXT_TIME[$id]=$((NOW + 200 + RANDOM % 1300))
            NEW_ACTIVE+=($id)
        fi
    done
    
    ACTIVE_IDS=("${NEW_ACTIVE[@]}")
    
    # --- STATS ---
    if (( NOW - LAST_STATS >= 1000 )); then
        echo "[$(date +%H:%M:%S)] Active: ${#ACTIVE_IDS[@]} | Completed: $TOTAL_COMPLETED | Pageviews: $TOTAL_PAGEVIEWS"
        LAST_STATS=$NOW
    fi
    
    sleep 0.0$((TICK_MS / 10))
done