#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASS=0
FAIL=0
BINARY="fizzbuzz"

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
echo -e "${YELLOW}=== Challenge 04: Loops (FizzBuzz) ===${NC}"
echo ""

# --- Structural checks ---
echo "Checking project structure..."
if [ -f "go.mod" ]; then
    pass "go.mod exists"
else
    fail "go.mod not found" "run: go mod init fizzbuzz" ""
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

# --- Test: N=15 full output ---
ACTUAL=$(./"${BINARY}" 15 2>&1)
EXPECTED=$'1\n2\nFizz\n4\nBuzz\nFizz\n7\n8\nFizz\nBuzz\n11\nFizz\n13\n14\nFizzBuzz'

if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "N=15: full output matches"
else
    fail "N=15: output mismatch" "$EXPECTED" "$ACTUAL"
fi

# --- Test: N=1 ---
ACTUAL=$(./"${BINARY}" 1 2>&1)
if [ "$ACTUAL" = "1" ]; then
    pass "N=1: prints '1'"
else
    fail "N=1" "1" "$ACTUAL"
fi

# --- Test: N=3 ---
ACTUAL=$(./"${BINARY}" 3 2>&1)
EXPECTED=$'1\n2\nFizz'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "N=3: prints '1 2 Fizz'"
else
    fail "N=3" "$EXPECTED" "$ACTUAL"
fi

# --- Test: N=5 ---
ACTUAL=$(./"${BINARY}" 5 2>&1)
EXPECTED=$'1\n2\nFizz\n4\nBuzz'
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "N=5: ends with 'Buzz'"
else
    fail "N=5" "$EXPECTED" "$ACTUAL"
fi

# --- Test: N=0 (no output) ---
ACTUAL=$(./"${BINARY}" 0 2>&1)
if [ -z "$ACTUAL" ]; then
    pass "N=0: no output"
else
    fail "N=0" "(empty)" "$ACTUAL"
fi

# --- Test: Negative N (no output) ---
ACTUAL=$(./"${BINARY}" -- -5 2>&1)
if [ -z "$ACTUAL" ]; then
    pass "N=-5: no output"
else
    fail "N=-5" "(empty)" "$ACTUAL"
fi

# --- Test: No arguments ---
ACTUAL=$(./"${BINARY}" 2>&1 || true)
EXPECTED="usage: fizzbuzz <n>"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "No args: usage message"
else
    fail "No args" "$EXPECTED" "$ACTUAL"
fi

# --- Test: Not a number ---
ACTUAL=$(./"${BINARY}" abc 2>&1 || true)
EXPECTED="error: not a number"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "abc: error message"
else
    fail "abc" "$EXPECTED" "$ACTUAL"
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
