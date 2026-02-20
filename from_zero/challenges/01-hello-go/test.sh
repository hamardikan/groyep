#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASS=0
FAIL=0
BINARY="hello"

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
echo -e "${YELLOW}=== Challenge 01: Hello, Go! ===${NC}"
echo ""

# --- Test 1: go.mod exists ---
echo "Checking project structure..."
if [ -f "go.mod" ]; then
    pass "go.mod exists"
else
    fail "go.mod not found" "Run: go mod init hello"
fi

# --- Test 2: main.go exists ---
if [ -f "main.go" ]; then
    pass "main.go exists"
else
    fail "main.go not found" "Create a main.go file"
fi

# --- Test 3: Build succeeds ---
echo "Building..."
if go build -o "${BINARY}" . 2>/dev/null; then
    pass "Build succeeded"
else
    fail "Build failed" "Fix syntax errors: go build -o ${BINARY} ."
fi

# --- Test 4: Full output matches ---
if [ -f "${BINARY}" ]; then
    echo "Running..."
    ACTUAL=$(./"${BINARY}" 2>&1)
    EXPECTED=$'Hello, World!\nHello, Go!'

    if [ "$ACTUAL" = "$EXPECTED" ]; then
        pass "Output matches exactly"
    else
        fail "Output mismatch" "Expected:\n${EXPECTED}\nGot:\n${ACTUAL}"
    fi

    # --- Test 5: First line ---
    FIRST_LINE=$(echo "$ACTUAL" | head -n 1)
    if [ "$FIRST_LINE" = "Hello, World!" ]; then
        pass "First line is 'Hello, World!'"
    else
        fail "First line wrong" "Expected: 'Hello, World!', Got: '${FIRST_LINE}'"
    fi

    # --- Test 6: Second line ---
    SECOND_LINE=$(echo "$ACTUAL" | sed -n '2p')
    if [ "$SECOND_LINE" = "Hello, Go!" ]; then
        pass "Second line is 'Hello, Go!'"
    else
        fail "Second line wrong" "Expected: 'Hello, Go!', Got: '${SECOND_LINE}'"
    fi

    # --- Test 7: Exactly 2 lines ---
    LINE_COUNT=$(echo "$ACTUAL" | wc -l | tr -d ' ')
    if [ "$LINE_COUNT" = "2" ]; then
        pass "Exactly 2 lines of output"
    else
        fail "Wrong number of lines" "Expected 2 lines, got ${LINE_COUNT}"
    fi
else
    fail "Binary not found" "Build must succeed first"
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
