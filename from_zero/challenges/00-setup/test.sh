#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

PASS=0
FAIL=0
BINARY="setup"

pass() {
    echo -e "  ${GREEN}PASS${NC}: $1"
    PASS=$((PASS + 1))
}

fail() {
    echo -e "  ${RED}FAIL${NC}: $1"
    echo -e "        $2"
    FAIL=$((FAIL + 1))
}

echo ""
echo -e "${YELLOW}=== Challenge 00: Setup ===${NC}"
echo ""

# --- Test 1: Go is installed ---
echo "Checking Go installation..."
if go version > /dev/null 2>&1; then
    pass "Go is installed ($(go version | awk '{print $3}'))"
else
    fail "Go is not installed" "Install Go from https://go.dev/dl/"
fi

# --- Test 2: go.mod exists ---
echo "Checking project structure..."
if [ -f "go.mod" ]; then
    pass "go.mod exists"
else
    fail "go.mod not found" "Run: go mod init <module-name>"
fi

# --- Test 3: main.go exists ---
if [ -f "main.go" ]; then
    pass "main.go exists"
else
    fail "main.go not found" "Create a main.go file with package main and func main()"
fi

# --- Test 4: Build succeeds ---
echo "Building..."
if go build -o "${BINARY}" . 2>/dev/null; then
    pass "Build succeeded"
else
    fail "Build failed" "Check your code for syntax errors: go build -o ${BINARY} ."
fi

# --- Test 5: Output is correct ---
if [ -f "${BINARY}" ]; then
    echo "Running..."
    ACTUAL=$(./"${BINARY}" 2>&1)
    EXPECTED="ready"

    if [ "$ACTUAL" = "$EXPECTED" ]; then
        pass "Output is 'ready'"
    else
        fail "Unexpected output" "Expected: '${EXPECTED}', Got: '${ACTUAL}'"
    fi
else
    fail "Binary not found" "Build must succeed before running"
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
