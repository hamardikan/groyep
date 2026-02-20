#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASS=0
FAIL=0
BINARY="wordfreq"

pass() {
    echo -e "  ${GREEN}PASS${NC}: $1"
    PASS=$((PASS + 1))
}

fail() {
    echo -e "  ${RED}FAIL${NC}: $1"
    echo -e "        Expected: '$2'"
    echo -e "        Got:      '$3'"
    FAIL=$((FAIL + 1))
}

echo ""
echo -e "${YELLOW}=== Challenge 06: Collections (Word Frequency) ===${NC}"
echo ""

# --- Structural checks ---
echo "Checking project structure..."
if [ -f "go.mod" ]; then
    pass "go.mod exists"
else
    fail "go.mod not found" "run: go mod init wordfreq" ""
fi

# --- Build ---
echo "Building..."
if go build -o "${BINARY}" . 2>/dev/null; then
    pass "Build succeeded"
else
    fail "Build failed" "compiles" "build errors"
    rm -f "${BINARY}"
    exit 1
fi

echo "Running tests..."

# --- Test: hello world hello go world hello ---
ACTUAL=$(./"${BINARY}" hello world hello go world hello 2>&1)
EXPECTED=$'go: 1\nhello: 3\nworld: 2'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "hello world hello go world hello"
else
    fail "three words with repeats" "$EXPECTED" "$ACTUAL"
fi

# --- Test: single word ---
ACTUAL=$(./"${BINARY}" apple 2>&1)
EXPECTED="apple: 1"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "single word 'apple'"
else
    fail "single word" "$EXPECTED" "$ACTUAL"
fi

# --- Test: a b a b a ---
ACTUAL=$(./"${BINARY}" a b a b a 2>&1)
EXPECTED=$'a: 3\nb: 2'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "a b a b a"
else
    fail "two alternating words" "$EXPECTED" "$ACTUAL"
fi

# --- Test: alphabetical sorting ---
ACTUAL=$(./"${BINARY}" cherry apple banana 2>&1)
EXPECTED=$'apple: 1\nbanana: 1\ncherry: 1'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "alphabetical sorting (cherry apple banana)"
else
    fail "alphabetical sorting" "$EXPECTED" "$ACTUAL"
fi

# --- Test: all same word ---
ACTUAL=$(./"${BINARY}" go go go 2>&1)
EXPECTED="go: 3"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "all same word (go go go)"
else
    fail "all same word" "$EXPECTED" "$ACTUAL"
fi

# --- Test: case sensitive ---
ACTUAL=$(./"${BINARY}" Hello hello HELLO 2>&1)
EXPECTED=$'HELLO: 1\nHello: 1\nhello: 1'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "case sensitive (Hello hello HELLO)"
else
    fail "case sensitive" "$EXPECTED" "$ACTUAL"
fi

# --- Test: no arguments ---
ACTUAL=$(./"${BINARY}" 2>&1 || true)
EXPECTED="usage: wordfreq <word1> [word2] ..."
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "no args → usage message"
else
    fail "no args → usage" "$EXPECTED" "$ACTUAL"
fi

# --- Cleanup ---
rm -f "${BINARY}"

# --- Summary ---
echo ""
TOTAL=$((PASS + FAIL))
if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}All ${TOTAL} tests passed!${NC}"
else
    echo -e "${RED}${FAIL} of ${TOTAL} tests failed.${NC}"
    exit 1
fi
