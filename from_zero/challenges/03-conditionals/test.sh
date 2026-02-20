#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASS=0
FAIL=0
BINARY="classify"

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
echo -e "${YELLOW}=== Challenge 03: Conditionals ===${NC}"
echo ""

# --- Structural checks ---
echo "Checking project structure..."
if [ -f "go.mod" ]; then
    pass "go.mod exists"
else
    fail "go.mod not found" "run: go mod init classify" ""
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

# Helper function for test cases
run_test() {
    local description="$1"
    local expected="$2"
    shift 2
    local actual
    actual=$(./"${BINARY}" "$@" 2>&1 || true)
    if [ "$actual" = "$expected" ]; then
        pass "$description"
    else
        fail "$description" "$expected" "$actual"
    fi
}

# --- Test cases ---
run_test "classify 5 → positive"          "positive"                5
run_test "classify 100 → positive"        "positive"                100
run_test "classify -3 → negative"         "negative"                -- -3
run_test "classify -999 → negative"       "negative"                -- -999
run_test "classify 0 → zero"             "zero"                    0
run_test "classify abc → error"          "error: not a number"     abc
run_test "classify 3.14 → error"         "error: not a number"     3.14

# No args case (special — don't pass any args)
ACTUAL=$(./"${BINARY}" 2>&1 || true)
EXPECTED="usage: classify <number>"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "no args → usage message"
else
    fail "no args → usage message" "$EXPECTED" "$ACTUAL"
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
