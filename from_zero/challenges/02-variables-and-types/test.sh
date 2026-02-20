#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASS=0
FAIL=0
BINARY="variables"

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
echo -e "${YELLOW}=== Challenge 02: Variables and Types ===${NC}"
echo ""

# --- Test 1: go.mod exists ---
echo "Checking project structure..."
if [ -f "go.mod" ]; then
    pass "go.mod exists"
else
    fail "go.mod not found" "Run: go mod init variables"
fi

# --- Test 2: Build succeeds ---
echo "Building..."
if go build -o "${BINARY}" . 2>/dev/null; then
    pass "Build succeeded"
else
    fail "Build failed" "Fix syntax errors: go build -o ${BINARY} ."
fi

if [ ! -f "${BINARY}" ]; then
    echo -e "${RED}Cannot run tests — build failed.${NC}"
    rm -f "${BINARY}"
    exit 1
fi

echo "Running..."

# --- Test 3: Full output ---
ACTUAL=$(./"${BINARY}" 2>&1)
EXPECTED=$'Name: Gopher\nAge: 10\nHeight: 1.75\nIsAwesome: true'

if [ "$ACTUAL" = "$EXPECTED" ]; then
    pass "Full output matches"
else
    fail "Output mismatch" "Expected:\n${EXPECTED}\nGot:\n${ACTUAL}"
fi

# --- Test 4: Name line ---
LINE=$(echo "$ACTUAL" | sed -n '1p')
if [ "$LINE" = "Name: Gopher" ]; then
    pass "Line 1: Name: Gopher"
else
    fail "Line 1 wrong" "Expected: 'Name: Gopher', Got: '${LINE}'"
fi

# --- Test 5: Age line ---
LINE=$(echo "$ACTUAL" | sed -n '2p')
if [ "$LINE" = "Age: 10" ]; then
    pass "Line 2: Age: 10"
else
    fail "Line 2 wrong" "Expected: 'Age: 10', Got: '${LINE}'"
fi

# --- Test 6: Height line ---
LINE=$(echo "$ACTUAL" | sed -n '3p')
if [ "$LINE" = "Height: 1.75" ]; then
    pass "Line 3: Height: 1.75"
else
    fail "Line 3 wrong" "Expected: 'Height: 1.75', Got: '${LINE}' (hint: use %.2f)"
fi

# --- Test 7: IsAwesome line ---
LINE=$(echo "$ACTUAL" | sed -n '4p')
if [ "$LINE" = "IsAwesome: true" ]; then
    pass "Line 4: IsAwesome: true"
else
    fail "Line 4 wrong" "Expected: 'IsAwesome: true', Got: '${LINE}'"
fi

# --- Test 8: Exactly 4 lines ---
LINE_COUNT=$(echo "$ACTUAL" | wc -l | tr -d ' ')
if [ "$LINE_COUNT" = "4" ]; then
    pass "Exactly 4 lines of output"
else
    fail "Wrong number of lines" "Expected 4 lines, got ${LINE_COUNT}"
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
