#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASS=0
FAIL=0
BINARY="calc"

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
echo -e "${YELLOW}=== Challenge 05: Functions (Calculator) ===${NC}"
echo ""

# --- Structural checks ---
echo "Checking project structure..."
if [ -f "go.mod" ]; then
    pass "go.mod exists"
else
    fail "go.mod not found" "run: go mod init calc" ""
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

# Helper function
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

# --- Arithmetic operations ---
run_test "10 add 5 = 15.00"              "Result: 15.00"              10 add 5
run_test "1.5 add 2.3 = 3.80"            "Result: 3.80"               1.5 add 2.3
run_test "10 sub 3 = 7.00"               "Result: 7.00"               10 sub 3
run_test "3 sub 10 = -7.00"              "Result: -7.00"              3 sub 10
run_test "4 mul 3 = 12.00"               "Result: 12.00"              4 mul 3
run_test "4 mul 2.5 = 10.00"             "Result: 10.00"              4 mul 2.5
run_test "7 div 2 = 3.50"                "Result: 3.50"               7 div 2
run_test "10 div 5 = 2.00"               "Result: 2.00"               10 div 5

# --- Error cases ---
run_test "division by zero"               "error: division by zero"    10 div 0
run_test "unknown operator 'pow'"         "error: unknown operator 'pow'" 10 pow 2
run_test "invalid number 'abc'"           "error: invalid number 'abc'" abc add 5
run_test "invalid number 'xyz'"           "error: invalid number 'xyz'" 5 add xyz

# --- Usage cases ---
# No arguments
ACTUAL=$(./"${BINARY}" 2>&1 || true)
EXPECTED="usage: calc <num1> <operator> <num2>"
if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "no args → usage"
else
    fail "no args → usage" "$EXPECTED" "$ACTUAL"
fi

run_test "too few args → usage"           "usage: calc <num1> <operator> <num2>" 1 2
run_test "too many args → usage"          "usage: calc <num1> <operator> <num2>" 1 add 2 3

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
